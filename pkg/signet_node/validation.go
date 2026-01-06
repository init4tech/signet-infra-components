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
	// Apply defaults to env
	args.Env.ApplyDefaults()
}

// ApplyDefaults sets default values for optional SignetNodeEnv fields
func (env *SignetNodeEnv) ApplyDefaults() {
	if env.RpcPort == 0 {
		env.RpcPort = DefaultSignetRpcPort
	}
	if env.WsRpcPort == 0 {
		env.WsRpcPort = DefaultSignetWsRpcPort
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
	// Note: AppLabels is not validated as it's optional and has a default struct value
	if err := args.Env.Validate(); err != nil {
		return fmt.Errorf("invalid signet node env: %w", err)
	}
	return nil
}

func (env *SignetNodeEnv) Validate() error {
	// Always required fields
	if env.BlobExplorerUrl == "" {
		return fmt.Errorf("blob explorer url is required")
	}
	if env.SignetStaticPath == "" {
		return fmt.Errorf("signet static path is required")
	}
	if env.SignetDatabasePath == "" {
		return fmt.Errorf("signet database path is required")
	}

	// Conditional validation: if ChainName is not set, require genesis and slot calculator vars
	if env.ChainName == "" {
		// Genesis configuration required
		if env.RollupGenesisJsonPath == "" {
			return fmt.Errorf("rollup genesis json path is required when chain name is not set")
		}
		if env.HostGenesisJsonPath == "" {
			return fmt.Errorf("host genesis json path is required when chain name is not set")
		}

		// Slot calculator configuration required
		if env.StartTimestamp <= 0 {
			return fmt.Errorf("start timestamp must be a positive integer when chain name is not set")
		}
		if env.SlotOffset < 0 {
			return fmt.Errorf("slot offset must be a non-negative integer when chain name is not set")
		}
		if env.SlotDuration <= 0 {
			return fmt.Errorf("slot duration must be a positive integer when chain name is not set")
		}
	}

	return nil
}
