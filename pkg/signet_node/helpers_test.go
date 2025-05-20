package signet_node

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

// TestCreateResourceLabels moved to pkg/utils/labels_test.go since the function was moved there

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
