package execution

import (
	"fmt"
)

// Validate validates the execution client arguments
// NOTE: This only works for pulumi.Int (not computed outputs) in tests.
func (args *ExecutionClientArgs) Validate() error {
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
	if args.ImagePullPolicy == "" {
		return fmt.Errorf("imagePullPolicy is required")
	}
	if args.JWTSecret == "" {
		return fmt.Errorf("jwtSecret is required")
	}
	if args.P2PPort == 0 {
		return fmt.Errorf("p2pPort is required")
	}
	if args.RPCPort == 0 {
		return fmt.Errorf("rpcPort is required")
	}
	if args.WSPort == 0 {
		return fmt.Errorf("wsPort is required")
	}
	if args.MetricsPort == 0 {
		return fmt.Errorf("metricsPort is required")
	}
	if args.AuthRPCPort == 0 {
		return fmt.Errorf("authRpcPort is required")
	}
	if args.DiscoveryPort == 0 {
		return fmt.Errorf("discoveryPort is required")
	}

	if args.P2PPort <= 0 {
		return fmt.Errorf("p2pPort must be greater than zero")
	}
	if args.RPCPort <= 0 {
		return fmt.Errorf("rpcPort must be greater than zero")
	}
	if args.WSPort <= 0 {
		return fmt.Errorf("wsPort must be greater than zero")
	}
	if args.MetricsPort <= 0 {
		return fmt.Errorf("metricsPort must be greater than zero")
	}
	if args.AuthRPCPort <= 0 {
		return fmt.Errorf("authRpcPort must be greater than zero")
	}
	if args.DiscoveryPort <= 0 {
		return fmt.Errorf("discoveryPort must be greater than zero")
	}
	return nil
}
