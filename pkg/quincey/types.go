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
	"strconv"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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

// Public-facing structs with base Go types
type QuinceyComponentArgs struct {
	// Namespace is the Kubernetes namespace where resources will be created
	Namespace string
	// Image is the container image to use for the Quincey service
	Image string
	// Env contains all environment variables for the Quincey service
	Env QuinceyEnv
	// Port is the port the service will listen on
	Port int
	// VirtualServiceHosts is the list of hosts for the virtual service
	VirtualServiceHosts []string
}

// Internal structs with Pulumi types for use within the component
type quinceyComponentArgsInternal struct {
	// Namespace is the Kubernetes namespace where resources will be created
	Namespace pulumi.StringInput
	// Image is the container image to use for the Quincey service
	Image pulumi.StringInput
	// Env contains all environment variables for the Quincey service
	Env quinceyEnvInternal
	// Port is the port the service will listen on
	Port pulumi.StringInput
	// VirtualServiceHosts is the list of hosts for the virtual service
	VirtualServiceHosts pulumi.StringArrayInput
}

// Public-facing environment struct with base Go types
type QuinceyEnv struct {
	QuinceyPort              int    `pulumi:"quinceyPort" validate:"required"`
	QuinceyKeyId             string `pulumi:"quinceyKeyId" validate:"required"`
	AwsAccessKeyId           string `pulumi:"awsAccessKeyId" validate:"required"`
	AwsSecretAccessKey       string `pulumi:"awsSecretAccessKey" validate:"required"`
	AwsDefaultRegion         string `pulumi:"awsDefaultRegion" validate:"required"`
	BlockQueryStart          int    `pulumi:"blockQueryStart" validate:"required"`
	BlockQueryCutoff         int    `pulumi:"blockQueryCutoff" validate:"required"`
	ChainOffset              int    `pulumi:"chainOffset" validate:"required"`
	HostRpcUrl               string `pulumi:"hostRpcUrl" validate:"required"`
	OauthIssuer              string `pulumi:"oauthIssuer" validate:"required"`
	OauthJwksUri             string `pulumi:"oauthJwksUri" validate:"required"`
	OtelExporterOtlpEndpoint string `pulumi:"otelExporterOtlpEndpoint"`
	OtelExporterOtlpProtocol string `pulumi:"otelExporterOtlpProtocol"`
	RustLog                  string `pulumi:"rustLog"`
	QuinceyBuilders          string `pulumi:"quinceyBuilders" validate:"required"`
}

// Internal environment struct with Pulumi types
type quinceyEnvInternal struct {
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

// Conversion function to convert public args to internal args
func (args QuinceyComponentArgs) toInternal() quinceyComponentArgsInternal {
	return quinceyComponentArgsInternal{
		Namespace:           pulumi.String(args.Namespace),
		Image:               pulumi.String(args.Image),
		Env:                 args.Env.toInternal(),
		Port:                pulumi.String(strconv.Itoa(args.Port)),
		VirtualServiceHosts: pulumi.ToStringArray(args.VirtualServiceHosts),
	}
}

// Conversion function to convert public env to internal env
func (e QuinceyEnv) toInternal() quinceyEnvInternal {
	return quinceyEnvInternal{
		QuinceyPort:              pulumi.String(strconv.Itoa(e.QuinceyPort)),
		QuinceyKeyId:             pulumi.String(e.QuinceyKeyId),
		AwsAccessKeyId:           pulumi.String(e.AwsAccessKeyId),
		AwsSecretAccessKey:       pulumi.String(e.AwsSecretAccessKey),
		AwsDefaultRegion:         pulumi.String(e.AwsDefaultRegion),
		BlockQueryStart:          pulumi.String(strconv.Itoa(e.BlockQueryStart)),
		BlockQueryCutoff:         pulumi.String(strconv.Itoa(e.BlockQueryCutoff)),
		ChainOffset:              pulumi.String(strconv.Itoa(e.ChainOffset)),
		HostRpcUrl:               pulumi.String(e.HostRpcUrl),
		OauthIssuer:              pulumi.String(e.OauthIssuer),
		OauthJwksUri:             pulumi.String(e.OauthJwksUri),
		OtelExporterOtlpEndpoint: pulumi.String(e.OtelExporterOtlpEndpoint),
		OtelExporterOtlpProtocol: pulumi.String(e.OtelExporterOtlpProtocol),
		RustLog:                  pulumi.String(e.RustLog),
		QuinceyBuilders:          pulumi.String(e.QuinceyBuilders),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal env
func (e quinceyEnvInternal) GetEnvMap() pulumi.StringMap {
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
