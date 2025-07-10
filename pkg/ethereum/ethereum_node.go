package ethereum

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewEthereumNodeComponent creates a new Ethereum node component that combines an execution client and a consensus client
func NewEthereumNodeComponent(ctx *pulumi.Context, args *EthereumNodeArgs, opts ...pulumi.ResourceOption) (*EthereumNodeComponent, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ethereum node args: %w", err)
	}

	component := &EthereumNodeComponent{
		Name:      args.Name,
		Namespace: args.Namespace,
	}

	// Create the execution client
	execClient, err := execution.NewExecutionClient(ctx, args.ExecutionClient, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client: %w", err)
	}
	component.ExecutionClient = execClient

	// Create the consensus client
	// Set the execution client endpoint in the consensus client args
	args.ConsensusClient.ExecutionClientEndpoint = fmt.Sprintf("http://%s-rpc.%s.svc.cluster.local:%d",
		args.Name, args.Namespace, args.ExecutionClient.RPCPort)

	consClient, err := consensus.NewConsensusClient(ctx, args.ConsensusClient, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus client: %w", err)
	}
	component.ConsensusClient = consClient

	return component, nil
}
