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
	if env.BuilderPort == nil {
		return fmt.Errorf("builder port is required")
	}
	if env.BuilderKey == nil {
		return fmt.Errorf("builder key is required")
	}
	if env.HostRpcUrl == nil {
		return fmt.Errorf("host RPC URL is required")
	}
	if env.RollupRpcUrl == nil {
		return fmt.Errorf("rollup RPC URL is required")
	}
	if env.ZenithAddress == nil {
		return fmt.Errorf("zenith address is required")
	}
	if env.QuinceyUrl == nil {
		return fmt.Errorf("quincey URL is required")
	}
	if env.OtelExporterOtlpEndpoint == nil {
		return fmt.Errorf("otel exporter otlp endpoint is required")
	}
	if env.OauthAudience == nil {
		return fmt.Errorf("oauth audience is required")
	}
	if env.OauthAuthenticateUrl == nil {
		return fmt.Errorf("oauth authenticate URL is required")
	}
	if env.OAuthClientId == nil {
		return fmt.Errorf("oauth client ID is required")
	}
	if env.OauthClientSecret == nil {
		return fmt.Errorf("oauth client secret is required")
	}
	if env.OauthIssuer == nil {
		return fmt.Errorf("oauth issuer is required")
	}
	if env.OauthTokenUrl == nil {
		return fmt.Errorf("oauth token URL is required")
	}
	if env.RustLog == nil {
		return fmt.Errorf("rust log is required")
	}
	if env.SlotOffset == nil {
		return fmt.Errorf("slot offset is required")
	}
	if env.StartTimestamp == nil {
		return fmt.Errorf("start timestamp is required")
	}
	if env.SubmitViaCallData == nil {
		return fmt.Errorf("submit via call data is required")
	}
	if env.TargetSlotTime == nil {
		return fmt.Errorf("target slot time is required")
	}
	if env.TxBroadcastUrls == nil {
		return fmt.Errorf("tx broadcast URLs is required")
	}
	if env.TxPoolCacheDuration == nil {
		return fmt.Errorf("tx pool cache duration is required")
	}
	if env.TxPoolUrl == nil {
		return fmt.Errorf("tx pool URL is required")
	}
	if env.ZenithAddress == nil {
		return fmt.Errorf("zenith address is required")
	}

	return nil
}
