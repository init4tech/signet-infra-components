package builder

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AppLabels represents the Kubernetes labels to be applied to the builder resources.
type AppLabels struct {
	Labels pulumi.StringMap
}

// BuilderComponent represents a Pulumi component that deploys a builder service.
type BuilderComponent struct {
	pulumi.ResourceState
	BuilderComponentArgs BuilderComponentArgs
	Deployment           *appsv1.Deployment
	Service              *corev1.Service
	ServiceAccount       *corev1.ServiceAccount
	ConfigMap            *corev1.ConfigMap
}

// BuilderComponentArgs contains the configuration for deploying a builder service.
type BuilderComponentArgs struct {
	Namespace  string     // k8s namespace to deploy the builder to
	AppLabels  AppLabels  // Labels to apply to the builder pod
	Name       string     // Builder name identifier
	Image      string     // Builder docker image
	BuilderEnv BuilderEnv // Builder environment variables
}

// BuilderEnv contains all environment variables required by the builder service.
type BuilderEnv struct {
	AuthTokenRefreshInterval pulumi.StringInput `pulumi:"authTokenRefreshInterval" validate:"required"`
	AwsAccountId             pulumi.StringInput `pulumi:"awsAccountId" validate:"required"`
	AwsAccessKeyId           pulumi.StringInput `pulumi:"awsAccessKeyId" validate:"required"`
	AwsRegion                pulumi.StringInput `pulumi:"awsRegion" validate:"required"`
	AwsSecretAccessKey       pulumi.StringInput `pulumi:"awsSecretAccessKey" validate:"required"`
	BlockConfirmationBuffer  pulumi.StringInput `pulumi:"blockConfirmationBuffer" validate:"required"`
	BlockQueryCutoff         pulumi.StringInput `pulumi:"blockQueryCutoff" validate:"required"`
	BlockQueryStart          pulumi.StringInput `pulumi:"blockQueryStart" validate:"required"`
	BuilderHelperAddress     pulumi.StringInput `pulumi:"builderHelperAddress" validate:"required"`
	BuilderKey               pulumi.StringInput `pulumi:"builderKey" validate:"required"`
	BuilderPort              pulumi.StringInput `pulumi:"builderPort" validate:"required"`
	BuilderRewardsAddress    pulumi.StringInput `pulumi:"builderRewardsAddress" validate:"required"`
	ChainOffset              pulumi.StringInput `pulumi:"chainOffset" validate:"required"`
	ConcurrentLimit          pulumi.StringInput `pulumi:"concurrentLimit" validate:"required"`
	HostChainId              pulumi.StringInput `pulumi:"hostChainId" validate:"required"`
	HostRpcUrl               pulumi.StringInput `pulumi:"hostRpcUrl" validate:"required"`
	OauthAudience            pulumi.StringInput `pulumi:"oauthAudience" validate:"required"`
	OauthAuthenticateUrl     pulumi.StringInput `pulumi:"oauthAuthenticateUrl" validate:"required"`
	OAuthClientId            pulumi.StringInput `pulumi:"oauthClientId" validate:"required"`
	OauthClientSecret        pulumi.StringInput `pulumi:"oauthClientSecret" validate:"required"`
	OauthIssuer              pulumi.StringInput `pulumi:"oauthIssuer" validate:"required"`
	OauthTokenUrl            pulumi.StringInput `pulumi:"oauthTokenUrl" validate:"required"`
	OtelExporterOtlpEndpoint pulumi.StringInput `pulumi:"otelExporterOtlpEndpoint"`
	QuinceyUrl               pulumi.StringInput `pulumi:"quinceyUrl" validate:"required"`
	RollupBlockGasLimit      pulumi.StringInput `pulumi:"rollupBlockGasLimit" validate:"required"`
	RollupChainId            pulumi.StringInput `pulumi:"rollupChainId" validate:"required"`
	RollupRpcUrl             pulumi.StringInput `pulumi:"rollupRpcUrl" validate:"required"`
	RustLog                  pulumi.StringInput `pulumi:"rustLog"`
	SlotOffset               pulumi.StringInput `pulumi:"slotOffset" validate:"required"`
	StartTimestamp           pulumi.StringInput `pulumi:"startTimestamp" validate:"required"`
	SubmitViaCallData        pulumi.StringInput `pulumi:"submitViaCallData" validate:"required"`
	TargetSlotTime           pulumi.StringInput `pulumi:"targetSlotTime" validate:"required"`
	TxBroadcastUrls          pulumi.StringInput `pulumi:"txBroadcastUrls" validate:"required"`
	TxPoolCacheDuration      pulumi.StringInput `pulumi:"txPoolCacheDuration" validate:"required"`
	TxPoolUrl                pulumi.StringInput `pulumi:"txPoolUrl" validate:"required"`
	ZenithAddress            pulumi.StringInput `pulumi:"zenithAddress" validate:"required"`
}

// GetEnvMap implements the utils.EnvProvider interface
// It creates a map of environment variables from the BuilderEnv struct
func (e BuilderEnv) GetEnvMap() pulumi.StringMap {
	// All fields are now StringInput, so we can use the standard reflection method
	return utils.CreateEnvMap(e)
}

type Builder interface {
	GetServiceURL() pulumi.StringOutput
	GetMetricsURL() pulumi.StringOutput
}

// Ensure BuilderComponent implements Builder
var _ Builder = &BuilderComponent{}
