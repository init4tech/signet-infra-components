package pylon

import (
	"fmt"
)

// Validate validates the PylonComponentArgs struct
func (args *PylonComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
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
	if env.PylonStartBlock == "" {
		return fmt.Errorf("pylonStartBlock is required")
	}

	if env.PylonSenders == "" {
		return fmt.Errorf("pylonSenders is required")
	}

	if env.PylonS3Url == "" {
		return fmt.Errorf("pylonS3Url is required")
	}

	if env.PylonS3Region == "" {
		return fmt.Errorf("pylonS3Region is required")
	}

	if env.PylonClUrl == "" {
		return fmt.Errorf("pylonClUrl is required")
	}

	if env.PylonBlobscanBaseUrl == "" {
		return fmt.Errorf("pylonBlobscanBaseUrl is required")
	}

	if env.PylonNetworkStartTimestamp == "" {
		return fmt.Errorf("pylonNetworkStartTimestamp is required")
	}

	if env.PylonNetworkSlotDuration == "" {
		return fmt.Errorf("pylonNetworkSlotDuration is required")
	}

	if env.PylonNetworkSlotOffset == "" {
		return fmt.Errorf("pylonNetworkSlotOffset is required")
	}

	if env.PylonRequestsPerSecond == "" {
		return fmt.Errorf("pylonRequestsPerSecond is required")
	}

	if env.PylonPort == "" {
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
