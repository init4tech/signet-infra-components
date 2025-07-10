package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIAMStatementValidate tests the validation of IAMStatement
func TestIAMStatementValidate(t *testing.T) {
	// Test valid statement
	validStmt := IAMStatement{
		Effect: "Allow",
		Action: []string{"kms:Sign", "kms:GetPublicKey"},
	}
	err := validStmt.Validate()
	assert.NoError(t, err)

	// Test missing effect
	invalidStmt1 := IAMStatement{
		Action: []string{"kms:Sign"},
	}
	err = invalidStmt1.Validate()
	assert.Error(t, err)
	assert.Equal(t, "effect is required", err.Error())

	// Test missing action
	invalidStmt2 := IAMStatement{
		Effect: "Allow",
	}
	err = invalidStmt2.Validate()
	assert.Error(t, err)
	assert.Equal(t, "action is required", err.Error())
}

// TestIAMPolicyValidate tests the validation of IAMPolicy
func TestIAMPolicyValidate(t *testing.T) {
	// Test valid policy
	validPolicy := IAMPolicy{
		Version: "2012-10-17",
		Statement: []IAMStatement{
			{
				Effect: "Allow",
				Action: []string{"kms:Sign"},
			},
		},
	}
	err := validPolicy.Validate()
	assert.NoError(t, err)

	// Test missing version
	invalidPolicy1 := IAMPolicy{
		Statement: []IAMStatement{
			{
				Effect: "Allow",
				Action: []string{"kms:Sign"},
			},
		},
	}
	err = invalidPolicy1.Validate()
	assert.Error(t, err)
	assert.Equal(t, "version is required", err.Error())

	// Test missing statement
	invalidPolicy2 := IAMPolicy{
		Version: "2012-10-17",
	}
	err = invalidPolicy2.Validate()
	assert.Error(t, err)
	assert.Equal(t, "statement is required", err.Error())
}

// TestKMSStatementValidate tests the validation of KMSStatement
func TestKMSStatementValidate(t *testing.T) {
	// Test valid statement
	validStmt := KMSStatement{
		Effect:   "Allow",
		Action:   []string{"kms:Sign", "kms:GetPublicKey"},
		Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
	}
	err := validStmt.Validate()
	assert.NoError(t, err)

	// Test missing effect
	invalidStmt1 := KMSStatement{
		Action:   []string{"kms:Sign"},
		Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
	}
	err = invalidStmt1.Validate()
	assert.Error(t, err)
	assert.Equal(t, "effect is required", err.Error())

	// Test missing action
	invalidStmt2 := KMSStatement{
		Effect:   "Allow",
		Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
	}
	err = invalidStmt2.Validate()
	assert.Error(t, err)
	assert.Equal(t, "action is required", err.Error())

	// Test missing resource
	invalidStmt3 := KMSStatement{
		Effect: "Allow",
		Action: []string{"kms:Sign"},
	}
	err = invalidStmt3.Validate()
	assert.Error(t, err)
	assert.Equal(t, "resource is required", err.Error())
}

// TestKMSPolicyValidate tests the validation of KMSPolicy
func TestKMSPolicyValidate(t *testing.T) {
	// Test valid policy
	validPolicy := KMSPolicy{
		Version: "2012-10-17",
		Statement: []KMSStatement{
			{
				Effect:   "Allow",
				Action:   []string{"kms:Sign", "kms:GetPublicKey"},
				Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			},
		},
	}
	err := validPolicy.Validate()
	assert.NoError(t, err)

	// Test missing version
	invalidPolicy1 := KMSPolicy{
		Statement: []KMSStatement{
			{
				Effect:   "Allow",
				Action:   []string{"kms:Sign"},
				Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			},
		},
	}
	err = invalidPolicy1.Validate()
	assert.Error(t, err)
	assert.Equal(t, "version is required", err.Error())

	// Test missing statement
	invalidPolicy2 := KMSPolicy{
		Version: "2012-10-17",
	}
	err = invalidPolicy2.Validate()
	assert.Error(t, err)
	assert.Equal(t, "statement is required", err.Error())
}

// TestKMSPolicyToInternal tests the conversion from public to internal types
func TestKMSPolicyToInternal(t *testing.T) {
	// Create a public policy
	publicPolicy := KMSPolicy{
		Version: "2012-10-17",
		Statement: []KMSStatement{
			{
				Effect:   "Allow",
				Action:   []string{"kms:Sign", "kms:GetPublicKey"},
				Resource: "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			},
		},
	}

	// Convert to internal
	internalPolicy := publicPolicy.toInternal()

	// Verify the conversion
	assert.Equal(t, publicPolicy.Version, internalPolicy.Version)
	assert.Len(t, internalPolicy.Statement, 1)
	assert.Equal(t, publicPolicy.Statement[0].Effect, internalPolicy.Statement[0].Effect)
	assert.Equal(t, publicPolicy.Statement[0].Action, internalPolicy.Statement[0].Action)
	assert.NotNil(t, internalPolicy.Statement[0].Resource)
}
