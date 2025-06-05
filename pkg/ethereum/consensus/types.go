package consensus

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ConsensusClientArgs represents the arguments for creating a consensus client
type ConsensusClientArgs struct {
	Name                    pulumi.StringInput
	Namespace               pulumi.StringInput
	StorageSize             pulumi.StringInput
	StorageClass            pulumi.StringInput
	Image                   pulumi.StringInput
	ImagePullPolicy         pulumi.StringInput
	JWTSecret               pulumi.StringInput
	NodeSelector            pulumi.StringMap
	Tolerations             corev1.TolerationArray
	P2PPort                 pulumi.IntInput
	BeaconAPIPort           pulumi.IntInput
	MetricsPort             pulumi.IntInput
	ExecutionClientEndpoint pulumi.StringInput
	Bootnodes               pulumi.StringArray
	AdditionalArgs          pulumi.StringArray
}

// ConsensusClientComponent represents a consensus client deployment
type ConsensusClientComponent struct {
	pulumi.ResourceState

	// Name is the base name for all resources
	Name pulumi.StringOutput
	// Namespace is the Kubernetes namespace
	Namespace pulumi.StringOutput
	// PVC is the persistent volume claim
	PVC *corev1.PersistentVolumeClaim
	// JWTSecret is the JWT secret
	JWTSecret *corev1.Secret
	// P2PService is the P2P service
	P2PService *corev1.Service
	// BeaconAPIService is the beacon API service
	BeaconAPIService *corev1.Service
	// StatefulSet is the stateful set
	StatefulSet *appsv1.StatefulSet
}
