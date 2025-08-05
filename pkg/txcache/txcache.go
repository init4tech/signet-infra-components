package txcache

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewTxCacheComponent(ctx *pulumi.Context, args TxCacheComponentArgs, opts ...pulumi.ResourceOption) (*TxCacheComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction cache component args: %w", err)
	}

	internalArgs := args.toInternal()

	component := &TxCacheComponent{
		ResourceState: pulumi.ResourceState{},
	}

	// Register the component with Pulumi
	err := ctx.RegisterComponentResource(ComponentResourceType, args.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register transaction cache component: %w", err)
	}

	appLabels := pulumi.StringMap{
		"app": pulumi.String(AppLabel),
	}

	serviceAccount, err := corev1.NewServiceAccount(ctx, ServiceAccountName, &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}

	// Create ConfigMap for environment variables
	configMapName := fmt.Sprintf("%s-env", args.Name)
	configMap, err := utils.CreateConfigMap(
		ctx,
		configMapName,
		internalArgs.Namespace,
		appLabels,
		internalArgs.Env,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment ConfigMap: %w", err)
	}

	// create the deployment for the quincey-server container to use the KMS key
	txCacheDeployment, err := appsv1.NewDeployment(ctx, DeploymentName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: pulumi.Int(ReplicaCount),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: serviceAccount.Metadata.Name(),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String(ContainerName),
							Image:           internalArgs.Image,
							ImagePullPolicy: pulumi.String(ImagePullPolicy),
							EnvFrom: corev1.EnvFromSourceArray{
								&corev1.EnvFromSourceArgs{
									ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
										Name: configMap.Metadata.Name(),
									},
								},
							},
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: internalArgs.Port,
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{configMap}), pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	// create a service for the deployment that allows the quincey-server to be accessed from the public internet over port 8080
	txCacheService, err := corev1.NewService(ctx, ServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: appLabels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       internalArgs.Port,
					TargetPort: internalArgs.Port,
				},
			},
			Type: pulumi.String(ServiceTypeClusterIP),
		},
	}, pulumi.DependsOn([]pulumi.Resource{txCacheDeployment}))
	if err != nil {
		return nil, err
	}

	virtualService, err := crd.NewCustomResource(ctx, VirtualServiceName, &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(VirtualServiceAPIVersion),
		Kind:       pulumi.String(VirtualServiceKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(AppLabel),
			Namespace: internalArgs.Namespace,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"hosts": []string{
					VirtualServiceHost,
				},
				"gateways": []string{
					GatewayName,
				},
				"http": []map[string]interface{}{
					{
						"match": []map[string]interface{}{
							{
								"uri": map[string]interface{}{
									"prefix": UriPrefix,
								},
							},
						},
						"route": []map[string]interface{}{
							{
								"destination": map[string]interface{}{
									"host": txCacheService.Metadata.Name(),
									"port": map[string]interface{}{
										"number": internalArgs.Port,
									},
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{txCacheService}))
	if err != nil {
		return nil, err
	}

	// create an RequestAuthentication policy for the virtual service
	// oauth access token is required to accees the service
	requestAuthentication, err := crd.NewCustomResource(ctx, JwtPolicyName, &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(RequestAuthAPIVersion),
		Kind:       pulumi.String(RequestAuthKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(JwtPolicyName),
			Namespace: internalArgs.Namespace,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": appLabels,
				},
				"jwtRules": []map[string]interface{}{
					{
						"issuer":  internalArgs.OauthIssuer,
						"jwksUri": internalArgs.OauthJwksUri,
						// List of output claim header objects
						"outputClaimToHeaders": []map[string]interface{}{
							{
								"claim":  JwtClaimSub,
								"header": JwtHeaderSub,
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{txCacheService}))
	if err != nil {
		return nil, err
	}

	// create a policy for the virtual service
	// only allow requests with a valid jwt token
	authorizationPolicy, err := crd.NewCustomResource(ctx, AuthPolicyName, &crd.CustomResourceArgs{
		ApiVersion: pulumi.String(AuthPolicyAPIVersion),
		Kind:       pulumi.String(AuthPolicyKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(AuthPolicyName),
			Namespace: internalArgs.Namespace,
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": appLabels,
				},
				"action": AuthActionAllow,
				"rules": []map[string]interface{}{
					{
						"from": []map[string]interface{}{
							{
								"source": map[string]interface{}{
									"requestPrincipals": []string{
										RequestPrincipalWildcard,
									},
								},
							},
						},
						"to": []map[string]interface{}{
							{
								"operation": map[string]interface{}{
									"paths": []string{
										BundlesPath,
										BundlesWildcardPath,
									},
									"methods": []string{
										HttpMethodGet,
									},
								},
							},
						},
					},
					{
						"to": []map[string]interface{}{
							{
								"operation": map[string]interface{}{
									"paths": []string{
										BundlesPath,
										BundlesWildcardPath,
									},
									"notMethods": []string{
										HttpMethodGet,
									},
								},
							},
						},
					},
					{
						"to": []map[string]interface{}{
							{
								"operation": map[string]interface{}{
									"notPaths": []string{
										BundlesPath,
										BundlesWildcardPath,
									},
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{txCacheService}))
	if err != nil {
		return nil, err
	}

	component.ServiceAccount = serviceAccount
	component.ConfigMap = configMap
	component.Deployment = txCacheDeployment
	component.Service = txCacheService
	component.VirtualService = virtualService
	component.RequestAuthentication = requestAuthentication
	component.AuthorizationPolicy = authorizationPolicy

	return component, nil
}
