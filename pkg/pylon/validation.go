package pylon

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ValidatePylon(ctx *pulumi.Context, args *PylonComponentArgs) error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if err := validateEnv(args.Env); err != nil {
		return err
	}

	return nil
}

func validateEnv(env *Env) error {
	if env.PylonStartBlock == nil {
		return fmt.Errorf("pylonStartBlock is required")
	}

	if env.PylonSenderAddress == nil {
		return fmt.Errorf("pylonSenderAddress is required")
	}

	if env.PylonS3Url == nil {
		return fmt.Errorf("pylonS3Url is required")
	}

	if env.PylonS3Region == nil {
		return fmt.Errorf("pylonS3Region is required")
	}

	if env.PylonConsensusClientUrl == nil {
		return fmt.Errorf("pylonConsensusClientUrl is required")
	}

	if env.PylonBlobscanBaseUrl == nil {
		return fmt.Errorf("pylonBlobscanBaseUrl is required")
	}

	if env.PylonNetworkStartTimestamp == nil {
		return fmt.Errorf("pylonNetworkStartTimestamp is required")
	}

	if env.PylonNetworkSlotDuration == nil {
		return fmt.Errorf("pylonNetworkSlotDuration is required")
	}

	if env.PylonNetworkSlotOffset == nil {
		return fmt.Errorf("pylonNetworkSlotOffset is required")
	}

	if env.PylonRequestsPerSecond == nil {
		return fmt.Errorf("pylonRequestsPerSecond is required")
	}

	if env.PylonRustLog == nil {
		return fmt.Errorf("pylonRustLog is required")
	}

	if env.PylonPort == nil {
		return fmt.Errorf("pylonPort is required")
	}

	if env.AwsAccessKeyId == nil {
		return fmt.Errorf("awsAccessKeyId is required")
	}

	if env.AwsSecretAccessKey == nil {
		return fmt.Errorf("awsSecretAccessKey is required")
	}

	if env.AwsRegion == nil {
		return fmt.Errorf("awsRegion is required")
	}

	if env.PostgresUser == nil {
		return fmt.Errorf("postgresUser is required")
	}

	if env.PostgresPassword == nil {
		return fmt.Errorf("postgresPassword is required")
	}

	return nil
}
