package quincey

// Resource defaults
const (
	// Port defaults
	DefaultQuinceyPort = 8080
	DefaultMetricsPort = 9000

	// Deployment defaults
	DefaultReplicas = 1

	// Component kind
	ComponentKind = "signet:index:Quincey"
)

// Resource names and identifiers
const (
	// ServiceName is the name of the Quincey service
	ServiceName = "quincey-server"
	// AppLabel is the label used to identify Quincey resources
	AppLabel = "quincey-server"
	// ComponentName is the name of this component
	ComponentName = "quincey"
)

// Resource name suffixes
const (
	ServiceSuffix        = "-service"
	DeploymentSuffix     = "-deployment"
	ServiceAccountSuffix = "-sa"
	ConfigMapSuffix      = "-configmap"
	VirtualServiceSuffix = "-vservice"
	RequestAuthSuffix    = "-request-auth"
	AuthPolicySuffix     = "-auth-policy"
)

// Service types
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

// JWT and OAuth constants
const (
	JWTTokenHeader     = "authorization"
	JWTTokenPrefix     = "Bearer "
	OAuthIssuerClaim   = "iss"
	DefaultAppSelector = "signet"
)
