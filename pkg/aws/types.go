package aws

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMStatement represents a statement in an IAM policy
type IAMStatement struct {
	Sid       string `json:"sid,omitempty"`
	Effect    string `json:"effect"`
	Principal struct {
		Service []string `json:"Service"`
	} `json:"Principal"`
	Action []string `json:"Action"`
}

// IAMPolicy represents an IAM policy document
type IAMPolicy struct {
	Version   string         `json:"Version"`
	Statement []IAMStatement `json:"Statement"`
}

// KMSStatement represents a statement in a KMS policy
type KMSStatement struct {
	Effect   string             `json:"Effect"`
	Action   []string           `json:"Action"`
	Resource pulumi.StringInput `json:"Resource"`
}

// KMSPolicy represents a KMS policy document
type KMSPolicy struct {
	Version   string         `json:"Version"`
	Statement []KMSStatement `json:"Statement"`
}
