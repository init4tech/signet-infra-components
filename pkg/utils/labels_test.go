package utils

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestCreateResourceLabels(t *testing.T) {
	// Test with standard labels
	appName := "test-app"
	resourceName := "test-resource"
	partOf := "test-system"

	labels := CreateResourceLabels(appName, resourceName, partOf, nil)

	// The "app" label is no longer included in standard labels per requirements
	assert.Equal(t, pulumi.String(resourceName), labels["app.kubernetes.io/name"])
	assert.Equal(t, pulumi.String(partOf), labels["app.kubernetes.io/part-of"])

	// Test with additional labels
	additionalLabels := pulumi.StringMap{
		"app":          pulumi.String(appName), // Now app should be explicitly provided as an additional label
		"environment":  pulumi.String("production"),
		"tier":         pulumi.String("backend"),
		"custom-label": pulumi.String("custom-value"),
	}

	mergedLabels := CreateResourceLabels(appName, resourceName, partOf, additionalLabels)

	// Verify standard labels
	assert.Equal(t, pulumi.String(appName), mergedLabels["app"]) // Now comes from additional labels
	assert.Equal(t, pulumi.String(resourceName), mergedLabels["app.kubernetes.io/name"])
	assert.Equal(t, pulumi.String(partOf), mergedLabels["app.kubernetes.io/part-of"])

	// Verify additional labels were merged
	assert.Equal(t, pulumi.String("production"), mergedLabels["environment"])
	assert.Equal(t, pulumi.String("backend"), mergedLabels["tier"])
	assert.Equal(t, pulumi.String("custom-value"), mergedLabels["custom-label"])

	// Test override behavior
	overrideLabels := pulumi.StringMap{
		"app": pulumi.String("override-app"),
	}

	overriddenLabels := CreateResourceLabels(appName, resourceName, partOf, overrideLabels)

	// Verify the override behavior
	assert.Equal(t, pulumi.String("override-app"), overriddenLabels["app"]) // Comes from overrideLabels
	assert.Equal(t, pulumi.String(resourceName), overriddenLabels["app.kubernetes.io/name"])
	assert.Equal(t, pulumi.String(partOf), overriddenLabels["app.kubernetes.io/part-of"])
}
