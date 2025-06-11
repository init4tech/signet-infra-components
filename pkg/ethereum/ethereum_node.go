package ethereum

import (
	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewEthereumNodeComponent creates a new Ethereum node component that combines an execution client and a consensus client
func NewEthereumNodeComponent(ctx *pulumi.Context, name string, args *EthereumNodeArgs, opts ...pulumi.ResourceOption) (*EthereumNodeComponent, error) {
	component := &EthereumNodeComponent{
		Name: name,
	}

	// Create the execution client
	execClient, err := execution.NewExecutionClient(ctx, args.ExecutionClient, opts...)
	if err != nil {
		return nil, err
	}
	component.ExecutionClient = execClient

	// Create the consensus client
	// Set the execution client endpoint in the consensus client args
	args.ConsensusClient.ExecutionClientEndpoint = pulumi.Sprintf("http://%s-rpc.%s.svc.cluster.local:%d",
		name, args.Namespace, args.ExecutionClient.RPCPort)

	consClient, err := consensus.NewConsensusClient(ctx, args.ConsensusClient, opts...)
	if err != nil {
		return nil, err
	}
	component.ConsensusClient = consClient

	return component, nil
}
