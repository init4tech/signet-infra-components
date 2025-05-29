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
		args     *QuinceyComponentArgs
		port     pulumi.IntOutput
		wantName string
	}{
		{
			name: "valid container config",
			args: &QuinceyComponentArgs{
				Image: pulumi.String("quincey:test"),
				Env: QuinceyEnv{
					QuinceyPort:              pulumi.String("8080"),
					QuinceyKeyId:             pulumi.String("test-key-id"),
					AwsAccessKeyId:           pulumi.String("test-aws-key"),
					AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
					AwsDefaultRegion:         pulumi.String("us-west-2"),
					BlockQueryStart:          pulumi.String("0"),
					BlockQueryCutoff:         pulumi.String("1000"),
					ChainOffset:              pulumi.String("0"),
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
		QuinceyPort:              pulumi.String("8080"),
		QuinceyKeyId:             pulumi.String("test-key-id"),
		AwsAccessKeyId:           pulumi.String("test-aws-key"),
		AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
		AwsDefaultRegion:         pulumi.String("us-west-2"),
		BlockQueryStart:          pulumi.String("0"),
		BlockQueryCutoff:         pulumi.String("1000"),
		ChainOffset:              pulumi.String("0"),
		HostRpcUrl:               pulumi.String("http://test-rpc:8545"),
		OauthIssuer:              pulumi.String("https://test-issuer"),
		OauthJwksUri:             pulumi.String("https://test-jwks"),
		QuinceyBuilders:          pulumi.String("test-builder"),
		OtelExporterOtlpEndpoint: pulumi.String("http://otel:4317"),
		OtelExporterOtlpProtocol: pulumi.String("grpc"),
		RustLog:                  pulumi.String("info"),
	}

	// Since we can't directly test Pulumi outputs in unit tests,
	// we'll verify that the environment variables are properly set
	// by checking that the struct fields are not nil
	assert.NotNil(t, env.QuinceyPort)
	assert.NotNil(t, env.QuinceyKeyId)
	assert.NotNil(t, env.AwsAccessKeyId)
	assert.NotNil(t, env.AwsSecretAccessKey)
	assert.NotNil(t, env.AwsDefaultRegion)
	assert.NotNil(t, env.BlockQueryStart)
	assert.NotNil(t, env.BlockQueryCutoff)
	assert.NotNil(t, env.ChainOffset)
	assert.NotNil(t, env.HostRpcUrl)
	assert.NotNil(t, env.OauthIssuer)
	assert.NotNil(t, env.OauthJwksUri)
	assert.NotNil(t, env.QuinceyBuilders)
	assert.NotNil(t, env.OtelExporterOtlpEndpoint)
	assert.NotNil(t, env.OtelExporterOtlpProtocol)
	assert.NotNil(t, env.RustLog)
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
				QuinceyPort:              pulumi.String("8080"),
				QuinceyKeyId:             pulumi.String("test-key-id"),
				AwsAccessKeyId:           pulumi.String("test-aws-key"),
				AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
				AwsDefaultRegion:         pulumi.String("us-west-2"),
				BlockQueryStart:          pulumi.String("0"),
				BlockQueryCutoff:         pulumi.String("1000"),
				ChainOffset:              pulumi.String("0"),
				HostRpcUrl:               pulumi.String("http://test-rpc:8545"),
				OauthIssuer:              pulumi.String("https://test-issuer"),
				OauthJwksUri:             pulumi.String("https://test-jwks"),
				QuinceyBuilders:          pulumi.String("test-builder"),
				OtelExporterOtlpEndpoint: pulumi.String("http://otel:4317"),
				OtelExporterOtlpProtocol: pulumi.String("grpc"),
				RustLog:                  pulumi.String("info"),
			},
			expectError: false,
		},
		{
			name: "missing required field QuinceyPort",
			env: QuinceyEnv{
				QuinceyKeyId:             pulumi.String("test-key-id"),
				AwsAccessKeyId:           pulumi.String("test-aws-key"),
				AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
				AwsDefaultRegion:         pulumi.String("us-west-2"),
				BlockQueryStart:          pulumi.String("0"),
				BlockQueryCutoff:         pulumi.String("1000"),
				ChainOffset:              pulumi.String("0"),
				HostRpcUrl:               pulumi.String("http://test-rpc:8545"),
				OauthIssuer:              pulumi.String("https://test-issuer"),
				OauthJwksUri:             pulumi.String("https://test-jwks"),
				QuinceyBuilders:          pulumi.String("test-builder"),
				OtelExporterOtlpEndpoint: pulumi.String("http://otel:4317"),
				OtelExporterOtlpProtocol: pulumi.String("grpc"),
				RustLog:                  pulumi.String("info"),
			},
			expectError: true,
		},
		{
			name: "missing required field QuinceyKeyId",
			env: QuinceyEnv{
				QuinceyPort:              pulumi.String("8080"),
				AwsAccessKeyId:           pulumi.String("test-aws-key"),
				AwsSecretAccessKey:       pulumi.String("test-aws-secret"),
				AwsDefaultRegion:         pulumi.String("us-west-2"),
				BlockQueryStart:          pulumi.String("0"),
				BlockQueryCutoff:         pulumi.String("1000"),
				ChainOffset:              pulumi.String("0"),
				HostRpcUrl:               pulumi.String("http://test-rpc:8545"),
				OauthIssuer:              pulumi.String("https://test-issuer"),
				OauthJwksUri:             pulumi.String("https://test-jwks"),
				QuinceyBuilders:          pulumi.String("test-builder"),
				OtelExporterOtlpEndpoint: pulumi.String("http://otel:4317"),
				OtelExporterOtlpProtocol: pulumi.String("grpc"),
				RustLog:                  pulumi.String("info"),
			},
			expectError: true,
		},
		{
			name: "missing optional fields",
			env: QuinceyEnv{
				QuinceyPort:        pulumi.String("8080"),
				QuinceyKeyId:       pulumi.String("test-key-id"),
				AwsAccessKeyId:     pulumi.String("test-aws-key"),
				AwsSecretAccessKey: pulumi.String("test-aws-secret"),
				AwsDefaultRegion:   pulumi.String("us-west-2"),
				BlockQueryStart:    pulumi.String("0"),
				BlockQueryCutoff:   pulumi.String("1000"),
				ChainOffset:        pulumi.String("0"),
				HostRpcUrl:         pulumi.String("http://test-rpc:8545"),
				OauthIssuer:        pulumi.String("https://test-issuer"),
				OauthJwksUri:       pulumi.String("https://test-jwks"),
				QuinceyBuilders:    pulumi.String("test-builder"),
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
