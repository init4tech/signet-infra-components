package pylon

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PylonComponentArgs struct {
	Namespace           string
	Name                string
	DbProjectName       string
	ExecutionJwt        pulumi.StringInput
	PylonBlobBucketName string
	Env                 *Env
}

type Env struct {
	PylonStartBlock            pulumi.StringInput
	PylonS3Url                 pulumi.StringInput
	PylonS3Region              pulumi.StringInput
	PylonSenderAddress         pulumi.StringInput
	PylonNetworkSlotDuration   pulumi.StringInput
	PylonNetworkSlotOffset     pulumi.StringInput
	PylonRequestsPerSecond     pulumi.StringInput
	PylonRustLog               pulumi.StringInput
	PylonPort                  pulumi.StringInput
	AwsAccessKeyId             pulumi.StringInput
	AwsSecretAccessKey         pulumi.StringInput
	AwsRegion                  pulumi.StringInput
	PostgresUser               pulumi.StringInput
	PostgresPassword           pulumi.StringInput
	PylonConsensusClientUrl    pulumi.StringInput
	PylonBlobscanBaseUrl       pulumi.StringInput
	PylonNetworkStartTimestamp pulumi.StringInput
}

type PylonComponent struct {
	pulumi.ResourceState
}
