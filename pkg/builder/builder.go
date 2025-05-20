// Package builder provides a Pulumi component for deploying a builder service to Kubernetes.
package builder

import (
	"encoding/json"
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewBuilder creates a new builder component with the given configuration.
func NewBuilder(ctx *pulumi.Context, args BuilderComponentArgs, opts ...pulumi.ResourceOption) (*BuilderComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid builder component args: %w", err)
	}

	component := &BuilderComponent{
		BuilderComponentArgs: args,
	}
	err := ctx.RegisterComponentResource("the-builder:index:Builder", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create service account
	sa, err := corev1.NewServiceAccount(ctx, fmt.Sprintf("%s-sa", args.Name), &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(fmt.Sprintf("%s-sa", args.Name)),
			Namespace: pulumi.String(args.Namespace),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}
	component.ServiceAccount = sa

	// Create IAM role
	assumeRolePolicy := IAMPolicy{
		Version: "2012-10-17",
		Statement: []IAMStatement{
			{
				Sid:    "AllowEksAuthToAssumeRoleForPodIdentity",
				Effect: "Allow",
				Principal: struct {
					Service []string `json:"Service"`
				}{
					Service: []string{
						"pods.eks.amazonaws.com",
						"ec2.amazonaws.com",
					},
				},
				Action: []string{
					"sts:AssumeRole",
					"sts:TagSession",
				},
			},
		},
	}

	assumeRolePolicyJSON, err := json.Marshal(assumeRolePolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal assume role policy: %w", err)
	}

	role, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", args.Name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
		Description:      pulumi.String(fmt.Sprintf("Role for %s pod to assume", args.Name)),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(fmt.Sprintf("%s-role", args.Name)),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM role: %w", err)
	}
	component.IAMRole = role

	// Create KMS policy
	policyJSON := CreateKMSPolicy(args.BuilderEnv.BuilderKey)

	policy, err := iam.NewPolicy(ctx, fmt.Sprintf("%s-policy", args.Name), &iam.PolicyArgs{
		Policy: policyJSON,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM policy: %w", err)
	}
	component.IAMPolicy = policy

	// Attach policy to role
	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-role-policy-attachment", args.Name), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: policy.Arn,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to attach policy to role: %w", err)
	}

	// Create ConfigMap for environment variables
	configMap, err := utils.CreateConfigMap(
		ctx,
		fmt.Sprintf("%s-env", args.Name),
		pulumi.String(args.Namespace),
		pulumi.StringMap{
			"app":                       pulumi.String(args.Name),
			"app.kubernetes.io/name":    pulumi.String(args.Name),
			"app.kubernetes.io/part-of": pulumi.String(args.Name),
		},
		args.BuilderEnv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment ConfigMap: %w", err)
	}
	component.ConfigMap = configMap

	// Create deployment
	deployment, err := appsv1.NewDeployment(ctx, fmt.Sprintf("%s-deployment", args.Name), &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(fmt.Sprintf("%s-deployment", args.Name)),
			Namespace: pulumi.String(args.Namespace),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(DefaultReplicas),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: args.AppLabels.Labels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: args.AppLabels.Labels,
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String(fmt.Sprintf("%s-sa", args.Name)),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String(args.Name),
							Image: pulumi.String(args.Image),
							EnvFrom: corev1.EnvFromSourceArray{
								&corev1.EnvFromSourceArgs{
									ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
										Name: component.ConfigMap.Metadata.Name(),
									},
								},
							},
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.All(args.BuilderEnv.BuilderPort).ApplyT(func(inputs []interface{}) int {
										portStr := inputs[0].(string)
										var port int
										if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
											// Default to 8080 if there's an error parsing the port
											return 8080
										}
										return port
									}).(pulumi.IntInput),
								},
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(MetricsPort),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"cpu":    pulumi.String("2"),
									"memory": pulumi.String("2Gi"),
								},
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String("1"),
									"memory": pulumi.String("1Gi"),
								},
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: pulumi.All(args.BuilderEnv.BuilderPort).ApplyT(func(inputs []interface{}) int {
										portStr := inputs[0].(string)
										var port int
										if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
											// Default to 8080 if there's an error parsing the port
											return 8080
										}
										return port
									}).(pulumi.IntInput),
								},
								InitialDelaySeconds: pulumi.Int(5),
								PeriodSeconds:       pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
								FailureThreshold:    pulumi.Int(3),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: pulumi.All(args.BuilderEnv.BuilderPort).ApplyT(func(inputs []interface{}) int {
										portStr := inputs[0].(string)
										var port int
										if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
											// Default to 8080 if there's an error parsing the port
											return 8080
										}
										return port
									}).(pulumi.IntInput),
								},
								InitialDelaySeconds: pulumi.Int(5),
								PeriodSeconds:       pulumi.Int(10),
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{role, policy}), pulumi.DeleteBeforeReplace(true), pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}
	component.Deployment = deployment

	// Create service
	service, err := corev1.NewService(ctx, fmt.Sprintf("%s-service", args.Name), &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(fmt.Sprintf("%s-service", args.Name)),
			Namespace: pulumi.String(args.Namespace),
			Annotations: pulumi.StringMap{
				"prometheus.io/scrape": pulumi.String("true"),
				"prometheus.io/port":   pulumi.Sprintf("%d", MetricsPort),
				"prometheus.io/path":   pulumi.String("/metrics"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: args.AppLabels.Labels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port: pulumi.All(args.BuilderEnv.BuilderPort).ApplyT(func(inputs []interface{}) int {
						portStr := inputs[0].(string)
						var port int
						if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
							// Default to 8080 if there's an error parsing the port
							return 8080
						}
						return port
					}).(pulumi.IntInput),
					TargetPort: pulumi.All(args.BuilderEnv.BuilderPort).ApplyT(func(inputs []interface{}) int {
						portStr := inputs[0].(string)
						var port int
						if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
							// Default to 8080 if there's an error parsing the port
							return 8080
						}
						return port
					}).(pulumi.IntInput),
					Name: pulumi.String("http"),
				},
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(MetricsPort),
					TargetPort: pulumi.Int(MetricsPort),
					Name:       pulumi.String("metrics"),
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{deployment}), pulumi.DeleteBeforeReplace(true), pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	component.Service = service

	// Create pod monitor
	_, err = crd.NewCustomResource(ctx, fmt.Sprintf("%s-svcmon", args.Name), &crd.CustomResourceArgs{
		ApiVersion: pulumi.String("monitoring.coreos.com/v1"),
		Kind:       pulumi.String("PodMonitor"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(fmt.Sprintf("%s-pod-monitor", args.Name)),
			Namespace: pulumi.String(args.Namespace),
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": args.AppLabels.Labels,
				},
				"namespaceSelector": map[string]interface{}{
					"any": true,
				},
				"podMetricsEndpoints": []map[string]interface{}{
					{
						"port": "metrics",
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create pod monitor: %w", err)
	}

	return component, nil
}

// GetServiceURL returns the URL of the builder service
func (c *BuilderComponent) GetServiceURL() pulumi.StringOutput {
	return pulumi.Sprintf("http://%s.%s.svc.cluster.local", c.Service.Metadata.Name(), c.Service.Metadata.Namespace())
}

// GetMetricsURL returns the URL of the builder metrics endpoint
func (c *BuilderComponent) GetMetricsURL() pulumi.StringOutput {
	return pulumi.Sprintf("http://%s.%s.svc.cluster.local:%d/metrics",
		c.Service.Metadata.Name(),
		c.Service.Metadata.Namespace(),
		MetricsPort)
}
