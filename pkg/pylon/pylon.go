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
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pylon component args: %w", err)
	}

	component := &PylonComponent{}

	err := ctx.RegisterComponentResource("signet:index:Pylon", args.Name, component, opts...)
	if err != nil {
		return nil, err
	}

	// Convert public args to internal args
	internalArgs := args.toInternal()

	stack := ctx.Stack()

	// Get the existing Route53 hosted zone for signet.sh
	dbProjectName := internalArgs.DbProjectName
	dbStackName := fmt.Sprintf("%s/%s", dbProjectName, stack)

	// TODO: this should be a stack reference to the pylon db stack- should this be an arg?
	// Need to think about how i want to handle the separation of the pylon db and pylon components
	thePylonDbStack, err := pulumi.NewStackReference(ctx, dbStackName, nil)
	if err != nil {
		return nil, err
	}

	// Get the database cluster endpoint (unused for now but needed for future implementation)
	_ = thePylonDbStack.GetStringOutput(pulumi.String("dbClusterEndpoint"))

	// Create the S3 bucket for blob storage (unused for now but needed for future implementation)
	_, err = s3.NewBucketV2(ctx, "pylon-blob-bucket", &s3.BucketV2Args{
		Bucket: internalArgs.PylonBlobBucketName,
	})
	if err != nil {
		return nil, err
	}

	// Convert environment to internal type for use with ethereum components
	internalEnv := args.Env.toInternal()

	// Create Ethereum node component
	ethereumNodeArgs := &ethereum.EthereumNodeArgs{
		Name:      args.Name,
		Namespace: args.Namespace,
		ExecutionClient: &execution.ExecutionClientArgs{
			Name:               args.Name,
			Namespace:          args.Namespace,
			StorageSize:        ExecutionClientStorageSize,
			StorageClass:       StorageClassAWSGP3,
			Image:              args.PylonImage,
			JWTSecret:          args.ExecutionJwt,
			P2PPort:            ExecutionP2PPort,
			RPCPort:            ExecutionRPCPort,
			WSPort:             ExecutionWSPort,
			MetricsPort:        ExecutionMetricsPort,
			AuthRPCPort:        ExecutionAuthRPCPort,
			DiscoveryPort:      ExecutionP2PPort,
			ExecutionClientEnv: internalEnv,
		},
		ConsensusClient: &consensus.ConsensusClientArgs{
			Name:            args.Name,
			Namespace:       args.Namespace,
			StorageSize:     ConsensusClientStorageSize,
			StorageClass:    StorageClassAWSGP3,
			Image:           ConsensusClientImage,
			ImagePullPolicy: ImagePullPolicyAlways,
			BeaconAPIPort:   ConsensusBeaconAPIPort,
			MetricsPort:     ConsensusMetricsPort,
		},
	}

	ethereumNode, err := ethereum.NewEthereumNodeComponent(ctx, ethereumNodeArgs, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.EthereumNode = ethereumNode

	return component, nil
}
