package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderEnvGetEnvMap(t *testing.T) {
	// Create a test BuilderEnv with some values
	env := BuilderEnv{
		BuilderPort:   8080,
		BuilderKey:    "test-key",
		HostRpcUrl:    "http://host-rpc",
		RollupRpcUrl:  "http://rollup-rpc",
		ZenithAddress: "0x123456",
		AwsRegion:     "us-west-2",
		AwsAccountId:  "123456789012",
	}

	// Get the environment variables map
	envMap := env.GetEnvMap()

	// Test that the map is not nil
	assert.NotNil(t, envMap)

	// Check that specific environment variables are in the map
	_, hasBuilderPort := envMap["BUILDER_PORT"]
	assert.True(t, hasBuilderPort, "BUILDER_PORT should be in the map")

	_, hasBuilderKey := envMap["BUILDER_KEY"]
	assert.True(t, hasBuilderKey, "BUILDER_KEY should be in the map")

	_, hasHostRpcUrl := envMap["HOST_RPC_URL"]
	assert.True(t, hasHostRpcUrl, "HOST_RPC_URL should be in the map")

	_, hasAwsRegion := envMap["AWS_REGION"]
	assert.True(t, hasAwsRegion, "AWS_REGION should be in the map")
}
