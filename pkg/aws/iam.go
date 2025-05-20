package aws

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMResources contains AWS IAM resources created for a component
type IAMResources struct {
	Role             *iam.Role
	Policy           *iam.Policy
	PolicyAttachment *iam.RolePolicyAttachment
}

// CreateIAMResources creates IAM resources (role, policy, and policy attachment) for a component
func CreateIAMResources(
	ctx *pulumi.Context,
	name string,
	serviceName string,
	keyArn pulumi.StringInput,
	parent pulumi.Resource,
) (*IAMResources, error) {
	// Create IAM role
	assumeRolePolicy := IAMPolicy{
		Version: "2012-10-17",
		Statement: []IAMStatement{
			{
				Sid:    "AllowEksAuthToAssumeRoleForPodIdentity",
				Effect: "Allow",
				Principal: struct {
					Service []string `json:"Service"`
				}{
					Service: []string{
						"pods.eks.amazonaws.com",
						"ec2.amazonaws.com",
					},
				},
				Action: []string{
					"sts:AssumeRole",
					"sts:TagSession",
				},
			},
		},
	}

	assumeRolePolicyJSON, err := json.Marshal(assumeRolePolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal assume role policy: %w", err)
	}

	role, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
		Description:      pulumi.String(fmt.Sprintf("Role for %s pod to assume", serviceName)),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(fmt.Sprintf("%s-role", name)),
		},
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM role: %w", err)
	}

	// Create KMS policy
	policyJSON := CreateKMSPolicy(keyArn)

	policy, err := iam.NewPolicy(ctx, fmt.Sprintf("%s-policy", name), &iam.PolicyArgs{
		Policy: policyJSON,
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM policy: %w", err)
	}

	// Attach policy to role
	policyAttachment, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-role-policy-attachment", name), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: policy.Arn,
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, fmt.Errorf("failed to attach policy to role: %w", err)
	}

	return &IAMResources{
		Role:             role,
		Policy:           policy,
		PolicyAttachment: policyAttachment,
	}, nil
}

// CreateKMSPolicy creates a KMS policy for the given key
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
