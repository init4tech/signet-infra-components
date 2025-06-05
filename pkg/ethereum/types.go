package ethereum

import (
	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// EthereumNodeArgs contains the configuration for both execution and consensus clients
type EthereumNodeArgs struct {
	// Name is the base name for all resources
	Name pulumi.StringInput `pulumi:"name"`
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace pulumi.StringInput `pulumi:"namespace"`
	// ExecutionClient contains the configuration for the execution client
	ExecutionClient *execution.ExecutionClientArgs `pulumi:"executionClient"`
	// ConsensusClient contains the configuration for the consensus client
	ConsensusClient *consensus.ConsensusClientArgs `pulumi:"consensusClient"`
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
