package pylon

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Validate validates the PylonComponentArgs struct
func (args *PylonComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if args.DbProjectName == "" {
		return fmt.Errorf("dbProjectName is required")
	}

	if args.ExecutionJwt == "" {
		return fmt.Errorf("executionJwt is required")
	}

	if args.PylonImage == "" {
		return fmt.Errorf("pylonImage is required")
	}

	if args.PylonBlobBucketName == "" {
		return fmt.Errorf("pylonBlobBucketName is required")
	}

	if err := validateEnv(args.Env); err != nil {
		return err
	}

	return nil
}

func ValidatePylon(ctx *pulumi.Context, args *PylonComponentArgs) error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if args.DbProjectName == "" {
		return fmt.Errorf("dbProjectName is required")
	}

	if args.ExecutionJwt == "" {
		return fmt.Errorf("executionJwt is required")
	}

	if args.PylonImage == "" {
		return fmt.Errorf("pylonImage is required")
	}

	if args.PylonBlobBucketName == "" {
		return fmt.Errorf("pylonBlobBucketName is required")
	}

	if err := validateEnv(args.Env); err != nil {
		return err
	}

	return nil
}

func validateEnv(env PylonEnv) error {
	if env.PylonStartBlock == 0 {
		return fmt.Errorf("pylonStartBlock is required")
	}

	if env.PylonSenderAddress == "" {
		return fmt.Errorf("pylonSenderAddress is required")
	}

	if env.PylonS3Url == "" {
		return fmt.Errorf("pylonS3Url is required")
	}

	if env.PylonS3Region == "" {
		return fmt.Errorf("pylonS3Region is required")
	}

	if env.PylonConsensusClientUrl == "" {
		return fmt.Errorf("pylonConsensusClientUrl is required")
	}

	if env.PylonBlobscanBaseUrl == "" {
		return fmt.Errorf("pylonBlobscanBaseUrl is required")
	}

	if env.PylonNetworkStartTimestamp == 0 {
		return fmt.Errorf("pylonNetworkStartTimestamp is required")
	}

	if env.PylonNetworkSlotDuration == 0 {
		return fmt.Errorf("pylonNetworkSlotDuration is required")
	}

	if env.PylonNetworkSlotOffset < 0 {
		return fmt.Errorf("pylonNetworkSlotOffset must be a non-negative integer")
	}

	if env.PylonRequestsPerSecond == 0 {
		return fmt.Errorf("pylonRequestsPerSecond is required")
	}

	if env.PylonPort == 0 {
		return fmt.Errorf("pylonPort is required")
	}

	if env.AwsAccessKeyId == "" {
		return fmt.Errorf("awsAccessKeyId is required")
	}

	if env.AwsSecretAccessKey == "" {
		return fmt.Errorf("awsSecretAccessKey is required")
	}

	if env.AwsRegion == "" {
		return fmt.Errorf("awsRegion is required")
	}

	if env.PylonDbUrl == "" {
		return fmt.Errorf("pylonDbUrl is required")
	}

	return nil
}
