package execution

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ExecutionClientArgs contains the configuration for an execution client
type ExecutionClientArgs struct {
	// Name is the base name for all resources
	Name string `pulumi:"name"`
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace string `pulumi:"namespace"`
	// StorageSize is the size of the persistent volume claim
	StorageSize string `pulumi:"storageSize"`
	// StorageClass is the Kubernetes storage class to use
	StorageClass string `pulumi:"storageClass"`
	// Image is the container image to use
	Image string `pulumi:"image"`
	// ImagePullPolicy is the Kubernetes image pull policy
	ImagePullPolicy string `pulumi:"imagePullPolicy"`
	// Resources contains the resource requests and limits
	Resources *corev1.ResourceRequirements `pulumi:"resources,optional"`
	// NodeSelector is the Kubernetes node selector
	NodeSelector map[string]string `pulumi:"nodeSelector,optional"`
	// Tolerations are the Kubernetes tolerations
	Tolerations []corev1.Toleration `pulumi:"tolerations,optional"`
	// JWTSecret is the JWT secret for authentication
	JWTSecret string `pulumi:"jwtSecret"`
	// P2PPort is the port for P2P communication
	P2PPort int `pulumi:"p2pPort"`
	// RPCPort is the port for RPC communication
	RPCPort int `pulumi:"rpcPort"`
	// WSPort is the port for WebSocket communication
	WSPort int `pulumi:"wsPort"`
	// MetricsPort is the port for metrics
	MetricsPort int `pulumi:"metricsPort"`
	// AuthRPCPort is the port for authenticated RPC
	AuthRPCPort int `pulumi:"authRpcPort"`
	// DiscoveryPort is the port for node discovery
	DiscoveryPort int `pulumi:"discoveryPort"`
	// Bootnodes is a list of bootnode URLs
	Bootnodes []string `pulumi:"bootnodes,optional"`
	// AdditionalArgs are additional command line arguments
	AdditionalArgs []string `pulumi:"additionalArgs,optional"`
}

// ExecutionClientComponent represents an execution client deployment
type ExecutionClientComponent struct {
	pulumi.ResourceState

	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace
	Namespace string
	// ConfigMap is the shared config map
	ConfigMap *corev1.ConfigMap
	// PVC is the persistent volume claim
	PVC *corev1.PersistentVolumeClaim
	// JWTSecret is the JWT secret
	JWTSecret *corev1.Secret
	// P2PService is the P2P service
	P2PService *corev1.Service
	// RPCService is the RPC service
	RPCService *corev1.Service
	// StatefulSet is the stateful set
	StatefulSet *appsv1.StatefulSet
}
