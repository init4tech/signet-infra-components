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

	component := &ConsensusClientComponent{}

	var name string
	pulumi.All(args.Name).ApplyT(func(values []interface{}) error {
		name = values[0].(string)
		return nil
	})

	err := ctx.RegisterComponentResource("signet:consensus:ConsensusClient", name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create PVC for data storage
	pvcName := fmt.Sprintf("%s-data", name)
	component.PVC, err = corev1.NewPersistentVolumeClaim(ctx, pvcName, &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(pvcName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(name, pvcName, name, nil),
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{
				pulumi.String("ReadWriteOnce"),
			},
			Resources: &corev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": args.StorageSize,
				},
			},
			StorageClassName: args.StorageClass,
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create PVC: %w", err)
	}

	// Create JWT secret
	jwtSecretName := fmt.Sprintf("%s-jwt", name)
	component.JWTSecret, err = corev1.NewSecret(ctx, jwtSecretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": args.JWTSecret,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(jwtSecretName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(name, jwtSecretName, name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT secret: %w", err)
	}

	// Create P2P service
	p2pServiceName := fmt.Sprintf("%s-p2p", name)
	component.P2PService, err = corev1.NewService(ctx, p2pServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(p2pServiceName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(name, p2pServiceName, name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("p2p"),
					Port:       args.P2PPort,
					TargetPort: args.P2PPort,
					Protocol:   pulumi.String("TCP"),
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create P2P service: %w", err)
	}

	// Create Beacon API service
	beaconAPIServiceName := fmt.Sprintf("%s-beacon-api", name)
	component.BeaconAPIService, err = corev1.NewService(ctx, beaconAPIServiceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(beaconAPIServiceName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(name, beaconAPIServiceName, name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(name),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("beacon-api"),
					Port:       args.BeaconAPIPort,
					TargetPort: args.BeaconAPIPort,
				},
				corev1.ServicePortArgs{
					Name:       pulumi.String("metrics"),
					Port:       args.MetricsPort,
					TargetPort: args.MetricsPort,
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create Beacon API service: %w", err)
	}

	// Create StatefulSet
	statefulSetName := name
	component.StatefulSet, err = appsv1.NewStatefulSet(ctx, statefulSetName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(statefulSetName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(name, statefulSetName, name, nil),
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(name),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String(name),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String("consensus"),
							Image:           args.Image,
							ImagePullPolicy: args.ImagePullPolicy,
							Command:         createConsensusClientCommand(args),
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									Name:          pulumi.String("p2p"),
									ContainerPort: args.P2PPort,
									Protocol:      pulumi.String("TCP"),
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("beacon-api"),
									ContainerPort: args.BeaconAPIPort,
								},
								corev1.ContainerPortArgs{
									Name:          pulumi.String("metrics"),
									ContainerPort: args.MetricsPort,
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

	component.Name = args.Name.ToStringOutput()
	component.Namespace = args.Namespace.ToStringOutput()
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

	// Add bootnodes
	pulumi.All(args.Bootnodes).ApplyT(func(nodes []interface{}) error {
		for _, bootnode := range nodes {
			cmd = append(cmd, pulumi.Sprintf("--boot-nodes=%s", bootnode.(string)))
		}
		return nil
	})

	// Add additional args
	pulumi.All(args.AdditionalArgs).ApplyT(func(args []interface{}) error {
		for _, arg := range args {
			cmd = append(cmd, pulumi.String(arg.(string)))
		}
		return nil
	})

	return cmd
}
