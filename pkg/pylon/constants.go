package pylon

// Storage constants
const (
	ExecutionClientStorageSize = "150Gi"
	ConsensusClientStorageSize = "100Gi"
	StorageClassAWSGP3         = "aws-gp3"
)

// Port constants
const (
	ExecutionP2PPort       = 30303
	ExecutionRPCPort       = 8545
	ExecutionWSPort        = 8546
	ExecutionMetricsPort   = 9001
	ExecutionAuthRPCPort   = 8551
	ConsensusBeaconAPIPort = 4000
	ConsensusMetricsPort   = 5054
)

// Image constants
const (
	ConsensusClientImage  = "sigp/lighthouse:latest"
	ImagePullPolicyAlways = "Always"
)
