package pylon

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewPylonComponent(ctx *pulumi.Context, args *PylonComponentArgs, opts ...pulumi.ResourceOption) (*PylonComponent, error) {
	component := &PylonComponent{}

	err := ctx.RegisterComponentResource("signet:index:Pylon", args.Name, component, opts...)
	if err != nil {
		return nil, err
	}

	namespace := args.Namespace
	stack := ctx.Stack()

	// Get the existing Route53 hosted zone for signet.sh
	dbProjectName := args.DbProjectName
	dbStackName := fmt.Sprintf("%s/%s", dbProjectName, stack)

	// TODO: this should be a stack reference to the pylon db stack- should this be an arg?
	// Need to think about how i want to handle the separation of the pylon db and pylon components
	thePylonDbStack, err := pulumi.NewStackReference(ctx, dbStackName, nil)
	if err != nil {
		return nil, err
	}

	dbClusterEndpoint := thePylonDbStack.GetStringOutput(pulumi.String("dbClusterEndpoint"))

	pylonBlobBucket, err := s3.NewBucketV2(ctx, "pylon-blob-bucket", &s3.BucketV2Args{
		Bucket: pulumi.Sprintf(args.PylonBlobBucketName),
	})
	if err != nil {
		return nil, err
	}

	// Define node names for both pairs
	executionClientName := fmt.Sprintf("%s-pylon", args.Name)
	consensusClientName := fmt.Sprintf("%s-pylon-cl", args.Name)

	pylonNodeEnv := pulumi.StringMap{
		"PYLON_START_BLOCK":             args.Env.PylonStartBlock,
		"PYLON_SENDERS":                 args.Env.PylonSenderAddress,
		"PYLON_DB_URL":                  pulumi.Sprintf("postgresql://%s:%s@%s:%s/pylon", args.Env.PostgresUser, args.Env.PostgresPassword, dbClusterEndpoint, "5432"),
		"PYLON_S3_URL":                  args.Env.PylonS3Url,
		"PYLON_S3_REGION":               args.Env.PylonS3Region,
		"PYLON_S3_BUCKET_NAME":          pylonBlobBucket.Bucket,
		"PYLON_CL_URL":                  args.Env.PylonConsensusClientUrl,
		"PYLON_BLOBSCAN_BASE_URL":       args.Env.PylonBlobscanBaseUrl,
		"PYLON_NETWORK_START_TIMESTAMP": args.Env.PylonNetworkStartTimestamp,
		"PYLON_NETWORK_SLOT_DURATION":   args.Env.PylonNetworkSlotDuration,
		"PYLON_NETWORK_SLOT_OFFSET":     args.Env.PylonNetworkSlotOffset,
		"PYLON_REQUESTS_PER_SECOND":     args.Env.PylonRequestsPerSecond,
		"PYLON_PORT":                    args.Env.PylonPort,
		"RUST_LOG":                      args.Env.PylonRustLog,
		"AWS_ACCESS_KEY_ID":             args.Env.AwsAccessKeyId,
		"AWS_SECRET_ACCESS_KEY":         args.Env.AwsSecretAccessKey,
		"AWS_DEFAULT_REGION":            args.Env.AwsRegion,
	}

	storageSize := pulumi.String("150Gi")

	// create a snapshot of the existing execution client

	_, err = corev1.NewPersistentVolumeClaim(ctx, fmt.Sprintf("%s-data", executionClientName), &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-data", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-data", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{pulumi.String("ReadWriteOnce")},
			Resources: &corev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": storageSize,
				},
			},
			StorageClassName: pulumi.String("aws-gp3"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a secret for the execution jwt
	secret, err := corev1.NewSecret(ctx, fmt.Sprintf("%s-execution-jwt", executionClientName), &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": args.ExecutionJwt,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-execution-jwt", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name": pulumi.Sprintf("%s-execution-jwt", executionClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	pylonEnvConfigMap, err := corev1.NewConfigMap(ctx, fmt.Sprintf("%s-env-configmap", executionClientName), &corev1.ConfigMapArgs{
		Data: pylonNodeEnv,
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-env-configmap", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-env-configmap", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a Service for external ports
	_, err = corev1.NewService(ctx, fmt.Sprintf("%s-p2pnet-service", executionClientName), &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.Sprintf("%s", executionClientName)},
			Type:     pulumi.String("NodePort"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port: pulumi.Int(30303),
					Name: pulumi.String("p2p-tcp"),
				},
				&corev1.ServicePortArgs{
					Port:     pulumi.Int(30303),
					Protocol: pulumi.String("UDP"),
					Name:     pulumi.String("p2p-udp"),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-p2pnet-service", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-p2pnet-service", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a service for internal ports
	executionInternalService, err := corev1.NewService(ctx, fmt.Sprintf("%s-internal-service", executionClientName), &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.Sprintf("%s", executionClientName)},
			Type:     pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port: pulumi.Int(9001),
					Name: pulumi.String("metrics"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(8551),
					Name: pulumi.String("p2p"),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-internal-service", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-internal-service", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create ingress for the reth rpc traffic on port 8545
	_, err = corev1.NewService(ctx, fmt.Sprintf("%s-rpc-service", executionClientName), &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.Sprintf("%s", executionClientName)},
			Type:     pulumi.String("NodePort"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:       pulumi.Int(8545),
					TargetPort: pulumi.Int(8545),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-rpc-service", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-rpc-service", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = corev1.NewPersistentVolumeClaim(ctx, fmt.Sprintf("%s-data", consensusClientName), &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-data", consensusClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-data", consensusClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", consensusClientName),
			},
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{pulumi.String("ReadWriteOnce")}, // This should match your requirements
			Resources: &corev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": storageSize,
				},
			},
			StorageClassName: pulumi.String("aws-gp3"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a secret for the execution jwt
	secret, err = corev1.NewSecret(ctx, fmt.Sprintf("%s-execution-jwt", consensusClientName), &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": args.ExecutionJwt,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-execution-jwt", consensusClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name": pulumi.Sprintf("%s-execution-jwt", consensusClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create an internal service for lighthouse validator client to connect to the beacon node
	lighthouseInternalService, err := corev1.NewService(ctx, fmt.Sprintf("%s-internal-service", consensusClientName), &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-internal-service", consensusClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-internal-service", consensusClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", consensusClientName),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.Sprintf("%s", consensusClientName)},
			Type:     pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port: pulumi.Int(8545),
					Name: pulumi.String("beacon-node"),
				},
				corev1.ServicePortArgs{
					Port: pulumi.Int(9000),
					Name: pulumi.String("p2p-tcp"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(9000),
					Protocol: pulumi.String("UDP"),
					Name:     pulumi.String("p2p-udp"),
				},
				corev1.ServicePortArgs{
					Port:     pulumi.Int(9001),
					Protocol: pulumi.String("UDP"),
					Name:     pulumi.String("quic"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	executionServiceIpString := executionInternalService.Spec.ClusterIP().Elem()
	lighthouseServiceIpString := lighthouseInternalService.Spec.ClusterIP().Elem()

	// create the metrics service
	_, err = corev1.NewService(ctx, fmt.Sprintf("%s-metrics-service", consensusClientName), &corev1.ServiceArgs{
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{"app": pulumi.Sprintf("%s", consensusClientName)},
			Type:     pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{

				corev1.ServicePortArgs{
					Port: pulumi.Int(5054),
					Name: pulumi.String("metrics"),
				},
			},
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-metrics-service", consensusClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-metrics-service", consensusClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", consensusClientName),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Define the StatefulSet for the 'reth' container with a configmap volume and a data persistent volume
	_, err = appsv1.NewStatefulSet(ctx, fmt.Sprintf("%s-set", executionClientName), &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s", executionClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app":                       pulumi.Sprintf("%s-set", executionClientName),
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-set", executionClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
			},
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.Sprintf("%s", executionClientName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app":                       pulumi.Sprintf("%s", executionClientName),
						"app.kubernetes.io/name":    pulumi.Sprintf("%s", executionClientName),
						"app.kubernetes.io/part-of": pulumi.Sprintf("%s", executionClientName),
					},
				},
				Spec: &corev1.PodSpecArgs{
					InitContainers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String("pylon-init"),
							Image: pulumi.String("637423570300.dkr.ecr.us-east-1.amazonaws.com/pylon-init:latest"),
							Command: pulumi.StringArray{
								pulumi.String("psql"),
								pulumi.String("-h"),
								dbClusterEndpoint,
								pulumi.String("-U"),
								args.Env.PostgresUser,
								pulumi.String("-d"),
								pulumi.String("pylon"),
								pulumi.String("-f"),
								pulumi.String("/init.sql"),
							},
							Env: corev1.EnvVarArray{
								corev1.EnvVarArgs{
									Name:  pulumi.String("PGPASSWORD"),
									Value: args.Env.PostgresPassword,
								},
							},
							ImagePullPolicy: pulumi.String("Always"),
						},
					},
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.Sprintf("%s", executionClientName),
							Image:           pulumi.String("637423570300.dkr.ecr.us-east-1.amazonaws.com/pylon:latest"),
							ImagePullPolicy: pulumi.String("Always"),
							Command: pulumi.StringArray{
								pulumi.String("pylon"),
								pulumi.String("node"),
								pulumi.String("--datadir=/root/.local/share/reth/pecorino"),
								pulumi.String("--chain=/network_configs/genesis.json"),
								pulumi.String("--http"),
								pulumi.String("--http.port=8545"),
								pulumi.String("--http.addr=0.0.0.0"),
								pulumi.String("--http.corsdomain=*"),
								pulumi.String("--http.api=admin,net,eth,web3,debug,txpool,trace"),
								pulumi.String("--ws"),
								pulumi.String("--ws.addr=0.0.0.0"),
								pulumi.String("--ws.port=8546"),
								pulumi.String("--ws.api=net,eth"),
								pulumi.String("--ws.origins=*"),
								executionServiceIpString.ApplyT(func(ip string) string {
									return "--nat=extip:" + ip
								}).(pulumi.StringOutput),
								pulumi.String("--authrpc.port=8551"),
								pulumi.String("--authrpc.jwtsecret=/etc/reth/execution-jwt/jwt.hex"),
								pulumi.String("--authrpc.addr=0.0.0.0"),
								pulumi.String("--metrics=0.0.0.0:9001"),
								pulumi.String("--discovery.port=30303"),
								pulumi.String("--port=30303"),
								pulumi.String("--bootnodes=enode://02b1442bab88aa1fefc8fff37bcbffa18abf91450850bfc3e2d7d29296144dc74e71400bacd7882061d5193990ccba5441fde0f69dc00d065b0237cd236a055d@10.100.1.188:30303"),
							},
							EnvFrom: corev1.EnvFromSourceArray{
								corev1.EnvFromSourceArgs{
									ConfigMapRef: &corev1.ConfigMapEnvSourceArgs{
										Name:     pulumi.Sprintf("%s-env-configmap", executionClientName),
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
									ContainerPort: pulumi.Int(8080),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									Name:      pulumi.Sprintf("%s-data", executionClientName),
									MountPath: pulumi.String("/root/.local/share/reth"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.Sprintf("%s-execution-jwt", executionClientName),
									MountPath: pulumi.String("/etc/reth/execution-jwt"),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"cpu":    pulumi.String("4"),
									"memory": pulumi.String("24Gi"),
								},
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String("2"),
									"memory": pulumi.String("16Gi"),
								},
							},
						},
					},
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-data", executionClientName),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: pulumi.Sprintf("%s-data", executionClientName),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-execution-jwt", executionClientName),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: secret.Metadata.Name(),
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{pylonEnvConfigMap}))
	if err != nil {
		return nil, err
	}

	spec := &corev1.PodSpecArgs{
		Containers: corev1.ContainerArray{
			corev1.ContainerArgs{
				Name:            pulumi.Sprintf("%s", consensusClientName),
				Image:           pulumi.String("637423570300.dkr.ecr.us-east-1.amazonaws.com/pecorino-lighthouse:latest"),
				ImagePullPolicy: pulumi.String("Always"),
				Command: pulumi.StringArray{
					pulumi.String("lighthouse"),
					pulumi.String("beacon_node"),
					pulumi.String("--debug-level=info"),
					pulumi.String("--datadir=/root/.lighthouse/holesky/beacon-data"),
					pulumi.String("--disable-enr-auto-update"),
					lighthouseServiceIpString.ApplyT(func(ip string) string {
						return "--enr-address=" + ip
					}).(pulumi.StringOutput),
					pulumi.String("--enr-udp-port=9000"),
					pulumi.String("--enr-tcp-port=9000"),
					pulumi.String("--listen-address=0.0.0.0"),
					pulumi.String("--port=9000"),
					pulumi.String("--http"),
					pulumi.String("--http-address=0.0.0.0"),
					pulumi.String("--http-port=4000"),
					pulumi.String("--slots-per-restore-point=32"),
					pulumi.String("--disable-packet-filter"),
					executionServiceIpString.ApplyT(func(ip string) string {
						return "--execution-endpoints=http://" + ip + ":8551"
					}).(pulumi.StringOutput),
					pulumi.String("--execution-jwt=/secrets/jwt.hex"),
					pulumi.String("--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776"),
					pulumi.String("--metrics"),
					pulumi.String("--metrics-address=0.0.0.0"),
					pulumi.String("--metrics-allow-origin=*"),
					pulumi.String("--metrics-port=5054"),
					pulumi.String("--enable-private-discovery"),
					pulumi.String("--testnet-dir=/network_configs"),
					pulumi.String("--boot-nodes=enr:-MS4QFGg1RsGSWdGQuBest1ST2az6Zea1-bJBoNsZKmbi5xVfPKtoA-YScWKSibelm35OQf8a8sAXZAwpMASChDM1FQBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpDYTi7gYAAAOADh9QUAAAAAgmlkgnY0gmlwhApkhyWEcXVpY4IjKYlzZWNwMjU2azGhArnO2f0iOAd09pC3Izvz2YU-oadRMuVgCe9jgOQEuG8qiHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo"),
					pulumi.String("--checkpoint-sync-url=http://cl-1-lighthouse-reth.kt-cloud.svc.cluster.local:4000"),
				},
				Ports: corev1.ContainerPortArray{
					corev1.ContainerPortArgs{
						ContainerPort: pulumi.Int(9000),
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
				},
				VolumeMounts: corev1.VolumeMountArray{
					corev1.VolumeMountArgs{
						Name:      pulumi.Sprintf("%s-data", consensusClientName),
						MountPath: pulumi.String("/root/.lighthouse/holesky"),
					},
					corev1.VolumeMountArgs{
						Name:      pulumi.Sprintf("%s-execution-jwt", consensusClientName),
						MountPath: pulumi.String("/secrets"),
					},
				},
				Resources: &corev1.ResourceRequirementsArgs{
					Limits: pulumi.StringMap{
						"cpu":    pulumi.String("4"),
						"memory": pulumi.String("16Gi"),
					},
					Requests: pulumi.StringMap{
						"cpu":    pulumi.String("2"),
						"memory": pulumi.String("12Gi"),
					},
				},
			},
		},
		DnsPolicy: pulumi.String("ClusterFirst"),
		Volumes: corev1.VolumeArray{
			corev1.VolumeArgs{
				Name: pulumi.Sprintf("%s-data", consensusClientName),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
					ClaimName: pulumi.Sprintf("%s-data", consensusClientName),
				},
			},
			corev1.VolumeArgs{
				Name: pulumi.Sprintf("%s-execution-jwt", consensusClientName),
				Secret: &corev1.SecretVolumeSourceArgs{
					SecretName: secret.Metadata.Name(),
				},
			},
		},
	}

	// Create a stateful set to run a lighthouse node with a configmap volume and a data persistent volume
	_, err = appsv1.NewStatefulSet(ctx, fmt.Sprintf("%s-set", consensusClientName), &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s", consensusClientName),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-set", consensusClientName),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", consensusClientName),
			},
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.Sprintf("%s", consensusClientName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app":                       pulumi.Sprintf("%s", consensusClientName),
						"app.kubernetes.io/name":    pulumi.Sprintf("%s", consensusClientName),
						"app.kubernetes.io/part-of": pulumi.Sprintf("%s", consensusClientName),
					},
				},
				Spec: spec,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a service to expose the pylon node
	_, err = corev1.NewService(ctx, "pylon-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("pylon-service"),
			Namespace: pulumi.String(namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String("pecorino-pylon-reth"),
			},
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port: pulumi.Int(8080),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return component, nil
}
