package quincey

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

// TestCreateContainer tests the container creation logic
func TestCreateContainer(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		args     *quinceyComponentArgsInternal
		port     pulumi.IntOutput
		wantName string
	}{
		{
			name: "valid container config",
			args: &quinceyComponentArgsInternal{
				Image: pulumi.String("quincey:test"),
				Env: quinceyEnvInternal{
					QuinceyPort:              pulumi.String("8080"),
					QuinceyKeyId:             pulumi.String("test-key-id"),
					AwsAccessKeyId:           pulumi.String("test-aws-key"),
					AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
					AwsDefaultRegion:         pulumi.String("us-west-2"),
					BlockQueryStart:          pulumi.String("1"),
					BlockQueryCutoff:         pulumi.String("1000"),
					ChainOffset:              pulumi.String("10"),
					HostRpcUrl:               pulumi.String("http://test-rpc:8545"),
					OauthIssuer:              pulumi.String("https://test-issuer"),
					OauthJwksUri:             pulumi.String("https://test-jwks"),
					QuinceyBuilders:          pulumi.String("test-builder"),
					OtelExporterOtlpEndpoint: pulumi.String("http://otel:4317"),
					OtelExporterOtlpProtocol: pulumi.String("grpc"),
					RustLog:                  pulumi.String("info"),
				},
			},
			port:     pulumi.Int(8080).ToIntOutput(),
			wantName: ServiceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := createContainer(tt.args, tt.port)
			assert.NotNil(t, container)
			assert.Equal(t, pulumi.String(tt.wantName), container.Name)
			assert.Equal(t, tt.args.Image, container.Image)
			assert.NotNil(t, container.EnvFrom)
		})
	}
}

// TestQuinceyEnv tests the environment variable handling
func TestQuinceyEnv(t *testing.T) {
	env := QuinceyEnv{
		QuinceyPort:              "8080",
		QuinceyKeyId:             "test-key-id",
		AwsAccessKeyId:           "test-aws-key",
		AwsSecretAccessKey:       "test-aws-secret",
		AwsDefaultRegion:         "us-west-2",
		BlockQueryStart:          "1",
		BlockQueryCutoff:         "1000",
		ChainOffset:              "10",
		HostRpcUrl:               "http://test-rpc:8545",
		OauthIssuer:              "https://test-issuer",
		OauthJwksUri:             "https://test-jwks",
		QuinceyBuilders:          "test-builder",
		OtelExporterOtlpEndpoint: "http://otel:4317",
		OtelExporterOtlpProtocol: "grpc",
		RustLog:                  "info",
	}

	// Test that the environment variables are properly set
	assert.Equal(t, "8080", env.QuinceyPort)
	assert.Equal(t, "test-key-id", env.QuinceyKeyId)
	assert.Equal(t, "test-aws-key", env.AwsAccessKeyId)
	assert.Equal(t, "test-aws-secret", env.AwsSecretAccessKey)
	assert.Equal(t, "us-west-2", env.AwsDefaultRegion)
	assert.Equal(t, "1", env.BlockQueryStart)
	assert.Equal(t, "1000", env.BlockQueryCutoff)
	assert.Equal(t, "10", env.ChainOffset)
	assert.Equal(t, "http://test-rpc:8545", env.HostRpcUrl)
	assert.Equal(t, "https://test-issuer", env.OauthIssuer)
	assert.Equal(t, "https://test-jwks", env.OauthJwksUri)
	assert.Equal(t, "test-builder", env.QuinceyBuilders)
	assert.Equal(t, "http://otel:4317", env.OtelExporterOtlpEndpoint)
	assert.Equal(t, "grpc", env.OtelExporterOtlpProtocol)
	assert.Equal(t, "info", env.RustLog)
}

// TestConstants tests the package constants
func TestConstants(t *testing.T) {
	assert.Equal(t, "quincey-server", ServiceName)
	assert.Equal(t, "quincey-server", AppLabel)
	assert.Equal(t, 9000, DefaultMetricsPort)
	assert.Equal(t, "quincey", ComponentName)
}

// TestQuinceyEnvValidation tests the validation of required fields in QuinceyEnv
func TestQuinceyEnvValidation(t *testing.T) {
	tests := []struct {
		name        string
		env         QuinceyEnv
		expectError bool
	}{
		{
			name: "valid env with all required fields",
			env: QuinceyEnv{
				QuinceyPort:              "8080",
				QuinceyKeyId:             "test-key-id",
				AwsAccessKeyId:           "test-aws-key",
				AwsSecretAccessKey:       "test-aws-secret",
				AwsDefaultRegion:         "us-west-2",
				BlockQueryStart:          "1",
				BlockQueryCutoff:         "1000",
				ChainOffset:              "10",
				HostRpcUrl:               "http://test-rpc:8545",
				OauthIssuer:              "https://test-issuer",
				OauthJwksUri:             "https://test-jwks",
				QuinceyBuilders:          "test-builder",
				OtelExporterOtlpEndpoint: "http://otel:4317",
				OtelExporterOtlpProtocol: "grpc",
				RustLog:                  "info",
			},
			expectError: false,
		},
		{
			name: "missing required field QuinceyPort",
			env: QuinceyEnv{
				QuinceyKeyId:             "test-key-id",
				AwsAccessKeyId:           "test-aws-key",
				AwsSecretAccessKey:       "test-aws-secret",
				AwsDefaultRegion:         "us-west-2",
				BlockQueryStart:          "1",
				BlockQueryCutoff:         "1000",
				ChainOffset:              "10",
				HostRpcUrl:               "http://test-rpc:8545",
				OauthIssuer:              "https://test-issuer",
				OauthJwksUri:             "https://test-jwks",
				QuinceyBuilders:          "test-builder",
				OtelExporterOtlpEndpoint: "http://otel:4317",
				OtelExporterOtlpProtocol: "grpc",
				RustLog:                  "info",
			},
			expectError: true,
		},
		{
			name: "missing required field QuinceyKeyId",
			env: QuinceyEnv{
				QuinceyPort:              "8080",
				AwsAccessKeyId:           "test-aws-key",
				AwsSecretAccessKey:       "test-aws-secret",
				AwsDefaultRegion:         "us-west-2",
				BlockQueryStart:          "1",
				BlockQueryCutoff:         "1000",
				ChainOffset:              "10",
				HostRpcUrl:               "http://test-rpc:8545",
				OauthIssuer:              "https://test-issuer",
				OauthJwksUri:             "https://test-jwks",
				QuinceyBuilders:          "test-builder",
				OtelExporterOtlpEndpoint: "http://otel:4317",
				OtelExporterOtlpProtocol: "grpc",
				RustLog:                  "info",
			},
			expectError: true,
		},
		{
			name: "missing optional fields",
			env: QuinceyEnv{
				QuinceyPort:        "8080",
				QuinceyKeyId:       "test-key-id",
				AwsAccessKeyId:     "test-aws-key",
				AwsSecretAccessKey: "test-aws-secret",
				AwsDefaultRegion:   "us-west-2",
				BlockQueryStart:    "1",
				BlockQueryCutoff:   "1000",
				ChainOffset:        "10",
				HostRpcUrl:         "http://test-rpc:8545",
				OauthIssuer:        "https://test-issuer",
				OauthJwksUri:       "https://test-jwks",
				QuinceyBuilders:    "test-builder",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.env.Validate()
			if tt.expectError {
				assert.Error(t, err, "Expected validation error for missing required field")
			} else {
				assert.NoError(t, err, "Expected no validation error for valid env")
			}
		})
	}
}
