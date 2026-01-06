package signet_node

// Resource defaults
const (
	// Storage defaults
	DefaultStorageSize  = "150Gi"
	DefaultStorageClass = "aws-gp3"

	// StatefulSet defaults
	DefaultReplicas = 1

	// Resource allocation defaults
	DefaultCPULimit      = "2"
	DefaultMemoryLimit   = "16Gi"
	DefaultCPURequest    = "2"
	DefaultMemoryRequest = "4Gi"

	// Port defaults
	MetricsPort          = 9001
	RpcPort              = 8545
	WsPort               = 8546
	AuthRpcPort          = 8551
	DiscoveryPort        = 30303
	ConsensusHttpPort    = 4000
	ConsensusMetricsPort = 5054
	HostIpcPort          = 8547
	RollupHttpPort       = 8645
	RollupWsPort         = 8646

	// SignetNode environment defaults
	DefaultSignetRpcPort   = 5959
	DefaultSignetWsRpcPort = 5960

	// Component name
	ComponentKind = "the-builder:index:SignetNode"

	// Default mount paths for volumes
	DefaultSignetNodeDataMountPath = "/root/.local/share/reth"
	DefaultRollupDataMountPath     = "/root/.local/share/exex"
	DefaultExecutionJwtMountPath   = "/etc/reth/execution-jwt"
)

// Resource name suffixes
const (
	ServiceSuffix        = "-service"
	StatefulSetSuffix    = "-set"
	ConfigMapSuffix      = "-configmap"
	PvcSuffix            = "-data"
	SecretSuffix         = "-secret"
	VirtualServiceSuffix = "-vservice"
)

// Resource names
const (
	SignetNodeName      = "signet-node"
	LighthouseName      = "lighthouse"
	ExecutionJwtName    = "execution-jwt"
	RollupDataName      = "rollup-data"
	ExecutionConfigName = "exex-configmap"
)
