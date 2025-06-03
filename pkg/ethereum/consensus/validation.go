package consensus

import (
	"fmt"
)

// Validate validates the consensus client arguments
func (args *ConsensusClientArgs) Validate() error {
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if args.StorageSize == "" {
		return fmt.Errorf("storageSize is required")
	}
	if args.StorageClass == "" {
		return fmt.Errorf("storageClass is required")
	}
	if args.Image == "" {
		return fmt.Errorf("image is required")
	}
	if args.ExecutionClientEndpoint == "" {
		return fmt.Errorf("executionClientEndpoint is required")
	}
	if args.P2PPort <= 0 {
		return fmt.Errorf("p2pPort must be greater than 0")
	}
	if args.MetricsPort <= 0 {
		return fmt.Errorf("metricsPort must be greater than 0")
	}
	if args.BeaconAPIPort <= 0 {
		return fmt.Errorf("beaconApiPort must be greater than 0")
	}
	return nil
}
