package builder

import (
	"strconv"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AppLabels represents the Kubernetes labels to be applied to the builder resources.
type AppLabels struct {
	Labels pulumi.StringMap
}

// Public-facing struct for builder component args
// All fields are base Go types
// Internal struct will use Pulumi types

type BuilderComponentArgs struct {
	Namespace  string     // k8s namespace to deploy the builder to
	AppLabels  AppLabels  // Labels to apply to the builder pod
	Name       string     // Builder name identifier
	Image      string     // Builder docker image
	BuilderEnv BuilderEnv // Builder environment variables
}

type builderComponentArgsInternal struct {
	Namespace  pulumi.StringInput
	AppLabels  AppLabels
	Name       string
	Image      pulumi.StringInput
	BuilderEnv builderEnvInternal
}

// Public-facing struct for builder environment variables
// All fields are base Go types
// Use int for numeric fields, string for others

type BuilderEnv struct {
	AuthTokenRefreshInterval string
	AwsAccountId             string
	AwsAccessKeyId           string
	AwsRegion                string
	AwsSecretAccessKey       string
	BlockConfirmationBuffer  int
	BlockQueryCutoff         int
	BlockQueryStart          int
	BuilderHelperAddress     string
	BuilderKey               string
	BuilderPort              int
	BuilderRewardsAddress    string
	ChainOffset              int
	ConcurrentLimit          int
	HostChainId              int
	HostRpcUrl               string
	OauthAudience            string
	OauthAuthenticateUrl     string
	OAuthClientId            string
	OauthClientSecret        string
	OauthIssuer              string
	OauthTokenUrl            string
	OtelExporterOtlpEndpoint string
	QuinceyUrl               string
	RollupBlockGasLimit      int
	RollupChainId            int
	RollupRpcUrl             string
	RustLog                  string
	SlotOffset               int
	StartTimestamp           int
	SubmitViaCallData        string
	TargetSlotTime           int
	TxBroadcastUrls          string
	TxPoolCacheDuration      int
	TxPoolUrl                string
	ZenithAddress            string
}

type builderEnvInternal struct {
	AuthTokenRefreshInterval pulumi.StringInput
	AwsAccountId             pulumi.StringInput
	AwsAccessKeyId           pulumi.StringInput
	AwsRegion                pulumi.StringInput
	AwsSecretAccessKey       pulumi.StringInput
	BlockConfirmationBuffer  pulumi.StringInput
	BlockQueryCutoff         pulumi.StringInput
	BlockQueryStart          pulumi.StringInput
	BuilderHelperAddress     pulumi.StringInput
	BuilderKey               pulumi.StringInput
	BuilderPort              pulumi.StringInput
	BuilderRewardsAddress    pulumi.StringInput
	ChainOffset              pulumi.StringInput
	ConcurrentLimit          pulumi.StringInput
	HostChainId              pulumi.StringInput
	HostRpcUrl               pulumi.StringInput
	OauthAudience            pulumi.StringInput
	OauthAuthenticateUrl     pulumi.StringInput
	OAuthClientId            pulumi.StringInput
	OauthClientSecret        pulumi.StringInput
	OauthIssuer              pulumi.StringInput
	OauthTokenUrl            pulumi.StringInput
	OtelExporterOtlpEndpoint pulumi.StringInput
	QuinceyUrl               pulumi.StringInput
	RollupBlockGasLimit      pulumi.StringInput
	RollupChainId            pulumi.StringInput
	RollupRpcUrl             pulumi.StringInput
	RustLog                  pulumi.StringInput
	SlotOffset               pulumi.StringInput
	StartTimestamp           pulumi.StringInput
	SubmitViaCallData        pulumi.StringInput
	TargetSlotTime           pulumi.StringInput
	TxBroadcastUrls          pulumi.StringInput
	TxPoolCacheDuration      pulumi.StringInput
	TxPoolUrl                pulumi.StringInput
	ZenithAddress            pulumi.StringInput
}

// Conversion function for BuilderComponentArgs
func (args BuilderComponentArgs) toInternal() builderComponentArgsInternal {
	return builderComponentArgsInternal{
		Namespace:  pulumi.String(args.Namespace),
		AppLabels:  args.AppLabels,
		Name:       args.Name,
		Image:      pulumi.String(args.Image),
		BuilderEnv: args.BuilderEnv.toInternal(),
	}
}

// Conversion function for BuilderEnv
func (e BuilderEnv) toInternal() builderEnvInternal {
	return builderEnvInternal{
		AuthTokenRefreshInterval: pulumi.String(e.AuthTokenRefreshInterval),
		AwsAccountId:             pulumi.String(e.AwsAccountId),
		AwsAccessKeyId:           pulumi.String(e.AwsAccessKeyId),
		AwsRegion:                pulumi.String(e.AwsRegion),
		AwsSecretAccessKey:       pulumi.String(e.AwsSecretAccessKey),
		BlockConfirmationBuffer:  pulumi.String(strconv.Itoa(e.BlockConfirmationBuffer)),
		BlockQueryCutoff:         pulumi.String(strconv.Itoa(e.BlockQueryCutoff)),
		BlockQueryStart:          pulumi.String(strconv.Itoa(e.BlockQueryStart)),
		BuilderHelperAddress:     pulumi.String(e.BuilderHelperAddress),
		BuilderKey:               pulumi.String(e.BuilderKey),
		BuilderPort:              pulumi.String(strconv.Itoa(e.BuilderPort)),
		BuilderRewardsAddress:    pulumi.String(e.BuilderRewardsAddress),
		ChainOffset:              pulumi.String(strconv.Itoa(e.ChainOffset)),
		ConcurrentLimit:          pulumi.String(strconv.Itoa(e.ConcurrentLimit)),
		HostChainId:              pulumi.String(strconv.Itoa(e.HostChainId)),
		HostRpcUrl:               pulumi.String(e.HostRpcUrl),
		OauthAudience:            pulumi.String(e.OauthAudience),
		OauthAuthenticateUrl:     pulumi.String(e.OauthAuthenticateUrl),
		OAuthClientId:            pulumi.String(e.OAuthClientId),
		OauthClientSecret:        pulumi.String(e.OauthClientSecret),
		OauthIssuer:              pulumi.String(e.OauthIssuer),
		OauthTokenUrl:            pulumi.String(e.OauthTokenUrl),
		OtelExporterOtlpEndpoint: pulumi.String(e.OtelExporterOtlpEndpoint),
		QuinceyUrl:               pulumi.String(e.QuinceyUrl),
		RollupBlockGasLimit:      pulumi.String(strconv.Itoa(e.RollupBlockGasLimit)),
		RollupChainId:            pulumi.String(strconv.Itoa(e.RollupChainId)),
		RollupRpcUrl:             pulumi.String(e.RollupRpcUrl),
		RustLog:                  pulumi.String(e.RustLog),
		SlotOffset:               pulumi.String(strconv.Itoa(e.SlotOffset)),
		StartTimestamp:           pulumi.String(strconv.Itoa(e.StartTimestamp)),
		SubmitViaCallData:        pulumi.String(e.SubmitViaCallData),
		TargetSlotTime:           pulumi.String(strconv.Itoa(e.TargetSlotTime)),
		TxBroadcastUrls:          pulumi.String(e.TxBroadcastUrls),
		TxPoolCacheDuration:      pulumi.String(strconv.Itoa(e.TxPoolCacheDuration)),
		TxPoolUrl:                pulumi.String(e.TxPoolUrl),
		ZenithAddress:            pulumi.String(e.ZenithAddress),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal env
func (e builderEnvInternal) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// GetEnvMap returns the environment variables as a pulumi.StringMap for the public BuilderEnv
func (e BuilderEnv) GetEnvMap() pulumi.StringMap {
	return e.toInternal().GetEnvMap()
}

type Builder interface {
	GetServiceURL() pulumi.StringOutput
	GetMetricsURL() pulumi.StringOutput
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

// Ensure BuilderComponent implements Builder
var _ Builder = &BuilderComponent{}
