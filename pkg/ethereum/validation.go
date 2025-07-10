package ethereum

import (
	"fmt"
)

// Validate validates the EthereumNodeArgs struct
func (args *EthereumNodeArgs) Validate() error {
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.ExecutionClient == nil {
		return fmt.Errorf("executionClient is required")
	}

	if args.ConsensusClient == nil {
		return fmt.Errorf("consensusClient is required")
	}

	// Validate execution client
	if err := args.ExecutionClient.Validate(); err != nil {
		return fmt.Errorf("execution client validation failed: %w", err)
	}

	// Validate consensus client
	if err := args.ConsensusClient.Validate(); err != nil {
		return fmt.Errorf("consensus client validation failed: %w", err)
	}

	return nil
}
