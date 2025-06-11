package consensus

import (
	"fmt"
)

// Validate validates the consensus client arguments
func (args *ConsensusClientArgs) Validate() error {
	if args.Name == nil {
		return fmt.Errorf("name is required")
	}
	if args.Namespace == nil {
		return fmt.Errorf("namespace is required")
	}
	if args.StorageSize == nil {
		return fmt.Errorf("storageSize is required")
	}
	if args.StorageClass == nil {
		return fmt.Errorf("storageClass is required")
	}
	if args.Image == nil {
		return fmt.Errorf("image is required")
	}
	if args.ImagePullPolicy == nil {
		return fmt.Errorf("imagePullPolicy is required")
	}
	if args.JWTSecret == nil {
		return fmt.Errorf("jwtSecret is required")
	}
	if args.P2PPort == nil {
		return fmt.Errorf("p2pPort is required")
	}
	if args.BeaconAPIPort == nil {
		return fmt.Errorf("beaconAPIPort is required")
	}
	if args.MetricsPort == nil {
		return fmt.Errorf("metricsPort is required")
	}
	if args.ExecutionClientEndpoint == nil {
		return fmt.Errorf("executionClientEndpoint is required")
	}
	return nil
}
