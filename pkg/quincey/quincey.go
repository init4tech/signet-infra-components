package quincey

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewQuinceyComponent creates a new Quincey deployment in the specified namespace.
// It sets up all necessary Kubernetes resources including the deployment, service,
// and Istio configurations.
//
// Example:
//
//	component, err := NewQuinceyComponent(ctx, "quincey", &QuinceyComponentArgs{
//	    Namespace: pulumi.String("default"),
//	    Image:     pulumi.String("quincey:latest"),
//	    Env: QuinceyEnv{
//	        QuinceyPort: "8080",
//	        // ... other environment variables
//	    },
//	})
func NewQuinceyComponent(ctx *pulumi.Context, name string, args *QuinceyComponentArgs, opts ...pulumi.ResourceOption) (*QuinceyComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid quincey component args: %w", err)
	}

	// Convert public args to internal args for use with Pulumi
	internalArgs := args.toInternal()

	component := &QuinceyComponent{
		ResourceState: pulumi.ResourceState{},
	}

	if err := ctx.RegisterComponentResource(ComponentKind, name, component); err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create service account
	serviceAccount, err := createServiceAccount(ctx, internalArgs.Namespace, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}

	// Create config map
	configMap, err := createConfigMap(ctx, &internalArgs, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create config map: %w", err)
	}

	// Create deployment
	deployment, err := createDeployment(ctx, &internalArgs, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Create service
	service, err := createService(ctx, &internalArgs, deployment, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	// Create virtual service
	virtualService, err := createVirtualService(ctx, &internalArgs, service, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual service: %w", err)
	}

	// Create request authentication
	requestAuth, err := createRequestAuthentication(ctx, &internalArgs, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create request authentication: %w", err)
	}

	// Create authorization policy
	authPolicy, err := createAuthorizationPolicy(ctx, &internalArgs, component)
	if err != nil {
		return nil, fmt.Errorf("failed to create authorization policy: %w", err)
	}

	component.Service = service
	component.ServiceAccount = serviceAccount
	component.ConfigMap = configMap
	component.Deployment = deployment
	component.VirtualService = virtualService
	component.RequestAuthentication = requestAuth
	component.AuthorizationPolicy = authPolicy

	return component, nil
}

// createServiceAccount creates a Kubernetes service account for the Quincey service
func createServiceAccount(ctx *pulumi.Context, namespace pulumi.StringInput, parent *QuinceyComponent) (*corev1.ServiceAccount, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	return corev1.NewServiceAccount(ctx, ServiceName, &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(ServiceName),
			Namespace: namespace,
			Labels:    labels,
		},
	}, pulumi.Parent(parent))
}

// createDeployment creates the Kubernetes deployment for the Quincey service
func createDeployment(ctx *pulumi.Context, args *quinceyComponentArgsInternal, parent *QuinceyComponent) (*appsv1.Deployment, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	containerPortInt := utils.ParsePortWithDefault(args.Env.QuinceyPort, DefaultQuinceyPort)

	return appsv1.NewDeployment(ctx, ServiceName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(ServiceName),
			Namespace: args.Namespace,
			Labels:    labels,
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: labels,
			},
			Replicas: pulumi.Int(DefaultReplicas),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: labels,
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String(ServiceName),
					Containers: corev1.ContainerArray{
						createContainer(args, containerPortInt),
					},
				},
			},
		},
	}, pulumi.Parent(parent))
}

// createConfigMap creates the ConfigMap for the Quincey service
func createConfigMap(ctx *pulumi.Context, args *quinceyComponentArgsInternal, parent *QuinceyComponent) (*corev1.ConfigMap, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	return utils.CreateConfigMap(ctx, ServiceName, args.Namespace, labels, &args.Env)
}

// createContainer creates the container specification for the Quincey service
func createContainer(args *quinceyComponentArgsInternal, port pulumi.IntOutput) *corev1.ContainerArgs {
	return &corev1.ContainerArgs{
		Name:  pulumi.String(ServiceName),
		Image: args.Image,
		EnvFrom: corev1.EnvFromSourceArray{
			&corev1.EnvFromSourceArgs{
				ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
					Name: pulumi.String(ServiceName),
				},
			},
		},
	}
}

// createService creates the Kubernetes service for the Quincey service
func createService(ctx *pulumi.Context, args *quinceyComponentArgsInternal, deployment *appsv1.Deployment, parent *QuinceyComponent) (*corev1.Service, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	containerPortInt := utils.ParsePortWithDefault(args.Env.QuinceyPort, DefaultQuinceyPort)

	return corev1.NewService(ctx, "quincey-server-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("quincey"),
			Namespace: args.Namespace,
			Labels:    labels,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: labels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       containerPortInt,
					TargetPort: containerPortInt,
				},
			},
			Type: pulumi.String(ServiceTypeClusterIP),
		},
	}, pulumi.DependsOn([]pulumi.Resource{deployment}), pulumi.Parent(parent))
}

// createVirtualService creates the Istio virtual service for the Quincey service
func createVirtualService(ctx *pulumi.Context, args *quinceyComponentArgsInternal, service *corev1.Service, parent *QuinceyComponent) (*crd.CustomResource, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	containerPortInt := utils.ParsePortWithDefault(args.Env.QuinceyPort, DefaultQuinceyPort)

	// Get the service URL using the existing method
	serviceURL := parent.GetServiceURL("quincey", args.Namespace)

	return crd.NewCustomResource(ctx, "quincey-vservice", &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(IstioNetworkingAPIVersion),
		Kind:       pulumi.String(VirtualServiceKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("quincey"),
			Namespace: args.Namespace,
			Labels:    labels,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"hosts": args.VirtualServiceHosts,
				"gateways": []string{
					"default/init4-api-gateway",
				},
				"http": []map[string]interface{}{
					{
						"match": []map[string]interface{}{
							{
								"uri": map[string]interface{}{
									"prefix": "/signBlock",
								},
							},
							{
								"uri": map[string]interface{}{
									"prefix": "/healthCheck",
								},
							},
						},
						"route": []map[string]interface{}{
							{
								"destination": map[string]interface{}{
									"host": serviceURL,
									"port": map[string]interface{}{
										"number": containerPortInt,
									},
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{service}), pulumi.Parent(parent))
}

// createRequestAuthentication creates the Istio request authentication policy
func createRequestAuthentication(ctx *pulumi.Context, args *quinceyComponentArgsInternal, parent *QuinceyComponent) (*crd.CustomResource, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	return crd.NewCustomResource(ctx, "quincey-authorization-policy", &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(IstioSecurityAPIVersion),
		Kind:       pulumi.String(RequestAuthenticationKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("quincey-jwt-policy"),
			Namespace: args.Namespace,
			Labels:    labels,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": labels,
				},
				"jwtRules": []map[string]interface{}{
					{
						"issuer":  args.Env.OauthIssuer,
						"jwksUri": args.Env.OauthJwksUri,
						"outputClaimToHeaders": []map[string]interface{}{
							{
								"claim":  "sub",
								"header": "x-jwt-claim-sub",
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(parent))
}

// createAuthorizationPolicy creates the Istio authorization policy
func createAuthorizationPolicy(ctx *pulumi.Context, args *quinceyComponentArgsInternal, parent *QuinceyComponent) (*crd.CustomResource, error) {
	labels := utils.CreateResourceLabels(ComponentName, ServiceName, DefaultAppSelector, nil)

	return crd.NewCustomResource(ctx, "quincey-authorization-policy", &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(IstioSecurityAPIVersion),
		Kind:       pulumi.String(AuthorizationPolicyKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("quincey-jwt-auth-policy"),
			Namespace: args.Namespace,
			Labels:    labels,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": labels,
				},
				"action": "ALLOW",
				"rules": []map[string]interface{}{
					{
						"from": []map[string]interface{}{
							{
								"source": map[string]interface{}{
									"requestPrincipals": []string{
										"*",
									},
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(parent))
}

// GetServiceURL returns the URL of the builder service
func (c *QuinceyComponent) GetServiceURL(name string, namespace pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf("%s.%s.svc.cluster.local", name, namespace)
}

// GetMetricsURL returns the URL of the builder metrics endpoint
func (c *QuinceyComponent) GetMetricsURL(name string, namespace pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf("http://%s.%s.svc.cluster.local:%d/metrics", name, namespace, DefaultMetricsPort)
}
