package aws

// IAM Policy version and effect constants
const (
	// IAMPolicyVersion is the standard IAM policy language version
	IAMPolicyVersion = "2012-10-17"

	// Policy effects
	EffectAllow = "Allow"
	EffectDeny  = "Deny"
)

// AWS Services
const (
	// EKS pod identity service
	EKSPodsService = "pods.eks.amazonaws.com"
	// EC2 service
	EC2Service = "ec2.amazonaws.com"
)

// STS Actions
const (
	STSAssumeRoleAction = "sts:AssumeRole"
	STSTagSessionAction = "sts:TagSession"
)

// KMS Actions
const (
	KMSSignAction         = "kms:Sign"
	KMSGetPublicKeyAction = "kms:GetPublicKey"
)

// IAM statement identifiers
const (
	EKSAssumeRoleStatementSid = "AllowEksAuthToAssumeRoleForPodIdentity"
)

// Resource name suffixes
const (
	RoleSuffix                 = "-role"
	PolicySuffix               = "-policy"
	RolePolicyAttachmentSuffix = "-role-policy-attachment"
)
