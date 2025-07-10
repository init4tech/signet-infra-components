package execution

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Public-facing structs with base Go types

// ExecutionClientArgs contains the configuration for an execution client
type ExecutionClientArgs struct {
	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace string
	// StorageSize is the size of the persistent volume claim
	StorageSize string
	// StorageClass is the Kubernetes storage class to use
	StorageClass string
	// Image is the container image to use
	Image string
	// ImagePullPolicy is the Kubernetes image pull policy
	ImagePullPolicy string
	// Resources contains the resource requests and limits
	Resources *corev1.ResourceRequirements
	// NodeSelector is the Kubernetes node selector
	NodeSelector pulumi.StringMap
	// Tolerations are the Kubernetes tolerations
	Tolerations corev1.TolerationArray
	// JWTSecret is the JWT secret for authentication
	JWTSecret string
	// P2PPort is the port for P2P communication
	P2PPort int
	// RPCPort is the port for RPC communication
	RPCPort int
	// WSPort is the port for WebSocket communication
	WSPort int
	// MetricsPort is the port for metrics
	MetricsPort int
	// AuthRPCPort is the port for authenticated RPC
	AuthRPCPort int
	// DiscoveryPort is the port for node discovery
	DiscoveryPort int
	// Bootnodes is a list of bootnode URLs
	Bootnodes []string
	// AdditionalArgs are additional command line arguments
	AdditionalArgs []string
	// Environment variables for the execution client, accepts a generic type that implements the utils.EnvProvider interface
	ExecutionClientEnv utils.EnvProvider
}

// Internal structs with Pulumi types

type executionClientArgsInternal struct {
	// Name is the base name for all resources
	Name pulumi.StringInput
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace pulumi.StringInput
	// StorageSize is the size of the persistent volume claim
	StorageSize pulumi.StringInput
	// StorageClass is the Kubernetes storage class to use
	StorageClass pulumi.StringInput
	// Image is the container image to use
	Image pulumi.StringInput
	// ImagePullPolicy is the Kubernetes image pull policy
	ImagePullPolicy pulumi.StringInput
	// Resources contains the resource requests and limits
	Resources *corev1.ResourceRequirements
	// NodeSelector is the Kubernetes node selector
	NodeSelector pulumi.StringMap
	// Tolerations are the Kubernetes tolerations
	Tolerations corev1.TolerationArray
	// JWTSecret is the JWT secret for authentication
	JWTSecret pulumi.StringInput
	// P2PPort is the port for P2P communication
	P2PPort pulumi.IntInput
	// RPCPort is the port for RPC communication
	RPCPort pulumi.IntInput
	// WSPort is the port for WebSocket communication
	WSPort pulumi.IntInput
	// MetricsPort is the port for metrics
	MetricsPort pulumi.IntInput
	// AuthRPCPort is the port for authenticated RPC
	AuthRPCPort pulumi.IntInput
	// DiscoveryPort is the port for node discovery
	DiscoveryPort pulumi.IntInput
	// Bootnodes is a list of bootnode URLs
	Bootnodes pulumi.StringArray
	// AdditionalArgs are additional command line arguments
	AdditionalArgs pulumi.StringArray
	// Environment variables for the execution client, accepts a generic type that implements the utils.EnvProvider interface
	ExecutionClientEnv utils.EnvProvider
}

// Conversion functions

// toInternal converts public args to internal args for use with Pulumi
func (args ExecutionClientArgs) toInternal() executionClientArgsInternal {
	return executionClientArgsInternal{
		Name:               pulumi.String(args.Name),
		Namespace:          pulumi.String(args.Namespace),
		StorageSize:        pulumi.String(args.StorageSize),
		StorageClass:       pulumi.String(args.StorageClass),
		Image:              pulumi.String(args.Image),
		ImagePullPolicy:    pulumi.String(args.ImagePullPolicy),
		Resources:          args.Resources,
		NodeSelector:       args.NodeSelector,
		Tolerations:        args.Tolerations,
		JWTSecret:          pulumi.String(args.JWTSecret),
		P2PPort:            pulumi.Int(args.P2PPort),
		RPCPort:            pulumi.Int(args.RPCPort),
		WSPort:             pulumi.Int(args.WSPort),
		MetricsPort:        pulumi.Int(args.MetricsPort),
		AuthRPCPort:        pulumi.Int(args.AuthRPCPort),
		DiscoveryPort:      pulumi.Int(args.DiscoveryPort),
		Bootnodes:          pulumi.ToStringArray(args.Bootnodes),
		AdditionalArgs:     pulumi.ToStringArray(args.AdditionalArgs),
		ExecutionClientEnv: args.ExecutionClientEnv,
	}
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
