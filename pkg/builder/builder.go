// Package builder provides a Pulumi component for deploying a builder service to Kubernetes.
package builder

import (
	"fmt"
	"strconv"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// parseBuilderPort converts a port string to an integer with a default fallback
func parseBuilderPort(portStr pulumi.StringInput) pulumi.IntOutput {
	return pulumi.All(portStr).ApplyT(func(inputs []interface{}) int {
		portStr := inputs[0].(string)
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
		// Default to 8080 if there's an error parsing the port
		return 8080
	}).(pulumi.IntOutput)
}

// NewBuilder creates a new builder component with the given configuration.
func NewBuilder(ctx *pulumi.Context, args BuilderComponentArgs, opts ...pulumi.ResourceOption) (*BuilderComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid builder component args: %w", err)
	}

	component := &BuilderComponent{
		BuilderComponentArgs: args,
	}
	err := ctx.RegisterComponentResource("signet:index:Builder", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create service account
	serviceAccountName := fmt.Sprintf("%s-sa", args.Name)
	sa, err := corev1.NewServiceAccount(ctx, serviceAccountName, &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceAccountName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, serviceAccountName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}
	component.ServiceAccount = sa

	// Create ConfigMap for environment variables
	configMapName := fmt.Sprintf("%s-env", args.Name)
	configMap, err := utils.CreateConfigMap(
		ctx,
		configMapName,
		pulumi.String(args.Namespace),
		utils.CreateResourceLabels(args.Name, configMapName, args.Name, nil),
		args.BuilderEnv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment ConfigMap: %w", err)
	}
	component.ConfigMap = configMap

	// Create pod labels with app label for routing
	podLabels := utils.CreateResourceLabels(args.Name, args.Name, args.Name, args.AppLabels.Labels)
	podLabels["app"] = pulumi.String(args.Name)

	// Create deployment
	deploymentName := fmt.Sprintf("%s-deployment", args.Name)
	deployment, err := appsv1.NewDeployment(ctx, deploymentName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(deploymentName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, deploymentName, args.Name, nil),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(DefaultReplicas),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: podLabels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: podLabels,
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String(serviceAccountName),
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
									ContainerPort: parseBuilderPort(args.BuilderEnv.BuilderPort),
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
									Port: parseBuilderPort(args.BuilderEnv.BuilderPort),
								},
								InitialDelaySeconds: pulumi.Int(5),
								PeriodSeconds:       pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
								FailureThreshold:    pulumi.Int(3),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: parseBuilderPort(args.BuilderEnv.BuilderPort),
								},
								InitialDelaySeconds: pulumi.Int(5),
								PeriodSeconds:       pulumi.Int(10),
							},
						},
					},
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true), pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}
	component.Deployment = deployment

	// Create service
	serviceName := fmt.Sprintf("%s-service", args.Name)
	service, err := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, serviceName, args.Name, nil),
			Annotations: pulumi.StringMap{
				"prometheus.io/scrape": pulumi.String("true"),
				"prometheus.io/port":   pulumi.Sprintf("%d", MetricsPort),
				"prometheus.io/path":   pulumi.String("/metrics"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: podLabels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       parseBuilderPort(args.BuilderEnv.BuilderPort),
					TargetPort: parseBuilderPort(args.BuilderEnv.BuilderPort),
					Name:       pulumi.String("http"),
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
	podMonitorName := fmt.Sprintf("%s-pod-monitor", args.Name)
	_, err = crd.NewCustomResource(ctx, fmt.Sprintf("%s-svcmon", args.Name), &crd.CustomResourceArgs{
		ApiVersion: pulumi.String("monitoring.coreos.com/v1"),
		Kind:       pulumi.String("PodMonitor"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(podMonitorName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, podMonitorName, args.Name, nil),
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": podLabels,
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
