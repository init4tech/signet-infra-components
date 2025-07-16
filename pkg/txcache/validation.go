package txcache

import "fmt"

// Validate validates the TxCacheComponentArgs struct
func (args TxCacheComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if args.Image == "" {
		return fmt.Errorf("image is required")
	}

	if args.Port == 0 {
		return fmt.Errorf("port is required")
	}

	if args.OauthIssuer == "" {
		return fmt.Errorf("oauthIssuer is required")
	}

	if args.OauthJwksUri == "" {
		return fmt.Errorf("oauthJwksUri is required")
	}

	if err := validateEnv(args.Env); err != nil {
		return err
	}

	return nil
}

// validateEnv validates the TxCacheEnv struct
// OtelExporterOtlpProtocol and OtelExporterOtlpEndpoint are optional
// so no validation needed for them
func validateEnv(env TxCacheEnv) error {
	if env.HttpPort == "" {
		return fmt.Errorf("httpPort is required")
	}

	if env.AwsAccessKeyId == "" {
		return fmt.Errorf("awsAccessKeyId is required")
	}

	if env.AwsSecretAccessKey == "" {
		return fmt.Errorf("awsSecretAccessKey is required")
	}

	if env.AwsRegion == "" {
		return fmt.Errorf("awsRegion is required")
	}

	if env.RustLog == "" {
		return fmt.Errorf("rustLog is required")
	}

	if env.BlockQueryStart == "" {
		return fmt.Errorf("blockQueryStart is required")
	}

	if env.BlockQueryCutoff == "" {
		return fmt.Errorf("blockQueryCutoff is required")
	}

	if env.SlotOffset == "" {
		return fmt.Errorf("slotOffset is required")
	}

	if env.ExpirationTimestampOffset == "" {
		return fmt.Errorf("expirationTimestampOffset is required")
	}

	if env.NetworkName == "" {
		return fmt.Errorf("networkName is required")
	}

	if env.Builders == "" {
		return fmt.Errorf("builders is required")
	}

	if env.SlotDuration == "" {
		return fmt.Errorf("slotDuration is required")
	}

	if env.StartTimestamp == "" {
		return fmt.Errorf("startTimestamp is required")
	}

	return nil
}
