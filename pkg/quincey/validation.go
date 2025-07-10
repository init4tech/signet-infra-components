package quincey

import (
	"fmt"
)

// Validate validates the QuinceyComponentArgs struct, ensuring all required fields are set
func (args *QuinceyComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if args.Image == "" {
		return fmt.Errorf("image is required")
	}
	if err := args.Env.Validate(); err != nil {
		return fmt.Errorf("env is invalid: %w", err)
	}
	if args.VirtualServiceHosts == nil {
		return fmt.Errorf("virtual service hosts is required")
	}
	if args.Port <= 0 {
		return fmt.Errorf("port must be a positive integer")
	}
	return nil
}

// Validate validates the QuinceyEnv struct, ensuring all required fields are set
func (env *QuinceyEnv) Validate() error {
	if env.QuinceyPort == "" {
		return fmt.Errorf("quincey port is required")
	}
	if env.QuinceyKeyId == "" {
		return fmt.Errorf("quincey key ID is required")
	}
	if env.AwsAccessKeyId == "" {
		return fmt.Errorf("AWS access key ID is required")
	}
	if env.AwsSecretAccessKey == "" {
		return fmt.Errorf("AWS secret access key is required")
	}
	if env.AwsDefaultRegion == "" {
		return fmt.Errorf("AWS default region is required")
	}
	if env.BlockQueryStart == "" {
		return fmt.Errorf("block query start is required")
	}
	if env.BlockQueryCutoff == "" {
		return fmt.Errorf("block query cutoff is required")
	}
	if env.ChainOffset == "" {
		return fmt.Errorf("chain offset is required")
	}
	if env.HostRpcUrl == "" {
		return fmt.Errorf("host RPC URL is required")
	}
	if env.OauthIssuer == "" {
		return fmt.Errorf("OAuth issuer is required")
	}
	if env.OauthJwksUri == "" {
		return fmt.Errorf("OAuth JWKS URI is required")
	}
	if env.QuinceyBuilders == "" {
		return fmt.Errorf("quincey builders is required")
	}
	return nil
}
