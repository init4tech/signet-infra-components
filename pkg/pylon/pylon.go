package pylon

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/ethereum"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
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

	// Create Ethereum node component
	ethereumNodeArgs := &ethereum.EthereumNodeArgs{
		Name:      pulumi.String(args.Name),
		Namespace: pulumi.String(namespace),
		ExecutionClient: &execution.ExecutionClientArgs{
			Name:          pulumi.String(args.Name),
			Namespace:     pulumi.String(namespace),
			StorageSize:   pulumi.String("150Gi"),
			StorageClass:  pulumi.String("aws-gp3"),
			Image:         pulumi.String("ghcr.io/paradigmxyz/reth:latest"),
			JWTSecret:     args.ExecutionJwt,
			P2PPort:       pulumi.Int(30303),
			RPCPort:       pulumi.Int(8545),
			WSPort:        pulumi.Int(8546),
			MetricsPort:   pulumi.Int(9001),
			AuthRPCPort:   pulumi.Int(8551),
			DiscoveryPort: pulumi.Int(30303),
		},
		ConsensusClient: &consensus.ConsensusClientArgs{
			Name:            pulumi.String(args.Name),
			Namespace:       pulumi.String(namespace),
			StorageSize:     pulumi.String("100Gi"),
			StorageClass:    pulumi.String("aws-gp3"),
			Image:           pulumi.String("sigp/lighthouse:latest"),
			ImagePullPolicy: pulumi.String("Always"),
			JWTSecret:       args.ExecutionJwt,
			P2PPort:         pulumi.Int(9000),
			BeaconAPIPort:   pulumi.Int(4000),
			MetricsPort:     pulumi.Int(5054),
		},
	}

	ethereumNode, err := ethereum.NewEthereumNodeComponent(ctx, args.Name, ethereumNodeArgs, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.EthereumNode = ethereumNode

	// Create Pylon environment config map
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

	pylonEnvConfigMap, err := corev1.NewConfigMap(ctx, fmt.Sprintf("%s-env-configmap", args.Name), &corev1.ConfigMapArgs{
		Data: pylonNodeEnv,
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.Sprintf("%s-env-configmap", args.Name),
			Namespace: pulumi.String(namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.Sprintf("%s-env-configmap", args.Name),
				"app.kubernetes.io/part-of": pulumi.Sprintf("%s", args.Name),
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.PylonEnvConfigMap = pylonEnvConfigMap

	return component, nil
}
