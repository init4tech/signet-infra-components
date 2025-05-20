package utils

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateResourceLabels creates a consistent set of Kubernetes labels for resources
func CreateResourceLabels(app, name, partOf string, additionalLabels pulumi.StringMap) pulumi.StringMap {
	labels := pulumi.StringMap{
		"app.kubernetes.io/name":    pulumi.String(name),
		"app.kubernetes.io/part-of": pulumi.String(partOf),
	}

	// Merge additional labels if provided
	if additionalLabels != nil {
		for k, v := range additionalLabels {
			labels[k] = v
		}
	}

	return labels
}
