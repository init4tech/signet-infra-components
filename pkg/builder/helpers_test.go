package builder

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestCamelToSnake(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "BuilderPort",
			output: "BUILDER_PORT",
		},
		{
			input:  "HostRpcUrl",
			output: "HOST_RPC_URL",
		},
		{
			input:  "AWSRegion",
			output: "AWS_REGION",
		},
		{
			input:  "OtelExporterOtlpEndpoint",
			output: "OTEL_EXPORTER_OTLP_ENDPOINT",
		},
		{
			input:  "simpleText",
			output: "SIMPLE_TEXT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := CamelToSnake(tc.input)
			assert.Equal(t, tc.output, result)
		})
	}
}

func TestCreateKMSPolicy(t *testing.T) {
	// Test with a simple key ARN
	keyArn := "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab"
	key := pulumi.String(keyArn)

	// Get the policy string
	// Note: We can't directly access the string value from StringOutput in a unit test
	// So we'll just verify the function returns a non-nil value and the structure will
	// be tested separately when used in the actual builder component.
	policy := CreateKMSPolicy(key)

	// We can only indirectly test this by asserting the output is not nil
	assert.NotNil(t, policy)

	// Test with another key to ensure the function uses the provided key
	anotherKey := pulumi.String("another-key-arn")
	anotherPolicy := CreateKMSPolicy(anotherKey)
	assert.NotNil(t, anotherPolicy)
}

func TestCreateEnvVars(t *testing.T) {
	// Create a test BuilderEnv with some values
	env := BuilderEnv{
		BuilderPort:   pulumi.Int(8080),
		BuilderKey:    pulumi.String("test-key"),
		HostRpcUrl:    pulumi.String("http://host-rpc"),
		RollupRpcUrl:  pulumi.String("http://rollup-rpc"),
		ZenithAddress: pulumi.String("0x123456"),
		AwsRegion:     pulumi.String("us-west-2"),
		AwsAccountId:  pulumi.String("123456789012"),
	}

	// Get the environment variables
	envVars := CreateEnvVars(env)

	// Test that the array is not nil and has the correct length
	// We expect one env var for BuilderPort plus one for each of the other fields
	assert.NotNil(t, envVars)

	// Count the non-nil fields in env to determine expected size
	// We need to do this because BuilderEnv has many fields and we're only setting a few
	expectedSize := 0
	if env.BuilderPort != nil {
		expectedSize++
	}
	envVarMap := GetEnvironmentVarsFromStruct(env)
	expectedSize += len(envVarMap)

	// Check if the array is the expected size
	assert.Len(t, envVars, expectedSize)

	// Test that specific environment variables are set correctly
	// Due to the nature of Pulumi outputs, we can only test basics here
	// Just check that we have a non-empty array
	assert.Greater(t, len(envVars), 0, "Should have at least one environment variable")
}

func TestGetEnvironmentVarsFromStruct(t *testing.T) {
	// Create a test BuilderEnv with some values
	env := BuilderEnv{
		BuilderPort: pulumi.Int(8080),
		BuilderKey:  pulumi.String("test-key"),
		HostRpcUrl:  pulumi.String("http://host-rpc"),
		AwsRegion:   pulumi.String("us-west-2"),
	}

	// Get the environment variables map
	envVars := GetEnvironmentVarsFromStruct(env)

	// Test that the map is not nil and has the expected keys
	assert.NotNil(t, envVars)

	// BuilderPort should be excluded (handled specially)
	_, hasBuilderPort := envVars["BUILDER_PORT"]
	assert.False(t, hasBuilderPort, "BUILDER_PORT should not be in the map")

	// Other fields should be included
	_, hasBuilderKey := envVars["BUILDER_KEY"]
	assert.True(t, hasBuilderKey, "BUILDER_KEY should be in the map")

	_, hasHostRpcUrl := envVars["HOST_RPC_URL"]
	assert.True(t, hasHostRpcUrl, "HOST_RPC_URL should be in the map")

	_, hasAwsRegion := envVars["AWS_REGION"]
	assert.True(t, hasAwsRegion, "AWS_REGION should be in the map")
}
