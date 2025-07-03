package pylon

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/ethereum"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
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
			Name:               pulumi.String(args.Name),
			Namespace:          pulumi.String(namespace),
			StorageSize:        pulumi.String(ExecutionClientStorageSize),
			StorageClass:       pulumi.String(StorageClassAWSGP3),
			Image:              args.PylonImage,
			JWTSecret:          args.ExecutionJwt,
			P2PPort:            pulumi.Int(ExecutionP2PPort),
			RPCPort:            pulumi.Int(ExecutionRPCPort),
			WSPort:             pulumi.Int(ExecutionWSPort),
			MetricsPort:        pulumi.Int(ExecutionMetricsPort),
			AuthRPCPort:        pulumi.Int(ExecutionAuthRPCPort),
			DiscoveryPort:      pulumi.Int(ExecutionP2PPort),
			ExecutionClientEnv: args.Env,
		},
		ConsensusClient: &consensus.ConsensusClientArgs{
			Name:            pulumi.String(args.Name),
			Namespace:       pulumi.String(namespace),
			StorageSize:     pulumi.String(ConsensusClientStorageSize),
			StorageClass:    pulumi.String(StorageClassAWSGP3),
			Image:           pulumi.String(ConsensusClientImage),
			ImagePullPolicy: pulumi.String(ImagePullPolicyAlways),
			JWTSecret:       args.ExecutionJwt,
			P2PPort:         pulumi.Int(ConsensusP2PPort),
			BeaconAPIPort:   pulumi.Int(ConsensusBeaconAPIPort),
			MetricsPort:     pulumi.Int(ConsensusMetricsPort),
		},
	}

	ethereumNode, err := ethereum.NewEthereumNodeComponent(ctx, args.Name, ethereumNodeArgs, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.EthereumNode = ethereumNode

	return component, nil
}
