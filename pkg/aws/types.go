// Package aws provides utilities for creating and managing AWS IAM resources
// for Kubernetes workloads running in EKS. It handles the creation of IAM roles,
// policies, and policy attachments that enable pod identity and KMS access.
package aws

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMStatement represents a statement in an IAM policy document.
// It defines a single permission statement that specifies what actions
// are allowed or denied on which resources.
type IAMStatement struct {
	// Sid is an optional identifier for the statement
	Sid string `json:"sid,omitempty"`
	// Effect specifies whether the statement allows or denies access
	Effect string `json:"effect"`
	// Principal specifies who is allowed or denied access
	Principal struct {
		// Service contains a list of AWS services that can assume this role
		Service []string `json:"Service"`
	} `json:"Principal"`
	// Action specifies the AWS actions that are allowed or denied
	Action []string `json:"Action"`
}

// IAMPolicy represents a complete IAM policy document.
// It contains a version and a list of statements that define the policy's permissions.
type IAMPolicy struct {
	// Version specifies the policy language version
	Version string `json:"Version"`
	// Statement contains the list of permission statements
	Statement []IAMStatement `json:"Statement"`
}

// KMSStatement represents a statement in a KMS policy document.
// It defines permissions specifically for AWS KMS operations.
type KMSStatement struct {
	// Effect specifies whether the statement allows or denies access
	Effect string `json:"Effect"`
	// Action specifies the KMS actions that are allowed or denied
	Action []string `json:"Action"`
	// Resource specifies the KMS key ARN that the permissions apply to
	Resource pulumi.StringInput `json:"Resource"`
}

// KMSPolicy represents a complete KMS policy document.
// It contains a version and a list of statements that define the KMS permissions.
type KMSPolicy struct {
	// Version specifies the policy language version
	Version string `json:"Version"`
	// Statement contains the list of KMS permission statements
	Statement []KMSStatement `json:"Statement"`
}
