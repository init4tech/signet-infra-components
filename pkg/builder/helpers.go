package builder

import (
	"encoding/json"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateKMSPolicy creates a KMS policy for the builder service.
// Exported for testing.
func CreateKMSPolicy(key pulumi.StringInput) pulumi.StringOutput {
	policy := KMSPolicy{
		Version: "2012-10-17",
		Statement: []KMSStatement{
			{
				Effect: "Allow",
				Action: []string{
					"kms:Sign",
					"kms:GetPublicKey",
				},
				Resource: key,
			},
		},
	}

	// Convert to JSON string output
	return pulumi.All(key).ApplyT(func(_ []interface{}) (string, error) {
		jsonBytes, err := json.Marshal(policy)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}).(pulumi.StringOutput)
}
