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
