package builder

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// Service ports
	MetricsPort = 9000

	// Deployment settings
	DefaultReplicas = 1
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
	IAMRole              *iam.Role
	IAMPolicy            *iam.Policy
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
	AuthTokenRefreshInterval pulumi.StringInput `pulumi:"authTokenRefreshInterval"`
	AWSAccountId             pulumi.StringInput `pulumi:"awsAccountId"`
	AWSAccessKeyId           pulumi.StringInput `pulumi:"awsAccessKeyId"`
	AWSRegion                pulumi.StringInput `pulumi:"awsRegion"`
	AWSSecretAccessKey       pulumi.StringInput `pulumi:"awsSecretAccessKey"`
	BlockConfirmationBuffer  pulumi.StringInput `pulumi:"blockConfirmationBuffer"`
	BlockQueryCutoff         pulumi.StringInput `pulumi:"blockQueryCutoff"`
	BlockQueryStart          pulumi.StringInput `pulumi:"blockQueryStart"`
	BuilderHelperAddress     pulumi.StringInput `pulumi:"builderHelperAddress"`
	BuilderKey               pulumi.StringInput `pulumi:"builderKey"`
	BuilderPort              pulumi.IntInput    `pulumi:"builderPort"`
	BuilderRewardsAddress    pulumi.StringInput `pulumi:"builderRewardsAddress"`
	ChainOffset              pulumi.StringInput `pulumi:"chainOffset"`
	ConcurrentLimit          pulumi.StringInput `pulumi:"concurrentLimit"`
	HostChainId              pulumi.StringInput `pulumi:"hostChainId"`
	HostRpcUrl               pulumi.StringInput `pulumi:"hostRpcUrl"`
	OauthAudience            pulumi.StringInput `pulumi:"oauthAudience"`
	OauthAuthenticateUrl     pulumi.StringInput `pulumi:"oauthAuthenticateUrl"`
	OAuthClientId            pulumi.StringInput `pulumi:"oauthClientId"`
	OauthClientSecret        pulumi.StringInput `pulumi:"oauthClientSecret"`
	OauthIssuer              pulumi.StringInput `pulumi:"oauthIssuer"`
	OauthTokenUrl            pulumi.StringInput `pulumi:"oauthTokenUrl"`
	OtelExporterOtlpEndpoint pulumi.StringInput `pulumi:"otelExporterOtlpEndpoint"`
	QuinceyUrl               pulumi.StringInput `pulumi:"quinceyUrl"`
	RollupBlockGasLimit      pulumi.StringInput `pulumi:"rollupBlockGasLimit"`
	RollupChainId            pulumi.StringInput `pulumi:"rollupChainId"`
	RollupRpcUrl             pulumi.StringInput `pulumi:"rollupRpcUrl"`
	RustLog                  pulumi.StringInput `pulumi:"rustLog"`
	SlotOffset               pulumi.StringInput `pulumi:"slotOffset"`
	StartTimestamp           pulumi.StringInput `pulumi:"startTimestamp"`
	SubmitViaCallData        pulumi.StringInput `pulumi:"submitViaCallData"`
	TargetSlotTime           pulumi.StringInput `pulumi:"targetSlotTime"`
	TxBroadcastUrls          pulumi.StringInput `pulumi:"txBroadcastUrls"`
	TxPoolCacheDuration      pulumi.StringInput `pulumi:"txPoolCacheDuration"`
	TxPoolUrl                pulumi.StringInput `pulumi:"txPoolUrl"`
	ZenithAddress            pulumi.StringInput `pulumi:"zenithAddress"`
}

type IAMStatement struct {
	Sid       string `json:"sid,omitempty"`
	Effect    string `json:"effect"`
	Principal struct {
		Service []string `json:"Service"`
	} `json:"Principal"`
	Action []string `json:"Action"`
}

type IAMPolicy struct {
	Version   string         `json:"Version"`
	Statement []IAMStatement `json:"Statement"`
}

type KMSStatement struct {
	Effect   string             `json:"Effect"`
	Action   []string           `json:"Action"`
	Resource pulumi.StringInput `json:"Resource"`
}

type KMSPolicy struct {
	Version   string         `json:"Version"`
	Statement []KMSStatement `json:"Statement"`
}

type Builder interface {
	GetServiceURL() pulumi.StringOutput
	GetMetricsURL() pulumi.StringOutput
}

// Ensure BuilderComponent implements Builder
var _ Builder = &BuilderComponent{}
