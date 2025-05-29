package signet_node

import (
	"fmt"
)

func (args *SignetNodeComponentArgs) Validate() error {
	if args.Namespace == nil {
		return fmt.Errorf("namespace is required")
	}
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}
	if err := args.Env.Validate(); err != nil {
		return fmt.Errorf("invalid signet node env: %w", err)
	}
	return nil
}

func (env *SignetNodeEnv) Validate() error {
	if env.HostZenithAddress == nil {
		return fmt.Errorf("host zenith address is required")
	}
	if env.HostOrdersContractAddress == nil {
		return fmt.Errorf("host orders contract address is required")
	}
	if env.SignetChainId == nil {
		return fmt.Errorf("signet chain id is required")
	}
	if env.BlobExplorerUrl == nil {
		return fmt.Errorf("blob explorer url is required")
	}
	if env.SignetStaticPath == nil {
		return fmt.Errorf("signet static path is required")
	}
	if env.SignetDatabasePath == nil {
		return fmt.Errorf("signet database path is required")
	}
	if env.RpcPort == nil {
		return fmt.Errorf("rpc port is required")
	}
	if env.WsRpcPort == nil {
		return fmt.Errorf("ws rpc port is required")
	}
	if env.TxForwardUrl == nil {
		return fmt.Errorf("tx forward url is required")
	}
	if env.GenesisJsonPath == nil {
		return fmt.Errorf("genesis json path is required")
	}
	if env.HostStartTimestamp == nil {
		return fmt.Errorf("host start timestamp is required")
	}
	if env.HostSlotOffset == nil {
		return fmt.Errorf("host slot offset is required")
	}
	if env.HostSlotDuration == nil {
		return fmt.Errorf("host slot duration is required")
	}
	return nil
}
