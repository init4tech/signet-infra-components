package builder

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NewAppLabels creates a new AppLabels instance with consistent Kubernetes labels
func NewAppLabels(app, name, partOf string, additionalLabels pulumi.StringMap) AppLabels {
	return AppLabels{
		Labels: utils.CreateResourceLabels(app, name, partOf, additionalLabels),
	}
}
