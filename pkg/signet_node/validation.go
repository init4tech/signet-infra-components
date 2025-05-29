package signet_node

import (
	"fmt"
)

func (args *SignetNodeComponentArgs) Validate() error {
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}
	if args.Namespace == nil {
		return fmt.Errorf("namespace is required")
	}
	if args.ExecutionJwt == nil {
		return fmt.Errorf("execution jwt is required")
	}
	if args.LighthousePvcSize == nil {
		return fmt.Errorf("lighthouse pvc size is required")
	}
	if args.RollupPvcSize == nil {
		return fmt.Errorf("rollup pvc size is required")
	}
	if args.ExecutionClientImage == nil {
		return fmt.Errorf("execution client image is required")
	}
	if args.ConsensusClientImage == nil {
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
	if env.BaseFeeRecipient == nil {
		return fmt.Errorf("base fee recipient is required")
	}
	if env.BlobExplorerUrl == nil {
		return fmt.Errorf("blob explorer url is required")
	}
	if env.GenesisJsonPath == nil {
		return fmt.Errorf("genesis json path is required")
	}
	if env.HostOrdersContractAddress == nil {
		return fmt.Errorf("host orders contract address is required")
	}
	if env.HostPassageContractAddress == nil {
		return fmt.Errorf("host passage contract address is required")
	}
	if env.HostSlotDuration == nil {
		return fmt.Errorf("host slot duration is required")
	}
	if env.HostSlotOffset == nil {
		return fmt.Errorf("host slot offset is required")
	}
	if env.HostStartTimestamp == nil {
		return fmt.Errorf("host start timestamp is required")
	}
	if env.HostTransactorContractAddress == nil {
		return fmt.Errorf("host transactor contract address is required")
	}
	if env.HostZenithAddress == nil {
		return fmt.Errorf("host zenith address is required")
	}
	if env.HostZenithDeployHeight == nil {
		return fmt.Errorf("host zenith deploy height is required")
	}
	if env.IpcEndpoint == nil {
		return fmt.Errorf("ipc endpoint is required")
	}
	if env.RpcPort == nil {
		return fmt.Errorf("rpc port is required")
	}
	if env.RuOrdersContractAddress == nil {
		return fmt.Errorf("ru orders contract address is required")
	}
	if env.RuPassageContractAddress == nil {
		return fmt.Errorf("ru passage contract address is required")
	}
	if env.RustLog == nil {
		return fmt.Errorf("rust log is required")
	}
	if env.SignetChainId == nil {
		return fmt.Errorf("signet chain id is required")
	}
	if env.SignetClUrl == nil {
		return fmt.Errorf("signet cl url is required")
	}
	if env.SignetDatabasePath == nil {
		return fmt.Errorf("signet database path is required")
	}
	if env.SignetPylonUrl == nil {
		return fmt.Errorf("signet pylon url is required")
	}
	if env.SignetStaticPath == nil {
		return fmt.Errorf("signet static path is required")
	}
	if env.TxForwardUrl == nil {
		return fmt.Errorf("tx forward url is required")
	}
	if env.WsRpcPort == nil {
		return fmt.Errorf("ws rpc port is required")
	}
	return nil
}
