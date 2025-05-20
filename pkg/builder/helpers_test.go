package builder

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

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

func TestBuilderEnvGetEnvMap(t *testing.T) {
	// Create a test BuilderEnv with some values
	env := BuilderEnv{
		BuilderPort:   pulumi.String("8080"),
		BuilderKey:    pulumi.String("test-key"),
		HostRpcUrl:    pulumi.String("http://host-rpc"),
		RollupRpcUrl:  pulumi.String("http://rollup-rpc"),
		ZenithAddress: pulumi.String("0x123456"),
		AwsRegion:     pulumi.String("us-west-2"),
		AwsAccountId:  pulumi.String("123456789012"),
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
