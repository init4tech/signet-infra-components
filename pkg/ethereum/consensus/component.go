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

	// Convert public args to internal args for use with Pulumi
	internalArgs := args.toInternal()

	component := &ConsensusClientComponent{
		Name:      args.Name,
		Namespace: args.Namespace,
	}

	err := ctx.RegisterComponentResource("signet:consensus:ConsensusClient", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create PVC for data storage
	pvcName := fmt.Sprintf("%s-data", args.Name)
	component.PVC, err = corev1.NewPersistentVolumeClaim(ctx, pvcName, &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
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

	pvcNameInput := component.PVC.Metadata.Name().ApplyT(func(name *string) pulumi.StringInput {
		return pulumi.String(*name)
	}).(pulumi.StringInput)

	// Create JWT secret
	jwtSecretName := fmt.Sprintf("%s-jwt", args.Name)
	component.JWTSecret, err = corev1.NewSecret(ctx, jwtSecretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": internalArgs.JWTSecret,
		},
		Metadata: &metav1.ObjectMetaArgs{
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
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create P2P service: %w", err)
	}

	// Create Beacon API service
	beaconAPIServiceName := fmt.Sprintf("%s-beacon-api", args.Name)
	component.BeaconAPIService, err = corev1.NewService(ctx, beaconAPIServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, beaconAPIServiceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(args.Name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("beacon-api"),
					Port:       internalArgs.BeaconAPIPort,
					TargetPort: internalArgs.BeaconAPIPort,
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("metrics"),
					Port:       internalArgs.MetricsPort,
					TargetPort: internalArgs.MetricsPort,
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create Beacon API service: %w", err)
	}

	// Create StatefulSet
	statefulSetName := args.Name
	component.StatefulSet, err = appsv1.NewStatefulSet(ctx, statefulSetName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
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
						corev1.ContainerArgs{
							Name:            pulumi.String("consensus"),
							Image:           internalArgs.Image,
							ImagePullPolicy: internalArgs.ImagePullPolicy,
							Command:         createConsensusClientCommand(args),
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									Name:          pulumi.String("p2p"),
									ContainerPort: internalArgs.P2PPort,
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("beacon-api"),
									ContainerPort: internalArgs.BeaconAPIPort,
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("metrics"),
									ContainerPort: internalArgs.MetricsPort,
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
								},
							},
							Resources: nil,
						},
					},
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.String("data"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: pvcNameInput,
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("jwt"),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: component.JWTSecret.Metadata.Name(),
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

// createConsensusClientCommand creates the command array for the consensus client
func createConsensusClientCommand(args *ConsensusClientArgs) pulumi.StringArray {
	cmd := pulumi.StringArray{
		pulumi.String("lighthouse"),
		pulumi.String("bn"),
		pulumi.String("--datadir=/data/lighthouse"),
		pulumi.String("--http"),
		pulumi.Sprintf("--http-port=%d", args.BeaconAPIPort),
		pulumi.String("--http-address=0.0.0.0"),
		pulumi.Sprintf("--execution-jwt=/etc/execution/jwt/jwt.hex"),
		pulumi.Sprintf("--execution-endpoint=%s", args.ExecutionClientEndpoint),
		pulumi.Sprintf("--port=%d", args.P2PPort),
		pulumi.String("--metrics"),
		pulumi.Sprintf("--metrics-port=%d", args.MetricsPort),
		pulumi.String("--metrics-address=0.0.0.0"),
		pulumi.String("--validator-monitor-auto"),
		pulumi.String("--suggested-fee-recipient=0x0000000000000000000000000000000000000000"),
		pulumi.String("--checkpoint-sync-url=https://mainnet.checkpoint.sigp.io"),
	}

	// Add bootnodes
	if args.Bootnodes != nil {
		for _, bootnode := range args.Bootnodes {
			cmd = append(cmd, pulumi.Sprintf("--boot-nodes=%s", bootnode))
		}
	}

	// Add additional args
	if args.AdditionalArgs != nil {
		for _, arg := range args.AdditionalArgs {
			cmd = append(cmd, pulumi.String(arg))
		}
	}

	return cmd
}
