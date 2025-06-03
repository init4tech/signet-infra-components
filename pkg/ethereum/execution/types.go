package execution

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ExecutionClientArgs contains the configuration for an execution client
type ExecutionClientArgs struct {
	// Name is the base name for all resources
	Name pulumi.StringInput `pulumi:"name"`
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace pulumi.StringInput `pulumi:"namespace"`
	// StorageSize is the size of the persistent volume claim
	StorageSize pulumi.StringInput `pulumi:"storageSize"`
	// StorageClass is the Kubernetes storage class to use
	StorageClass pulumi.StringInput `pulumi:"storageClass"`
	// Image is the container image to use
	Image pulumi.StringInput `pulumi:"image"`
	// ImagePullPolicy is the Kubernetes image pull policy
	ImagePullPolicy pulumi.StringInput `pulumi:"imagePullPolicy"`
	// Resources contains the resource requests and limits
	Resources *corev1.ResourceRequirements `pulumi:"resources,optional"`
	// NodeSelector is the Kubernetes node selector
	NodeSelector pulumi.StringMap `pulumi:"nodeSelector,optional"`
	// Tolerations are the Kubernetes tolerations
	Tolerations corev1.TolerationArray `pulumi:"tolerations,optional"`
	// JWTSecret is the JWT secret for authentication
	JWTSecret pulumi.StringInput `pulumi:"jwtSecret"`
	// P2PPort is the port for P2P communication
	P2PPort pulumi.IntInput `pulumi:"p2pPort"`
	// RPCPort is the port for RPC communication
	RPCPort pulumi.IntInput `pulumi:"rpcPort"`
	// WSPort is the port for WebSocket communication
	WSPort pulumi.IntInput `pulumi:"wsPort"`
	// MetricsPort is the port for metrics
	MetricsPort pulumi.IntInput `pulumi:"metricsPort"`
	// AuthRPCPort is the port for authenticated RPC
	AuthRPCPort pulumi.IntInput `pulumi:"authRpcPort"`
	// DiscoveryPort is the port for node discovery
	DiscoveryPort pulumi.IntInput `pulumi:"discoveryPort"`
	// Bootnodes is a list of bootnode URLs
	Bootnodes pulumi.StringArray `pulumi:"bootnodes,optional"`
	// AdditionalArgs are additional command line arguments
	AdditionalArgs pulumi.StringArray `pulumi:"additionalArgs,optional"`
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
