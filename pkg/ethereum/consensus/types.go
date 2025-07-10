package consensus

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Public-facing structs with base Go types

// ConsensusClientArgs represents the arguments for creating a consensus client
type ConsensusClientArgs struct {
	Name                    string
	Namespace               string
	StorageSize             string
	StorageClass            string
	Image                   string
	ImagePullPolicy         string
	JWTSecret               string
	NodeSelector            pulumi.StringMap
	Tolerations             corev1.TolerationArray
	P2PPort                 int
	BeaconAPIPort           int
	MetricsPort             int
	ExecutionClientEndpoint string
	Bootnodes               []string
	AdditionalArgs          []string
}

// Internal structs with Pulumi types

type consensusClientArgsInternal struct {
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

// Conversion functions

// toInternal converts public args to internal args for use with Pulumi
func (args ConsensusClientArgs) toInternal() consensusClientArgsInternal {
	return consensusClientArgsInternal{
		Name:                    pulumi.String(args.Name),
		Namespace:               pulumi.String(args.Namespace),
		StorageSize:             pulumi.String(args.StorageSize),
		StorageClass:            pulumi.String(args.StorageClass),
		Image:                   pulumi.String(args.Image),
		ImagePullPolicy:         pulumi.String(args.ImagePullPolicy),
		JWTSecret:               pulumi.String(args.JWTSecret),
		NodeSelector:            args.NodeSelector,
		Tolerations:             args.Tolerations,
		P2PPort:                 pulumi.Int(args.P2PPort),
		BeaconAPIPort:           pulumi.Int(args.BeaconAPIPort),
		MetricsPort:             pulumi.Int(args.MetricsPort),
		ExecutionClientEndpoint: pulumi.String(args.ExecutionClientEndpoint),
		Bootnodes:               pulumi.ToStringArray(args.Bootnodes),
		AdditionalArgs:          pulumi.ToStringArray(args.AdditionalArgs),
	}
}

// ConsensusClientComponent represents a consensus client deployment
type ConsensusClientComponent struct {
	pulumi.ResourceState

	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace
	Namespace string
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
