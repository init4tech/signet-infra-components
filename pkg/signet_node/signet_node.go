package signet_node

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewSignetNode(ctx *pulumi.Context, args SignetNodeComponentArgs, opts ...pulumi.ResourceOption) (*SignetNodeComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid signet node component args: %w", err)
	}

	// Convert public args to internal args for use with Pulumi
	internalArgs := args.toInternal()

	component := &SignetNodeComponent{
		SignetNodeComponentArgs: args,
	}
	err := ctx.RegisterComponentResource("signet:index:SignetNode", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	hostDatabasePvc, err := CreatePersistentVolumeClaim(
		ctx,
		"signet-node-data",
		internalArgs.Namespace,
		internalArgs.ExecutionPvcSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signet node db data pvc: %w", err)
	}

	rollupDatabasePvc, err := CreatePersistentVolumeClaim(
		ctx,
		"rollup-data",
		internalArgs.Namespace,
		internalArgs.RollupPvcSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rollup data pvc: %w", err)
	}

	secretName := "execution-jwt"
	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": internalArgs.ExecutionJwt,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, secretName, args.Name, nil),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create execution jwt secret: %w", err)
	}

	// Create ConfigMap for execution environment variables
	executionConfigMapName := "exex-configmap"
	executionConfigMap, err := utils.CreateConfigMap(
		ctx,
		executionConfigMapName,
		internalArgs.Namespace,
		utils.CreateResourceLabels(args.Name, executionConfigMapName, args.Name, nil),
		internalArgs.Env,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution configmap: %w", err)
	}
	component.SignetNodeConfigMap = executionConfigMap

	// SERVICE
	executionClientName := "signet-node"
	executionClientServiceName := fmt.Sprintf("%s-service", executionClientName)

	signetNodeService, err := corev1.NewService(ctx, executionClientServiceName, &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.String("signet-node-execution-set")},
			Type:     pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port: pulumi.Int(DiscoveryPort),
					Name: pulumi.String("disc"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(MetricsPort),
					Name: pulumi.String("metrics"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(RpcPort),
					Name: pulumi.String("host-http"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(AuthRpcPort),
					Name: pulumi.String("p2p"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(WsPort),
					Name: pulumi.String("host-ws"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(HostIpcPort),
					Name: pulumi.String("host-ipc"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(RollupHttpPort),
					Name: pulumi.String("rollup-http"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(RollupWsPort),
					Name: pulumi.String("rollup-ws"),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, executionClientServiceName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client service: %w", err)
	}

	// STATEFUL SET
	hostStatefulSetName := "signet-node-execution-set"
	hostStatefulSetResourceName := fmt.Sprintf("%s-set", hostStatefulSetName)

	// Create pod labels with app label for stateful set
	executionPodLabels := utils.CreateResourceLabels(args.Name, hostStatefulSetName, args.Name, nil)
	executionPodLabels["app"] = pulumi.String(hostStatefulSetName)

	// Define the StatefulSet for the 'reth' container with a configmap volume and a data persistent volume
	signetNodeStatefulSet, err := appsv1.NewStatefulSet(ctx, hostStatefulSetResourceName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    utils.CreateResourceLabels(args.Name, hostStatefulSetResourceName, args.Name, nil),
			Namespace: internalArgs.Namespace,
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(hostStatefulSetName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels:    executionPodLabels,
					Namespace: internalArgs.Namespace,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String(hostStatefulSetName),
							Image:           internalArgs.ExecutionClientImage,
							ImagePullPolicy: pulumi.String("Always"),
							Command:         internalArgs.ExecutionClientStartCommand,
							EnvFrom: corev1.EnvFromSourceArray{
								corev1.EnvFromSourceArgs{
									ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
										Name:     executionConfigMap.Metadata.Name(),
										Optional: pulumi.Bool(true),
									},
								},
							},
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(30303),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(30303),
									Protocol:      pulumi.String("UDP"),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(9001),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8545),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8551),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8546),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8547),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8645),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8646),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									Name:      pulumi.String("signet-node-data"),
									MountPath: pulumi.String("/root/.local/share/reth"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.String("rollup-data"),
									MountPath: pulumi.String("/root/.local/share/exex"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.String("execution-jwt"),
									MountPath: pulumi.String("/etc/reth/execution-jwt"),
								},
							},
							Resources: NewResourceRequirements("2", "16Gi", "2", "4Gi"),
						},
					},
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.String("signet-node-data"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: hostDatabasePvc.Metadata.Name().Elem().ToStringOutput(),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("execution-jwt"),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: secret.Metadata.Name(),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("rollup-data"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: rollupDatabasePvc.Metadata.Name().Elem().ToStringOutput(),
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client statefulset: %w", err)
	}

	// LIGHTHOUSE
	consensusClientName := "lighthouse"

	lighthousePvc, err := CreatePersistentVolumeClaim(
		ctx,
		"real-lighthouse-data",
		internalArgs.Namespace,
		internalArgs.LighthousePvcSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse data pvc: %w", err)
	}

	// Create ConfigMap for consensus environment variables
	consensusEnv := ConsensusEnv{
		Example: "example",
	}

	consensusConfigMapName := "consensus-configmap-env-config"
	consensusConfigMap, err := utils.CreateConfigMap(
		ctx,
		consensusConfigMapName,
		internalArgs.Namespace,
		utils.CreateResourceLabels(args.Name, consensusConfigMapName, args.Name, nil),
		consensusEnv.toInternal(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus configmap: %w", err)
	}
	component.LighthouseConfigMap = consensusConfigMap

	lighthouseServiceName := fmt.Sprintf("%s-service", consensusClientName)
	lighthouseService, err := corev1.NewService(ctx, lighthouseServiceName, &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.String(consensusClientName)},
			Type:     pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:     pulumi.Int(9000),
					Name:     pulumi.String("udp"),
					Protocol: pulumi.String("UDP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(9000),
					Name:     pulumi.String("tcp"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(9001),
					Name:     pulumi.String("udp2"),
					Protocol: pulumi.String("UDP"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(5054),
					Name: pulumi.String("metrics"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(8551),
					Name: pulumi.String("exec"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(5052),
					Name: pulumi.String("http"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(4000),
					Name: pulumi.String("http-cl"),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, lighthouseServiceName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse internal service: %w", err)
	}

	lighthouseStatefulSetName := "lighthouse"

	// Create pod labels with app label for stateful set
	lighthousePodLabels := utils.CreateResourceLabels(args.Name, lighthouseStatefulSetName, args.Name, nil)
	lighthousePodLabels["app"] = pulumi.String(lighthouseStatefulSetName)

	signetNodeServiceName := signetNodeService.Metadata.Name().Elem().ToStringOutput()
	consensusClientStartCommand := signetNodeServiceName.ApplyT(func(name string) []string {
		return append(args.ConsensusClientStartCommand, fmt.Sprintf("--execution-endpoints=http://%s:8551", name))
	}).(pulumi.StringArrayOutput)

	lighthouseStatefulSet, err := appsv1.NewStatefulSet(ctx, lighthouseStatefulSetName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    utils.CreateResourceLabels(args.Name, lighthouseStatefulSetName, args.Name, nil),
			Namespace: internalArgs.Namespace,
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(lighthouseStatefulSetName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Namespace: internalArgs.Namespace,
					Labels:    lighthousePodLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:    pulumi.String(lighthouseStatefulSetName),
							Image:   internalArgs.ConsensusClientImage,
							Command: consensusClientStartCommand,
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8551),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(9000),
									Protocol:      pulumi.String("UDP"),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(9001),
									Protocol:      pulumi.String("UDP"),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(5054),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(5052),
								},
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(4000),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									Name:      pulumi.Sprintf("%s-config", lighthouseStatefulSetName),
									MountPath: pulumi.String("/etc/lighthouse"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.String("real-lighthouse-db"),
									MountPath: pulumi.String("/root/.lighthouse/holesky"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.Sprintf("%s-execution-jwt", lighthouseStatefulSetName),
									MountPath: pulumi.String("/secrets"),
								},
							},
							Resources: NewResourceRequirements("2", "16Gi", "2", "4Gi"),
						},
					},
					DnsPolicy: pulumi.String("ClusterFirst"),
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-config", lighthouseStatefulSetName),
							ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
								Name: consensusConfigMap.Metadata.Name(),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("real-lighthouse-db"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: lighthousePvc.Metadata.Name().Elem().ToStringOutput(),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-execution-jwt", lighthouseStatefulSetName),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: secret.Metadata.Name(),
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse statefulset: %w", err)
	}

	// Create a VirtualService resource to route traffic to the signet nodes
	// This enables the service mesh to route traffic from rpc.havarti.signet.sh
	// to the signet-rpc service in the cluster
	// VirtualService spec definition: https://istio.io/latest/docs/reference/config/networking/virtual-service/
	virtualServiceName := "signet-rpc"
	signetNodeServiceUrl := signetNodeServiceName.ApplyT(func(name string) string {
		return fmt.Sprintf("%s.%s.svc.cluster.local", name, internalArgs.Namespace)
	})
	virtualService, err := crd.NewCustomResource(ctx, "signet-rpc-vservice", &crd.CustomResourceArgs{
		ApiVersion: pulumi.String("networking.istio.io/v1alpha3"),
		Kind:       pulumi.String("VirtualService"),
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, virtualServiceName, args.Name, nil),
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"hosts": []string{
					"rpc.mainnet.signet.sh",
				},
				"gateways": []string{
					"default/init4-api-gateway",
				},
				"http": []map[string]interface{}{
					{
						"match": []map[string]interface{}{
							{
								"uri": map[string]interface{}{
									"prefix": "/",
								},
							},
							{
								"uri": map[string]interface{}{
									"prefix": "/healthcheck",
								},
							},
						},
						"route": []map[string]interface{}{
							{
								"destination": map[string]interface{}{
									"port": map[string]interface{}{
										"number": args.Env.RpcPort,
									},
									"host": signetNodeServiceUrl,
								},
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create signet rpc virtual service: %w", err)
	}

	component.SignetNodeService = signetNodeService
	component.LighthouseService = lighthouseService
	component.SignetNodeStatefulSet = signetNodeStatefulSet
	component.LighthouseStatefulSet = lighthouseStatefulSet
	component.SignetNodeConfigMap = executionConfigMap
	component.LighthouseConfigMap = consensusConfigMap
	component.JwtSecret = secret
	component.HostDatabasePvc = hostDatabasePvc
	component.RollupDatabasePvc = rollupDatabasePvc
	component.LighthousePvc = lighthousePvc
	component.SignetNodeVirtualService = virtualService

	return component, nil
}
