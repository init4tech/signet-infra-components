package signet_node

import (
	"fmt"
)

// ApplyDefaults sets default values for optional fields
func (args *SignetNodeComponentArgs) ApplyDefaults() {
	if args.SignetNodeDataMountPath == "" {
		args.SignetNodeDataMountPath = DefaultSignetNodeDataMountPath
	}
	if args.RollupDataMountPath == "" {
		args.RollupDataMountPath = DefaultRollupDataMountPath
	}
	if args.ExecutionJwtMountPath == "" {
		args.ExecutionJwtMountPath = DefaultExecutionJwtMountPath
	}
}

func (args *SignetNodeComponentArgs) Validate() error {
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if args.ExecutionJwt == "" {
		return fmt.Errorf("execution jwt is required")
	}
	if args.ExecutionPvcSize == "" {
		return fmt.Errorf("execution pvc size is required")
	}
	if args.LighthousePvcSize == "" {
		return fmt.Errorf("lighthouse pvc size is required")
	}
	if args.RollupPvcSize == "" {
		return fmt.Errorf("rollup pvc size is required")
	}
	if args.ExecutionClientImage == "" {
		return fmt.Errorf("execution client image is required")
	}
	if args.ConsensusClientImage == "" {
		return fmt.Errorf("consensus client image is required")
	}
	if args.ExecutionClientStartCommand == nil {
		return fmt.Errorf("execution client start command is required")
	}
	if args.ConsensusClientStartCommand == nil {
		return fmt.Errorf("consensus client start command is required")
	}
	// Note: AppLabels is optional and has a default zero value (empty map is fine)

	if err := args.Env.Validate(); err != nil {
		return fmt.Errorf("invalid signet node env: %w", err)
	}
	return nil
}

func (env *SignetNodeEnv) Validate() error {
	if env.ChainName == "" {
		return fmt.Errorf("chainName is required")
	}
	if env.IpcEndpoint == "" {
		return fmt.Errorf("ipcEndpoint is required")
	}
	if env.RpcPort <= 0 {
		return fmt.Errorf("rpcPort must be a positive integer")
	}
	// RustLog is optional (*string), so we don't validate emptiness
	if env.SignetChainId <= 0 {
		return fmt.Errorf("signetChainId must be a positive integer")
	}
	if env.SignetClUrl == "" {
		return fmt.Errorf("signetClUrl is required")
	}
	if env.SignetDatabasePath == "" {
		return fmt.Errorf("signetDatabasePath is required")
	}
	if env.SignetPylonUrl == "" {
		return fmt.Errorf("signetPylonUrl is required")
	}
	if env.SignetStaticPath == "" {
		return fmt.Errorf("signetStaticPath is required")
	}
	if env.TxForwardUrl == "" {
		return fmt.Errorf("txForwardUrl is required")
	}
	if env.WsRpcPort <= 0 {
		return fmt.Errorf("wsRpcPort must be a positive integer")
	}
	return nil
}
