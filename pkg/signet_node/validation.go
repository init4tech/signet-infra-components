package signet_node

import (
	"fmt"
)

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
	if env.BaseFeeRecipient == "" {
		return fmt.Errorf("base fee recipient is required")
	}
	if env.BlobExplorerUrl == "" {
		return fmt.Errorf("blob explorer url is required")
	}
	if env.GenesisJsonPath == "" {
		return fmt.Errorf("genesis json path is required")
	}
	if env.HostOrdersContractAddress == "" {
		return fmt.Errorf("host orders contract address is required")
	}
	if env.HostPassageContractAddress == "" {
		return fmt.Errorf("host passage contract address is required")
	}
	if env.HostSlotDuration <= 0 {
		return fmt.Errorf("host slot duration must be a positive integer")
	}
	if env.HostSlotOffset < 0 {
		return fmt.Errorf("host slot offset must be a positive integer")
	}
	if env.HostStartTimestamp <= 0 {
		return fmt.Errorf("host start timestamp must be a positive integer")
	}
	if env.HostTransactorContractAddress == "" {
		return fmt.Errorf("host transactor contract address is required")
	}
	if env.HostZenithAddress == "" {
		return fmt.Errorf("host zenith address is required")
	}
	if env.HostZenithDeployHeight == "" {
		return fmt.Errorf("host zenith deploy height is required")
	}
	if env.IpcEndpoint == "" {
		return fmt.Errorf("ipc endpoint is required")
	}
	if env.RpcPort <= 0 {
		return fmt.Errorf("rpc port must be a positive integer")
	}
	if env.RuOrdersContractAddress == "" {
		return fmt.Errorf("ru orders contract address is required")
	}
	if env.RuPassageContractAddress == "" {
		return fmt.Errorf("ru passage contract address is required")
	}
	if env.SignetChainId <= 0 {
		return fmt.Errorf("signet chain id must be a positive integer")
	}
	if env.SignetClUrl == "" {
		return fmt.Errorf("signet cl url is required")
	}
	if env.SignetDatabasePath == "" {
		return fmt.Errorf("signet database path is required")
	}
	if env.SignetPylonUrl == "" {
		return fmt.Errorf("signet pylon url is required")
	}
	if env.SignetStaticPath == "" {
		return fmt.Errorf("signet static path is required")
	}
	if env.TxForwardUrl == "" {
		return fmt.Errorf("tx forward url is required")
	}
	if env.WsRpcPort <= 0 {
		return fmt.Errorf("ws rpc port must be a positive integer")
	}
	return nil
}
