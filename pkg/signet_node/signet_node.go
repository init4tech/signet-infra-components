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

	component := &SignetNodeComponent{
		SignetNodeComponentArgs: args,
	}
	err := ctx.RegisterComponentResource("signet:index:SignetNode", args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	storageSize := pulumi.String("150Gi")

	_, err = CreatePersistentVolumeClaim(
		ctx,
		"signet-node-data",
		args.Namespace,
		storageSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signet node db data pvc: %w", err)
	}

	_, err = CreatePersistentVolumeClaim(
		ctx,
		"rollup-data",
		args.Namespace,
		storageSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rollup data pvc: %w", err)
	}

	secretName := "execution-jwt"
	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
		StringData: pulumi.StringMap{
			"jwt.hex": args.ExecutionJwt,
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(secretName),
			Namespace: args.Namespace,
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
		args.Namespace,
		utils.CreateResourceLabels(args.Name, executionConfigMapName, args.Name, nil),
		args.Env,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution configmap: %w", err)
	}
	component.SignetNodeConfigMap = executionConfigMap

	// SERVICE
	executionClientName := "signet-node"
	executionClientServiceName := fmt.Sprintf("%s-service", executionClientName)

	hostExecutionClientService, err := corev1.NewService(ctx, executionClientServiceName, &corev1.ServiceArgs{
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
			Name:      pulumi.String(executionClientServiceName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, executionClientServiceName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client service: %w", err)
	}

	executionServiceIpString := hostExecutionClientService.Spec.ClusterIP().Elem()

	// STATEFUL SET
	hostStatefulSetName := "signet-node-execution-set"
	hostStatefulSetResourceName := fmt.Sprintf("%s-set", hostStatefulSetName)

	// Create pod labels with app label for stateful set
	executionPodLabels := utils.CreateResourceLabels(args.Name, hostStatefulSetName, args.Name, nil)
	executionPodLabels["app"] = pulumi.String(hostStatefulSetName)

	// Define the StatefulSet for the 'reth' container with a configmap volume and a data persistent volume
	_, err = appsv1.NewStatefulSet(ctx, hostStatefulSetResourceName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(hostStatefulSetName),
			Labels:    utils.CreateResourceLabels(args.Name, hostStatefulSetResourceName, args.Name, nil),
			Namespace: args.Namespace,
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
					Namespace: args.Namespace,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String(hostStatefulSetName),
							Image:           args.ExecutionClientImage,
							ImagePullPolicy: pulumi.String("Always"),
							Command: pulumi.StringArray{
								pulumi.String("signet"),
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
								ClaimName: pulumi.String("signet-node-data"),
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
								ClaimName: pulumi.String("rollup-data"),
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

	_, err = CreatePersistentVolumeClaim(
		ctx,
		"real-lighthouse-data",
		args.Namespace,
		storageSize,
		"aws-gp3",
		component,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse data pvc: %w", err)
	}

	// Create ConfigMap for consensus environment variables
	consensusEnv := ConsensusEnv{
		Example: pulumi.String("example"),
	}

	consensusConfigMapName := "consensus-configmap-env-config"
	consensusConfigMap, err := utils.CreateConfigMap(
		ctx,
		consensusConfigMapName,
		args.Namespace,
		utils.CreateResourceLabels(args.Name, consensusConfigMapName, args.Name, nil),
		consensusEnv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus configmap: %w", err)
	}
	component.LighthouseConfigMap = consensusConfigMap

	lighthouseServiceName := fmt.Sprintf("%s-service", consensusClientName)
	lighthouseInternalService, err := corev1.NewService(ctx, lighthouseServiceName, &corev1.ServiceArgs{
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
			Name:      pulumi.String(lighthouseServiceName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, lighthouseServiceName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse internal service: %w", err)
	}

	lighthouseServiceIpString := lighthouseInternalService.Spec.ClusterIP().Elem()

	lighthouseStatefulSet := "lighthouse"
	lighthouseStatefulSetResourceName := fmt.Sprintf("%s-set", lighthouseStatefulSet)

	// Create pod labels with app label for stateful set
	lighthousePodLabels := utils.CreateResourceLabels(args.Name, lighthouseStatefulSet, args.Name, nil)
	lighthousePodLabels["app"] = pulumi.String(lighthouseStatefulSet)

	_, err = appsv1.NewStatefulSet(ctx, lighthouseStatefulSetResourceName, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(lighthouseStatefulSet),
			Labels:    utils.CreateResourceLabels(args.Name, lighthouseStatefulSetResourceName, args.Name, nil),
			Namespace: args.Namespace,
		},
		Spec: &appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(lighthouseStatefulSet),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Namespace: args.Namespace,
					Labels:    lighthousePodLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String(lighthouseStatefulSet),
							Image: args.ConsensusClientImage,
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
									Name:      pulumi.Sprintf("%s-config", lighthouseStatefulSet),
									MountPath: pulumi.String("/etc/lighthouse"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.String("real-lighthouse-db"),
									MountPath: pulumi.String("/root/.lighthouse/holesky"),
								},
								corev1.VolumeMountArgs{
									Name:      pulumi.Sprintf("%s-execution-jwt", lighthouseStatefulSet),
									MountPath: pulumi.String("/secrets"),
								},
							},
							Resources: NewResourceRequirements("2", "16Gi", "2", "4Gi"),
						},
					},
					DnsPolicy: pulumi.String("ClusterFirst"),
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-config", lighthouseStatefulSet),
							ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
								Name: consensusConfigMap.Metadata.Name(),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.String("real-lighthouse-db"),
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: pulumi.String("real-lighthouse-data"),
							},
						},
						corev1.VolumeArgs{
							Name: pulumi.Sprintf("%s-execution-jwt", lighthouseStatefulSet),
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
	_, err = crd.NewCustomResource(ctx, "signet-rpc-vservice", &crd.CustomResourceArgs{
		ApiVersion: pulumi.String("networking.istio.io/v1alpha3"),
		Kind:       pulumi.String("VirtualService"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(virtualServiceName),
			Namespace: args.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, virtualServiceName, args.Name, nil),
		},
		OtherFields: map[string]interface{}{
			"spec": map[string]interface{}{
				"hosts": []string{
					"rpc.pecorino.signet.sh",
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
									"host": "signet-node-service.kt-cloud.svc.cluster.local",
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

	return component, nil
}
