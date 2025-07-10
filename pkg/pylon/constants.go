package pylon

const (
	ServiceTypeClusterIP = "ClusterIP"
)

// Istio API versions and kinds
const (
	IstioNetworkingAPIVersion = "networking.istio.io/v1alpha3"
	IstioSecurityAPIVersion   = "security.istio.io/v1beta1"
	VirtualServiceKind        = "VirtualService"
	RequestAuthenticationKind = "RequestAuthentication"
	AuthorizationPolicyKind   = "AuthorizationPolicy"
)

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
	ConsensusP2PPort       = 9000
	ConsensusBeaconAPIPort = 4000
	ConsensusMetricsPort   = 5054
)

// Database constants
const (
	PostgreSQLPort = "5432"
)

// Image constants
const (
	ConsensusClientImage  = "sigp/lighthouse:latest"
	ImagePullPolicyAlways = "Always"
)
