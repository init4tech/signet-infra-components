package erpcproxy

const (
	// Component identification
	ComponentKind = "signet:erpc-proxy:ErpcProxy"

	// Resource name suffixes
	ServiceAccountSuffix = "-sa"
	ConfigMapSuffix      = "-config"
	SecretSuffix         = "-secrets"
	DeploymentSuffix     = "-deployment"
	ServiceSuffix        = "-service"

	// Default values
	DefaultReplicas       = 1
	DefaultHttpPort       = 4000
	DefaultMetricsPort    = 4001
	DefaultImage          = "ghcr.io/erpc/erpc:latest"
	DefaultLogLevel       = "info"
	DefaultMaxTimeoutMs   = 30000
	DefaultMemoryRequest  = "256Mi"
	DefaultMemoryLimit    = "2Gi"
	DefaultCpuRequest     = "100m"
	DefaultCpuLimit       = "1000m"
	DefaultGoGC           = "40"
	DefaultGoMemLimit     = "1900MiB"

	// Environment variable names
	EnvGoGC       = "GOGC"
	EnvGoMemLimit = "GOMEMLIMIT"

	// Config file name
	ConfigFileName = "erpc.yaml"
	ConfigMountPath = "/etc/erpc"
)