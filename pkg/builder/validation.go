package builder

import (
	"fmt"
)

// Validate validates the BuilderComponentArgs
func (args *BuilderComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if args.Name == "" {
		return fmt.Errorf("name is required")
	}
	if args.Image == "" {
		return fmt.Errorf("image is required")
	}
	if args.AppLabels.Labels == nil {
		return fmt.Errorf("app labels are required")
	}
	return args.BuilderEnv.Validate()
}

// Validate validates the BuilderEnv
func (env *BuilderEnv) Validate() error {
	if env.AuthTokenRefreshInterval == "" {
		return fmt.Errorf("auth token refresh interval is required")
	}
	if env.AwsAccountId == "" {
		return fmt.Errorf("aws account id is required")
	}
	if env.AwsAccessKeyId == "" {
		return fmt.Errorf("aws access key id is required")
	}
	if env.AwsRegion == "" {
		return fmt.Errorf("aws region is required")
	}
	if env.AwsSecretAccessKey == "" {
		return fmt.Errorf("aws secret access key is required")
	}
	if env.BlockConfirmationBuffer == 0 {
		return fmt.Errorf("block confirmation buffer is required")
	}
	if env.BlockQueryCutoff == 0 {
		return fmt.Errorf("block query cutoff is required")
	}
	if env.BlockQueryStart == 0 {
		return fmt.Errorf("block query start is required")
	}
	if env.BuilderHelperAddress == "" {
		return fmt.Errorf("builder helper address is required")
	}
	if env.BuilderKey == "" {
		return fmt.Errorf("builder key is required")
	}
	if env.BuilderPort == 0 {
		return fmt.Errorf("builder port is required")
	}
	if env.BuilderRewardsAddress == "" {
		return fmt.Errorf("builder rewards address is required")
	}
	if env.ChainOffset == 0 {
		return fmt.Errorf("chain offset is required")
	}
	if env.ConcurrentLimit == 0 {
		return fmt.Errorf("concurrent limit is required")
	}
	if env.HostChainId == 0 {
		return fmt.Errorf("host chain id is required")
	}
	if env.HostRpcUrl == "" {
		return fmt.Errorf("host RPC URL is required")
	}
	if env.OauthAudience == "" {
		return fmt.Errorf("oauth audience is required")
	}
	if env.OauthAuthenticateUrl == "" {
		return fmt.Errorf("oauth authenticate URL is required")
	}
	if env.OAuthClientId == "" {
		return fmt.Errorf("oauth client ID is required")
	}
	if env.OauthClientSecret == "" {
		return fmt.Errorf("oauth client secret is required")
	}
	if env.OauthIssuer == "" {
		return fmt.Errorf("oauth issuer is required")
	}
	if env.OauthTokenUrl == "" {
		return fmt.Errorf("oauth token URL is required")
	}
	if env.OtelExporterOtlpEndpoint == "" {
		return fmt.Errorf("otel exporter otlp endpoint is required")
	}
	if env.QuinceyUrl == "" {
		return fmt.Errorf("quincey URL is required")
	}
	if env.RollupBlockGasLimit == 0 {
		return fmt.Errorf("rollup block gas limit is required")
	}
	if env.RollupChainId == 0 {
		return fmt.Errorf("rollup chain id is required")
	}
	if env.RollupRpcUrl == "" {
		return fmt.Errorf("rollup RPC URL is required")
	}
	if env.RustLog == "" {
		return fmt.Errorf("rust log is required")
	}
	if env.SlotOffset == 0 {
		return fmt.Errorf("slot offset is required")
	}
	if env.StartTimestamp == 0 {
		return fmt.Errorf("start timestamp is required")
	}
	if env.SubmitViaCallData == "" {
		return fmt.Errorf("submit via call data is required")
	}
	if env.TargetSlotTime == 0 {
		return fmt.Errorf("target slot time is required")
	}
	if env.TxBroadcastUrls == "" {
		return fmt.Errorf("tx broadcast URLs is required")
	}
	if env.TxPoolCacheDuration == 0 {
		return fmt.Errorf("tx pool cache duration is required")
	}
	if env.TxPoolUrl == "" {
		return fmt.Errorf("tx pool URL is required")
	}
	if env.ZenithAddress == "" {
		return fmt.Errorf("zenith address is required")
	}

	return nil
}
