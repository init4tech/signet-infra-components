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

	component := &ExecutionClientComponent{
		Name:      args.Name,
		Namespace: args.Namespace,
	}

	err := ctx.RegisterComponentResource("ethereum:index:ExecutionClient", args.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create PVC for data storage
	pvcName := fmt.Sprintf("%s-data", args.Name)
	pvc, err := utils.CreatePersistentVolumeClaim(
		ctx,
		pvcName,
		pulumi.String(args.Namespace),
		pulumi.String(args.StorageSize),
		args.StorageClass,
		utils.CreateResourceLabels(args.Name, pvcName, args.Name, nil),
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pvc: %w", err)
	}
	component.PVC = pvc

	// Create JWT secret
	jwtSecretName := fmt.Sprintf("%s-jwt", args.Name)
	jwtSecret, err := corev1.NewSecret(ctx, jwtSecretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": pulumi.String(args.JWTSecret),
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(jwtSecretName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, jwtSecretName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt secret: %w", err)
	}
	component.JWTSecret = jwtSecret

	// Create P2P service
	p2pServiceName := fmt.Sprintf("%s-p2p", args.Name)
	p2pService, err := corev1.NewService(ctx, p2pServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(p2pServiceName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, p2pServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.P2PPort),
					Name:     pulumi.String("p2p-tcp"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.P2PPort),
					Name:     pulumi.String("p2p-udp"),
					Protocol: pulumi.String("UDP"),
				},
			},
			Type: pulumi.String("NodePort"),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create p2p service: %w", err)
	}
	component.P2PService = p2pService

	// Create RPC service
	rpcServiceName := fmt.Sprintf("%s-rpc", args.Name)
	rpcService, err := corev1.NewService(ctx, rpcServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(rpcServiceName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, rpcServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.RPCPort),
					Name:     pulumi.String("rpc"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.WSPort),
					Name:     pulumi.String("ws"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.MetricsPort),
					Name:     pulumi.String("metrics"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.AuthRPCPort),
					Name:     pulumi.String("auth-rpc"),
					Protocol: pulumi.String("TCP"),
				},
			},
			Type: pulumi.String("ClusterIP"),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create rpc service: %w", err)
	}
	component.RPCService = rpcService

	// Create StatefulSet
	statefulSetName := fmt.Sprintf("%s-set", args.Name)
	statefulSet, err := appsv1.NewStatefulSet(ctx, statefulSetName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(statefulSetName),
			Namespace: pulumi.String(args.Namespace),
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
					Labels: utils.CreateResourceLabels(args.Name, statefulSetName, args.Name, pulumi.StringMap{
						"app": pulumi.String(args.Name),
					}),
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String(args.Name),
							Image:           pulumi.String(args.Image),
							ImagePullPolicy: pulumi.String(args.ImagePullPolicy),
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
							Command: createExecutionClientCommand(args),
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									Name:          pulumi.String("p2p-tcp"),
									ContainerPort: pulumi.Int(args.P2PPort),
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("p2p-udp"),
									ContainerPort: pulumi.Int(args.P2PPort),
									Protocol:      pulumi.String("UDP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("rpc"),
									ContainerPort: pulumi.Int(args.RPCPort),
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("ws"),
									ContainerPort: pulumi.Int(args.WSPort),
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("metrics"),
									ContainerPort: pulumi.Int(args.MetricsPort),
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("auth-rpc"),
									ContainerPort: pulumi.Int(args.AuthRPCPort),
									Protocol:      pulumi.String("TCP"),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									Name:      pulumi.String("data"),
									MountPath: pulumi.String("/data"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.String("jwt"),
									MountPath: pulumi.String("/etc/execution/jwt"),
									ReadOnly:  pulumi.Bool(true),
								},
							},
						},
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
					NodeSelector: pulumi.StringMap{
						"kubernetes.io/os": pulumi.String("linux"),
					},
					Tolerations: corev1.TolerationArray{
						corev1.TolerationArgs{
							Key:      pulumi.String("node-role.kubernetes.io/control-plane"),
							Operator: pulumi.String("Exists"),
							Effect:   pulumi.String("NoSchedule"),
						},
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create statefulset: %w", err)
	}
	component.StatefulSet = statefulSet

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
