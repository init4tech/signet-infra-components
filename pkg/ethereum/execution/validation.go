package execution

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Validate validates the execution client arguments
// NOTE: This only works for pulumi.Int (not computed outputs) in tests.
func (args *ExecutionClientArgs) Validate() error {
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
	if args.JWTSecret == nil {
		return fmt.Errorf("jwtSecret is required")
	}
	if args.P2PPort == nil {
		return fmt.Errorf("p2pPort is required")
	}
	if args.RPCPort == nil {
		return fmt.Errorf("rpcPort is required")
	}
	if args.WSPort == nil {
		return fmt.Errorf("wsPort is required")
	}
	if args.MetricsPort == nil {
		return fmt.Errorf("metricsPort is required")
	}
	if args.AuthRPCPort == nil {
		return fmt.Errorf("authRpcPort is required")
	}
	if args.DiscoveryPort == nil {
		return fmt.Errorf("discoveryPort is required")
	}

	// Only works for pulumi.Int, not computed outputs
	getInt := func(input pulumi.IntInput) int {
		var v int
		input.ToIntOutput().ApplyT(func(i int) int {
			v = i
			return i
		})
		return v
	}

	p2pPort := getInt(args.P2PPort)
	rpcPort := getInt(args.RPCPort)
	wsPort := getInt(args.WSPort)
	metricsPort := getInt(args.MetricsPort)
	authRpcPort := getInt(args.AuthRPCPort)
	discoveryPort := getInt(args.DiscoveryPort)

	if p2pPort <= 0 {
		return fmt.Errorf("p2pPort must be greater than zero")
	}
	if rpcPort <= 0 {
		return fmt.Errorf("rpcPort must be greater than zero")
	}
	if wsPort <= 0 {
		return fmt.Errorf("wsPort must be greater than zero")
	}
	if metricsPort <= 0 {
		return fmt.Errorf("metricsPort must be greater than zero")
	}
	if authRpcPort <= 0 {
		return fmt.Errorf("authRpcPort must be greater than zero")
	}
	if discoveryPort <= 0 {
		return fmt.Errorf("discoveryPort must be greater than zero")
	}
	return nil
}
