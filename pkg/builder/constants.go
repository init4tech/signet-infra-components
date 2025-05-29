package builder

// Resource defaults
const (
	// Port defaults
	DefaultBuilderPort = 8080
	MetricsPort        = 9000

	// Deployment defaults
	DefaultReplicas = 1

	// Resource allocation defaults
	DefaultCPULimit      = "2"
	DefaultMemoryLimit   = "2Gi"
	DefaultCPURequest    = "1"
	DefaultMemoryRequest = "1Gi"

	// Component kind
	ComponentKind = "signet:index:Builder"
)

// Resource name suffixes
const (
	ServiceSuffix        = "-service"
	DeploymentSuffix     = "-deployment"
	ServiceAccountSuffix = "-sa"
	ConfigMapSuffix      = "-env"
	PodMonitorSuffix     = "-pod-monitor"
	ServiceMonitorSuffix = "-svcmon"
)

// Health check paths
const (
	HealthCheckPath = "/healthcheck"
	MetricsPath     = "/metrics"
)

// Probe settings
const (
	ProbeInitialDelaySeconds = 5
	ProbePeriodSeconds       = 10
	LivenessProbePeriod      = 1
	ProbeTimeoutSeconds      = 1
	ProbeFailureThreshold    = 3
)

// Prometheus annotations
const (
	PrometheusScrapeAnnotation = "prometheus.io/scrape"
	PrometheusPortAnnotation   = "prometheus.io/port"
	PrometheusPathAnnotation   = "prometheus.io/path"
)
