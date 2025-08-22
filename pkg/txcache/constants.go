package txcache

// Component and resource names
const (
	ComponentResourceType = "signet:txcache:TransactionCache"
	AppLabel              = "tx-cache"
	ContainerName         = "tx-cache-container"
	ServiceAccountName    = "tx-cache-service-account"
	DeploymentName        = "tx-cache"
	ServiceName           = "tx-cache-service"
	VirtualServiceName    = "tx-cache-virtual-service"
	JwtPolicyName         = "tx-cache-jwt-policy"
	AuthPolicyName        = "tx-cache-authorization-policy"
)

// Kubernetes API versions and kinds
const (
	VirtualServiceAPIVersion = "networking.istio.io/v1alpha3"
	VirtualServiceKind       = "VirtualService"
	RequestAuthAPIVersion    = "security.istio.io/v1beta1"
	RequestAuthKind          = "RequestAuthentication"
	AuthPolicyAPIVersion     = "security.istio.io/v1beta1"
	AuthPolicyKind           = "AuthorizationPolicy"
)

// Service configuration
const (
	ServiceTypeClusterIP = "ClusterIP"
	ImagePullPolicy      = "Always"
	ReplicaCount         = 1
)

// Istio configuration
const (
	VirtualServiceHost = "transactions.pecorino.signet.sh"
	GatewayName        = "default/init4-api-gateway"
	UriPrefix          = "/"
)

// JWT configuration
const (
	JwtClaimSub              = "sub"
	JwtHeaderSub             = "x-jwt-claim-sub"
	RequestPrincipalWildcard = "*"
)

// HTTP methods and paths
const (
	HttpMethodGet       = "GET"
	BundlesPath         = "/bundles"
	BundlesWildcardPath = "/bundles/*"
)

// Authorization policy actions
const (
	AuthActionAllow = "ALLOW"
)
