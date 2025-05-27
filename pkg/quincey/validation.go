package quincey

import (
	"fmt"
)

// Validate validates the QuinceyComponentArgs struct, ensuring all required fields are set
func (args *QuinceyComponentArgs) Validate() error {
	if args.Namespace == nil {
		return fmt.Errorf("namespace is required")
	}
	if args.Image == nil {
		return fmt.Errorf("image is required")
	}
	if args.Env.Validate() != nil {
		return fmt.Errorf("env is invalid: %w", args.Env.Validate())
	}
	if args.VirtualServiceHosts == nil {
		return fmt.Errorf("virtual service hosts is required")
	}
	return nil
}

// Validate validates the QuinceyEnv struct, ensuring all required fields are set
func (env *QuinceyEnv) Validate() error {
	if env.QuinceyPort == nil {
		return fmt.Errorf("quincey port is required")
	}
	if env.QuinceyKeyId == nil {
		return fmt.Errorf("quincey key ID is required")
	}
	if env.AwsAccessKeyId == nil {
		return fmt.Errorf("AWS access key ID is required")
	}
	if env.AwsSecretAccessKey == nil {
		return fmt.Errorf("AWS secret access key is required")
	}
	if env.AwsDefaultRegion == nil {
		return fmt.Errorf("AWS default region is required")
	}
	if env.BlockQueryStart == nil {
		return fmt.Errorf("block query start is required")
	}
	if env.BlockQueryCutoff == nil {
		return fmt.Errorf("block query cutoff is required")
	}
	if env.ChainOffset == nil {
		return fmt.Errorf("chain offset is required")
	}
	if env.HostRpcUrl == nil {
		return fmt.Errorf("host RPC URL is required")
	}
	if env.OauthIssuer == nil {
		return fmt.Errorf("OAuth issuer is required")
	}
	if env.OauthJwksUri == nil {
		return fmt.Errorf("OAuth JWKS URI is required")
	}
	if env.QuinceyBuilders == nil {
		return fmt.Errorf("quincey builders is required")
	}
	return nil
}
