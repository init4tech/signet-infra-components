// Package quincey provides a Pulumi component for deploying and managing the Quincey service
// in a Kubernetes cluster. It handles the creation and configuration of all necessary
// Kubernetes resources including deployments, services, and Istio configurations.
//
// The package provides a high-level interface for deploying the Quincey service with
// proper Kubernetes resource configuration, including:
// - Deployment with proper resource limits and environment variables
// - Service for internal communication
// - Istio VirtualService for external access
// - Authentication and authorization policies
package quincey

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Common constants used throughout the package
const (
	// ServiceName is the name of the Quincey service
	ServiceName = "quincey-server"
	// AppLabel is the label used to identify Quincey resources
	AppLabel = "quincey-server"
	// DefaultMetricsPort is the default port for metrics
	DefaultMetricsPort = 9000
	// ComponentName is the name of this component
	ComponentName = "quincey"
)

// QuinceyComponent represents a complete Quincey deployment in Kubernetes.
// It manages all the necessary Kubernetes resources for running the Quincey service.
type QuinceyComponent struct {
	Service               *corev1.Service
	ServiceAccount        *corev1.ServiceAccount
	Deployment            *appsv1.Deployment
	VirtualService        *crd.CustomResource
	RequestAuthentication *crd.CustomResource
	AuthorizationPolicy   *crd.CustomResource
	ConfigMap             *corev1.ConfigMap
	pulumi.ResourceState
}

// QuinceyComponentArgs contains the configuration for creating a new QuinceyComponent.
type QuinceyComponentArgs struct {
	// Namespace is the Kubernetes namespace where resources will be created
	Namespace pulumi.StringInput
	// Image is the container image to use for the Quincey service
	Image pulumi.StringInput
	// Env contains all environment variables for the Quincey service
	Env QuinceyEnv
	// Port is the port the service will listen on
	Port pulumi.StringInput
	// VirtualServiceHosts is the list of hosts for the virtual service
	VirtualServiceHosts pulumi.StringArrayInput
}

// QuinceyEnv contains all environment variables needed by the Quincey service.
// It implements the utils.EnvProvider interface for automatic environment variable handling.
type QuinceyEnv struct {
	QuinceyPort              pulumi.StringInput `pulumi:"quinceyPort" validate:"required"`
	QuinceyKeyId             pulumi.StringInput `pulumi:"quinceyKeyId" validate:"required"`
	AwsAccessKeyId           pulumi.StringInput `pulumi:"awsAccessKeyId" validate:"required"`
	AwsSecretAccessKey       pulumi.StringInput `pulumi:"awsSecretAccessKey" validate:"required"`
	AwsDefaultRegion         pulumi.StringInput `pulumi:"awsDefaultRegion" validate:"required"`
	BlockQueryStart          pulumi.StringInput `pulumi:"blockQueryStart" validate:"required"`
	BlockQueryCutoff         pulumi.StringInput `pulumi:"blockQueryCutoff" validate:"required"`
	ChainOffset              pulumi.StringInput `pulumi:"chainOffset" validate:"required"`
	HostRpcUrl               pulumi.StringInput `pulumi:"hostRpcUrl" validate:"required"`
	OauthIssuer              pulumi.StringInput `pulumi:"oauthIssuer" validate:"required"`
	OauthJwksUri             pulumi.StringInput `pulumi:"oauthJwksUri" validate:"required"`
	OtelExporterOtlpEndpoint pulumi.StringInput `pulumi:"otelExporterOtlpEndpoint"`
	OtelExporterOtlpProtocol pulumi.StringInput `pulumi:"otelExporterOtlpProtocol"`
	RustLog                  pulumi.StringInput `pulumi:"rustLog"`
	QuinceyBuilders          pulumi.StringInput `pulumi:"quinceyBuilders" validate:"required"`
}

// GetEnvMap implements the utils.EnvProvider interface
func (e *QuinceyEnv) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// Quincey defines the interface for interacting with a Quincey deployment.
type Quincey interface {
	// GetServiceURL returns the internal Kubernetes service URL for the Quincey service
	GetServiceURL() pulumi.StringOutput
	// GetMetricsURL returns the URL for accessing the metrics endpoint
	GetMetricsURL() pulumi.StringOutput
}

// Ensure QuinceyComponent implements Quincey
var _ Quincey = &QuinceyComponent{}
