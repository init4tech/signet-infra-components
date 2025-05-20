package signet_node

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HostZenithAddress", "HOST_ZENITH_ADDRESS"},
		{"RpcPort", "RPC_PORT"},
		{"SignetChainId", "SIGNET_CHAIN_ID"},
		{"BaseFeeRecipient", "BASE_FEE_RECIPIENT"},
		{"HostStartTimestamp", "HOST_START_TIMESTAMP"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := CamelToSnake(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetEnvironmentVarsFromStruct(t *testing.T) {
	env := SignetNodeEnv{
		HostZenithAddress: pulumi.String("0x123"),
		RpcPort:           pulumi.String("8545"),
		WsRpcPort:         pulumi.String("8546"),
	}

	result := GetEnvironmentVarsFromStruct(env)

	assert.NotNil(t, result["HOST_ZENITH_ADDRESS"])
	assert.NotNil(t, result["RPC_PORT"])
	assert.NotNil(t, result["WS_RPC_PORT"])
	assert.Equal(t, 3, len(result))
}

func TestCreateResourceLabels(t *testing.T) {
	name := "test-resource"
	labels := CreateResourceLabels(name)

	assert.Equal(t, pulumi.String(name), labels["app"])
	assert.Equal(t, pulumi.String(name), labels["app.kubernetes.io/name"])
	assert.Equal(t, pulumi.String(name), labels["app.kubernetes.io/part-of"])
}

func TestGetResourceName(t *testing.T) {
	tests := []struct {
		baseName     string
		resourceType string
		expected     string
	}{
		{"signet-node", "service", "signet-node-service"},
		{"lighthouse", "statefulset", "lighthouse-set"},
		{"execution", "configmap", "execution-configmap"},
		{"data", "pvc", "data-data"},
		{"jwt", "secret", "jwt-secret"},
	}

	for _, test := range tests {
		t.Run(test.resourceType, func(t *testing.T) {
			result := GetResourceName(test.baseName, test.resourceType)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestResourceRequirements(t *testing.T) {
	// Test with default values
	resources := NewResourceRequirements("", "", "", "")

	// Check that resources is not nil
	assert.NotNil(t, resources)

	// Verify the structure indirectly by checking key presence
	stringMap := resources.Limits.(pulumi.StringMap)
	assert.Contains(t, stringMap, "cpu")
	assert.Contains(t, stringMap, "memory")
	assert.Equal(t, pulumi.String(DefaultCPULimit), stringMap["cpu"])
	assert.Equal(t, pulumi.String(DefaultMemoryLimit), stringMap["memory"])

	requestsMap := resources.Requests.(pulumi.StringMap)
	assert.Contains(t, requestsMap, "cpu")
	assert.Contains(t, requestsMap, "memory")
	assert.Equal(t, pulumi.String(DefaultCPURequest), requestsMap["cpu"])
	assert.Equal(t, pulumi.String(DefaultMemoryRequest), requestsMap["memory"])

	// Test with custom values
	customResources := NewResourceRequirements("4", "32Gi", "2", "16Gi")

	customLimitsMap := customResources.Limits.(pulumi.StringMap)
	assert.Equal(t, pulumi.String("4"), customLimitsMap["cpu"])
	assert.Equal(t, pulumi.String("32Gi"), customLimitsMap["memory"])

	customRequestsMap := customResources.Requests.(pulumi.StringMap)
	assert.Equal(t, pulumi.String("2"), customRequestsMap["cpu"])
	assert.Equal(t, pulumi.String("16Gi"), customRequestsMap["memory"])
}
