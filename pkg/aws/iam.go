// Package aws provides utilities for creating and managing AWS IAM resources
// for Kubernetes workloads running in EKS. It handles the creation of IAM roles,
// policies, and policy attachments that enable pod identity and KMS access.
//
// The package is designed to work with Pulumi and provides a high-level interface
// for managing AWS IAM resources, particularly for services that need to interact
// with AWS KMS for cryptographic operations.
package aws

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMResources contains AWS IAM resources created for a component.
// These resources work together to provide the necessary permissions
// for a Kubernetes workload to interact with AWS services.
type IAMResources struct {
	// Role is the IAM role that can be assumed by the Kubernetes workload
	Role *iam.Role
	// Policy defines the permissions granted to the role
	Policy *iam.Policy
	// PolicyAttachment connects the policy to the role
	PolicyAttachment *iam.RolePolicyAttachment
}

// CreateIAMResources creates IAM resources (role, policy, and policy attachment) for a component.
// It sets up the necessary permissions for a Kubernetes workload to interact with AWS KMS.
//
// Parameters:
//   - ctx: The Pulumi context
//   - name: Base name for the IAM resources
//   - serviceName: Name of the service that will use these IAM resources
//   - keyArn: ARN of the KMS key that the service needs to access
//   - parent: Parent Pulumi resource for dependency tracking
//
// Returns:
//   - *IAMResources: The created IAM resources
//   - error: Any error that occurred during creation
//
// Example:
//
//	resources, err := CreateIAMResources(ctx, "my-service", "my-service", keyArn, parent)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to create IAM resources: %w", err)
//	}
func CreateIAMResources(
	ctx *pulumi.Context,
	name string,
	serviceName string,
	keyArn pulumi.StringInput,
	parent pulumi.Resource,
) (*IAMResources, error) {
	// Create IAM role with assume role policy for EKS pod identity
	assumeRolePolicy := IAMPolicy{
		Version: IAMPolicyVersion,
		Statement: []IAMStatement{
			{
				Sid:    EKSAssumeRoleStatementSid,
				Effect: EffectAllow,
				Principal: struct {
					Service []string `json:"Service"`
				}{
					Service: []string{
						EKSPodsService,
						EC2Service,
					},
				},
				Action: []string{
					STSAssumeRoleAction,
					STSTagSessionAction,
				},
			},
		},
	}

	assumeRolePolicyJSON, err := json.Marshal(assumeRolePolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal assume role policy: %w", err)
	}

	role, err := iam.NewRole(ctx, fmt.Sprintf("%s%s", name, RoleSuffix), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(assumeRolePolicyJSON),
		Description:      pulumi.String(fmt.Sprintf("Role for %s pod to assume", serviceName)),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(fmt.Sprintf("%s%s", name, RoleSuffix)),
		},
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM role: %w", err)
	}

	// Create KMS policy for the specified key
	policyJSON := CreateKMSPolicy(keyArn)

	policy, err := iam.NewPolicy(ctx, fmt.Sprintf("%s%s", name, PolicySuffix), &iam.PolicyArgs{
		Policy: policyJSON,
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM policy: %w", err)
	}

	// Attach the KMS policy to the role
	policyAttachment, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s%s", name, RolePolicyAttachmentSuffix), &iam.RolePolicyAttachmentArgs{
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

// CreateKMSPolicy creates a KMS policy document that grants permissions to sign messages
// and retrieve public keys using the specified KMS key.
//
// Parameters:
//   - key: The ARN of the KMS key to create the policy for
//
// Returns:
//   - pulumi.StringOutput: A Pulumi output containing the JSON policy document
//
// The policy grants the following permissions:
//   - kms:Sign: Allows signing messages using the KMS key
//   - kms:GetPublicKey: Allows retrieving the public key associated with the KMS key
func CreateKMSPolicy(key pulumi.StringInput) pulumi.StringOutput {
	// Create internal policy with Pulumi types
	policy := kmsPolicyInternal{
		Version: IAMPolicyVersion,
		Statement: []kmsStatementInternal{
			{
				Effect: EffectAllow,
				Action: []string{
					KMSSignAction,
					KMSGetPublicKeyAction,
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

// CreateKMSPolicyFromPublic creates a KMS policy document from public types
// and converts it to internal types for use with Pulumi.
//
// Parameters:
//   - policy: The public KMSPolicy struct
//
// Returns:
//   - pulumi.StringOutput: A Pulumi output containing the JSON policy document
func CreateKMSPolicyFromPublic(policy KMSPolicy) (pulumi.StringOutput, error) {
	if err := policy.Validate(); err != nil {
		return pulumi.StringOutput{}, fmt.Errorf("invalid KMS policy: %w", err)
	}

	// Convert to internal format
	internalPolicy := policy.toInternal()

	// Convert to JSON string output
	// Since we're working with internal types that have Pulumi inputs,
	// we need to handle the conversion carefully
	statements := internalPolicy.Statement
	if len(statements) == 0 {
		return pulumi.StringOutput{}, fmt.Errorf("policy must have at least one statement")
	}

	// For now, we'll create a simple policy with the first statement
	// In a more complex scenario, you might need to handle multiple statements differently
	firstStmt := statements[0]

	return pulumi.All(firstStmt.Resource).ApplyT(func(values []interface{}) (string, error) {
		resource := values[0].(string)
		// Create a temporary struct for marshaling
		tempPolicy := struct {
			Version   string `json:"Version"`
			Statement []struct {
				Effect   string   `json:"Effect"`
				Action   []string `json:"Action"`
				Resource string   `json:"Resource"`
			} `json:"Statement"`
		}{
			Version: internalPolicy.Version,
			Statement: []struct {
				Effect   string   `json:"Effect"`
				Action   []string `json:"Action"`
				Resource string   `json:"Resource"`
			}{
				{
					Effect:   firstStmt.Effect,
					Action:   firstStmt.Action,
					Resource: resource,
				},
			},
		}

		jsonBytes, err := json.Marshal(tempPolicy)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}).(pulumi.StringOutput), nil
}

func CreateEcrDeployPolicy(ctx *pulumi.Context, awsAccountId string, repositoryName string) (*iam.GetPolicyDocumentResult, error) {
	policyDocument, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Version: pulumi.StringRef("2012-10-17"),
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Actions: []string{
					"ecr:GetDownloadUrlForLayer",
					"ecr:BatchGetImage",
					"ecr:BatchCheckLayerAvailability",
					"ecr:PutImage",
					"ecr:InitiateLayerUpload",
					"ecr:UploadLayerPart",
					"ecr:CompleteLayerUpload",
				},
				Resources: []string{
					fmt.Sprintf("arn:aws:ecr:us-east-1:%s:repository/%s", awsAccountId, repositoryName),
				},
			},
			{
				Effect: pulumi.StringRef("Allow"),
				Actions: []string{
					"ecr:GetAuthorizationToken",
				},
				Resources: []string{
					"*",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return policyDocument, nil

}

func CreateGithubAssumeRolePolicy(ctx *pulumi.Context, awsAccountId string, githubOrganization string, githubRepository string) (*iam.GetPolicyDocumentResult, error) {
	policyDocument, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Version: pulumi.StringRef("2012-10-17"),
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Federated",
						Identifiers: []string{fmt.Sprintf("arn:aws:iam::%s:oidc-provider/token.actions.githubusercontent.com", awsAccountId)},
					},
				},
				Actions: []string{"sts:AssumeRoleWithWebIdentity"},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringEquals",
						Variable: "token.actions.githubusercontent.com:aud",
						Values:   []string{"sts.amazonaws.com"},
					},
					{
						Test:     "StringLike",
						Variable: "token.actions.githubusercontent.com:sub",
						Values:   []string{fmt.Sprintf("repo:%s/%s:*", githubOrganization, githubRepository)},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return policyDocument, nil
}
