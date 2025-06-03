package consensus

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewConsensusClient creates a new consensus client component
func NewConsensusClient(ctx *pulumi.Context, args *ConsensusClientArgs, opts ...pulumi.ResourceOption) (*ConsensusClientComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid consensus client args: %w", err)
	}

	component := &ConsensusClientComponent{
		Name:      args.Name,
		Namespace: args.Namespace,
	}

	err := ctx.RegisterComponentResource("ethereum:index:ConsensusClient", args.Name, component, opts...)
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

	// Create Beacon API service
	beaconAPIServiceName := fmt.Sprintf("%s-beacon-api", args.Name)
	beaconAPIService, err := corev1.NewService(ctx, beaconAPIServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(beaconAPIServiceName),
			Namespace: pulumi.String(args.Namespace),
			Labels:    utils.CreateResourceLabels(args.Name, beaconAPIServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.BeaconAPIPort),
					Name:     pulumi.String("beacon-api"),
					Protocol: pulumi.String("TCP"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(args.MetricsPort),
					Name:     pulumi.String("metrics"),
					Protocol: pulumi.String("TCP"),
				},
			},
			Type: pulumi.String("ClusterIP"),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create beacon api service: %w", err)
	}
	component.BeaconAPIService = beaconAPIService

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
									"cpu":    pulumi.String("4"),
									"memory": pulumi.String("8Gi"),
								},
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String("2"),
									"memory": pulumi.String("4Gi"),
								},
							},
							Command: createConsensusClientCommand(args),
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
									Name:          pulumi.String("beacon-api"),
									ContainerPort: pulumi.Int(args.BeaconAPIPort),
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("metrics"),
									ContainerPort: pulumi.Int(args.MetricsPort),
									Protocol:      pulumi.String("TCP"),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									Name:      pulumi.String("data"),
									MountPath: pulumi.String("/data"),
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

// createConsensusClientCommand creates the command array for the consensus client
func createConsensusClientCommand(args *ConsensusClientArgs) pulumi.StringArray {
	cmd := pulumi.StringArray{
		pulumi.String("--datadir=/data"),
		pulumi.Sprintf("--execution-jwt=/etc/execution/jwt/jwt.hex"),
		pulumi.Sprintf("--execution-endpoint=%s", args.ExecutionClientEndpoint),
		pulumi.Sprintf("--port=%d", args.P2PPort),
		pulumi.Sprintf("--metrics-port=%d", args.MetricsPort),
		pulumi.Sprintf("--http-port=%d", args.BeaconAPIPort),
		pulumi.String("--http-address=0.0.0.0"),
		pulumi.String("--metrics-address=0.0.0.0"),
		pulumi.String("--validator-monitor-auto"),
		pulumi.String("--suggested-fee-recipient=0x0000000000000000000000000000000000000000"),
	}

	// Add bootnodes if provided
	for _, bootnode := range args.Bootnodes {
		cmd = append(cmd, pulumi.Sprintf("--boot-nodes=%s", bootnode))
	}

	// Add additional args if provided
	for _, arg := range args.AdditionalArgs {
		cmd = append(cmd, pulumi.String(arg))
	}

	return cmd
}
