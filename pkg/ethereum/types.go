package ethereum

import (
	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Public-facing structs with base Go types

// EthereumNodeArgs contains the configuration for both execution and consensus clients
type EthereumNodeArgs struct {
	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace string
	// ExecutionClient contains the configuration for the execution client
	ExecutionClient *execution.ExecutionClientArgs
	// ConsensusClient contains the configuration for the consensus client
	ConsensusClient *consensus.ConsensusClientArgs
}

// Internal structs with Pulumi types

type ethereumNodeArgsInternal struct {
	// Name is the base name for all resources
	Name pulumi.StringInput
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace pulumi.StringInput
	// ExecutionClient contains the configuration for the execution client
	ExecutionClient *execution.ExecutionClientArgs
	// ConsensusClient contains the configuration for the consensus client
	ConsensusClient *consensus.ConsensusClientArgs
}

// Conversion functions

// toInternal converts public args to internal args for use with Pulumi
func (args EthereumNodeArgs) toInternal() ethereumNodeArgsInternal {
	return ethereumNodeArgsInternal{
		Name:            pulumi.String(args.Name),
		Namespace:       pulumi.String(args.Namespace),
		ExecutionClient: args.ExecutionClient,
		ConsensusClient: args.ConsensusClient,
	}
}

// EthereumNodeComponent represents a complete Ethereum node deployment
type EthereumNodeComponent struct {
	pulumi.ResourceState

	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace
	Namespace string
	// ExecutionClient is the execution client component
	ExecutionClient *execution.ExecutionClientComponent
	// ConsensusClient is the consensus client component
	ConsensusClient *consensus.ConsensusClientComponent
}
