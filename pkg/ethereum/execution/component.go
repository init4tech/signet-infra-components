package execution

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewExecutionClient creates a new execution client component
func NewExecutionClient(ctx *pulumi.Context, args *ExecutionClientArgs, opts ...pulumi.ResourceOption) (*ExecutionClientComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid execution client args: %w", err)
	}

	// Convert public args to internal args for use with Pulumi
	internalArgs := args.toInternal()

	component := &ExecutionClientComponent{
		Name:      args.Name,
		Namespace: args.Namespace,
	}

	err := ctx.RegisterComponentResource("signet:index:ExecutionClient", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create PVC for data storage
	pvcName := fmt.Sprintf("%s-data", args.Name)
	component.PVC, err = corev1.NewPersistentVolumeClaim(ctx, pvcName, &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(pvcName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, pvcName, args.Name, nil),
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{
				pulumi.String("ReadWriteOnce"),
			},
			Resources: &corev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": internalArgs.StorageSize,
				},
			},
			StorageClassName: internalArgs.StorageClass,
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create PVC: %w", err)
	}

	// Create JWT secret
	jwtSecretName := fmt.Sprintf("%s-jwt", args.Name)
	component.JWTSecret, err = corev1.NewSecret(ctx, jwtSecretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": internalArgs.JWTSecret,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(jwtSecretName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, jwtSecretName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT secret: %w", err)
	}

	// Create P2P service
	p2pServiceName := fmt.Sprintf("%s-p2p", args.Name)
	component.P2PService, err = corev1.NewService(ctx, p2pServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(p2pServiceName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, p2pServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("p2p"),
					Port:       internalArgs.P2PPort,
					TargetPort: internalArgs.P2PPort,
					Protocol:   pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("discovery"),
					Port:       internalArgs.DiscoveryPort,
					TargetPort: internalArgs.DiscoveryPort,
					Protocol:   pulumi.String("UDP"),
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create P2P service: %w", err)
	}

	// Create RPC service
	rpcServiceName := fmt.Sprintf("%s-rpc", args.Name)
	component.RPCService, err = corev1.NewService(ctx, rpcServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(rpcServiceName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, rpcServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("rpc"),
					Port:       internalArgs.RPCPort,
					TargetPort: internalArgs.RPCPort,
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("ws"),
					Port:       internalArgs.WSPort,
					TargetPort: internalArgs.WSPort,
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("metrics"),
					Port:       internalArgs.MetricsPort,
					TargetPort: internalArgs.MetricsPort,
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("auth-rpc"),
					Port:       internalArgs.AuthRPCPort,
					TargetPort: internalArgs.AuthRPCPort,
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC service: %w", err)
	}

	// Create ConfigMap for environment variables if ExecutionClientEnv is provided
	var configMap *corev1.ConfigMap
	if args.ExecutionClientEnv != nil {
		configMapName := fmt.Sprintf("%s-env", args.Name)
		configMap, err = corev1.NewConfigMap(ctx, configMapName, &corev1.ConfigMapArgs{
			Data: args.ExecutionClientEnv.GetEnvMap(),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(configMapName),
				Namespace: internalArgs.Namespace,
				Labels:    utils.CreateResourceLabels(args.Name, configMapName, args.Name, nil),
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, fmt.Errorf("failed to create ConfigMap: %w", err)
		}
		component.ConfigMap = configMap
	}

	// Create StatefulSet
	statefulSetName := args.Name

	// Prepare container spec
	containerSpec := corev1.ContainerArgs{
		Name:            pulumi.String("execution"),
		Image:           internalArgs.Image,
		ImagePullPolicy: internalArgs.ImagePullPolicy,
		Command:         createExecutionClientCommand(args),
		Ports: corev1.ContainerPortArray{
			corev1.ContainerPortArgs{
				Name:          pulumi.String("p2p"),
				ContainerPort: internalArgs.P2PPort,
				Protocol:      pulumi.String("TCP"),
			},
			corev1.ContainerPortArgs{
				Name:          pulumi.String("discovery"),
				ContainerPort: internalArgs.DiscoveryPort,
				Protocol:      pulumi.String("UDP"),
			},
			corev1.ContainerPortArgs{
				Name:          pulumi.String("rpc"),
				ContainerPort: internalArgs.RPCPort,
			},
			corev1.ContainerPortArgs{
				Name:          pulumi.String("ws"),
				ContainerPort: internalArgs.WSPort,
			},
			corev1.ContainerPortArgs{
				Name:          pulumi.String("metrics"),
				ContainerPort: internalArgs.MetricsPort,
			},
			corev1.ContainerPortArgs{
				Name:          pulumi.String("auth-rpc"),
				ContainerPort: internalArgs.AuthRPCPort,
			},
		},
		VolumeMounts: corev1.VolumeMountArray{
			corev1.VolumeMountArgs{
				Name:      pulumi.String("data"),
				MountPath: pulumi.String("/data"),
			},
			corev1.VolumeMountArgs{
				Name:      pulumi.String("jwt"),
				MountPath: pulumi.String("/jwt"),
			},
		},
		Resources: nil,
	}

	// Add EnvFrom only if ConfigMap exists
	if configMap != nil {
		containerSpec.EnvFrom = corev1.EnvFromSourceArray{
			&corev1.EnvFromSourceArgs{
				ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
					Name: configMap.Metadata.Name(),
				},
			},
		}
	}

	component.StatefulSet, err = appsv1.NewStatefulSet(ctx, statefulSetName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(statefulSetName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, statefulSetName, args.Name, nil),
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(args.Name),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String(args.Name),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						containerSpec,
					},
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.String("data"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: pulumi.String(pvcName),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("jwt"),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: pulumi.String(jwtSecretName),
							},
						},
					},
					NodeSelector: args.NodeSelector,
					Tolerations:  args.Tolerations,
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create StatefulSet: %w", err)
	}

	return component, nil
}

// createExecutionClientCommand creates the command array for the execution client
func createExecutionClientCommand(args *ExecutionClientArgs) pulumi.StringArray {
	cmd := pulumi.StringArray{
		pulumi.String("--datadir=/data"),
		pulumi.String("--http"),
		pulumi.Sprintf("--http.port=%d", args.RPCPort),
		pulumi.String("--http.addr=0.0.0.0"),
		pulumi.String("--http.corsdomain=*"),
		pulumi.String("--http.api=admin,net,eth,web3,debug,txpool,trace"),
		pulumi.String("--ws"),
		pulumi.Sprintf("--ws.port=%d", args.WSPort),
		pulumi.String("--ws.addr=0.0.0.0"),
		pulumi.String("--ws.api=net,eth"),
		pulumi.String("--ws.origins=*"),
		pulumi.Sprintf("--authrpc.port=%d", args.AuthRPCPort),
		pulumi.String("--authrpc.jwtsecret=/etc/execution/jwt/jwt.hex"),
		pulumi.String("--authrpc.addr=0.0.0.0"),
		pulumi.Sprintf("--metrics=0.0.0.0:%d", args.MetricsPort),
		pulumi.Sprintf("--discovery.port=%d", args.DiscoveryPort),
		pulumi.Sprintf("--port=%d", args.P2PPort),
	}

	// Add bootnodes if provided
	for _, bootnode := range args.Bootnodes {
		cmd = append(cmd, pulumi.Sprintf("--bootnodes=%s", bootnode))
	}

	// Add additional args if provided
	for _, arg := range args.AdditionalArgs {
		cmd = append(cmd, pulumi.String(arg))
	}

	return cmd
}
