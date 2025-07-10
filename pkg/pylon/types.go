package pylon

import (
	"github.com/init4tech/signet-infra-components/pkg/ethereum"
	"github.com/init4tech/signet-infra-components/pkg/utils"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Public-facing structs with base Go types
type PylonComponentArgs struct {
	Namespace           string
	Name                string
	DbProjectName       string
	ExecutionJwt        string
	PylonImage          string
	PylonBlobBucketName string
	Env                 PylonEnv
}

// Internal structs with Pulumi types for use within the component
type pylonComponentArgsInternal struct {
	Namespace           pulumi.StringInput
	Name                pulumi.StringInput
	DbProjectName       pulumi.StringInput
	ExecutionJwt        pulumi.StringInput
	PylonImage          pulumi.StringInput
	PylonBlobBucketName pulumi.StringInput
	Env                 pylonEnvInternal
}

// Public-facing environment struct with base Go types
type PylonEnv struct {
	PylonStartBlock            string `pulumi:"pylonStartBlock" validate:"required"`
	PylonS3Url                 string `pulumi:"pylonS3Url" validate:"required"`
	PylonS3Region              string `pulumi:"pylonS3Region" validate:"required"`
	PylonSenderAddress         string `pulumi:"pylonSenderAddress" validate:"required"`
	PylonNetworkSlotDuration   string `pulumi:"pylonNetworkSlotDuration" validate:"required"`
	PylonNetworkSlotOffset     string `pulumi:"pylonNetworkSlotOffset" validate:"required"`
	PylonRequestsPerSecond     string `pulumi:"pylonRequestsPerSecond" validate:"required"`
	PylonRustLog               string `pulumi:"pylonRustLog"`
	PylonPort                  string `pulumi:"pylonPort" validate:"required"`
	AwsAccessKeyId             string `pulumi:"awsAccessKeyId" validate:"required"`
	AwsSecretAccessKey         string `pulumi:"awsSecretAccessKey" validate:"required"`
	AwsRegion                  string `pulumi:"awsRegion" validate:"required"`
	PylonDbUrl                 string `pulumi:"pylonDbUrl" validate:"required"`
	PylonConsensusClientUrl    string `pulumi:"pylonConsensusClientUrl" validate:"required"`
	PylonBlobscanBaseUrl       string `pulumi:"pylonBlobscanBaseUrl" validate:"required"`
	PylonNetworkStartTimestamp string `pulumi:"pylonNetworkStartTimestamp" validate:"required"`
}

// Internal environment struct with Pulumi types
type pylonEnvInternal struct {
	PylonStartBlock            pulumi.StringInput `pulumi:"pylonStartBlock" validate:"required"`
	PylonS3Url                 pulumi.StringInput `pulumi:"pylonS3Url" validate:"required"`
	PylonS3Region              pulumi.StringInput `pulumi:"pylonS3Region" validate:"required"`
	PylonSenderAddress         pulumi.StringInput `pulumi:"pylonSenderAddress" validate:"required"`
	PylonNetworkSlotDuration   pulumi.StringInput `pulumi:"pylonNetworkSlotDuration" validate:"required"`
	PylonNetworkSlotOffset     pulumi.StringInput `pulumi:"pylonNetworkSlotOffset" validate:"required"`
	PylonRequestsPerSecond     pulumi.StringInput `pulumi:"pylonRequestsPerSecond" validate:"required"`
	PylonRustLog               pulumi.StringInput `pulumi:"pylonRustLog"`
	PylonPort                  pulumi.StringInput `pulumi:"pylonPort" validate:"required"`
	AwsAccessKeyId             pulumi.StringInput `pulumi:"awsAccessKeyId" validate:"required"`
	AwsSecretAccessKey         pulumi.StringInput `pulumi:"awsSecretAccessKey" validate:"required"`
	AwsRegion                  pulumi.StringInput `pulumi:"awsRegion" validate:"required"`
	PylonDbUrl                 pulumi.StringInput `pulumi:"pylonDbUrl" validate:"required"`
	PylonConsensusClientUrl    pulumi.StringInput `pulumi:"pylonConsensusClientUrl" validate:"required"`
	PylonBlobscanBaseUrl       pulumi.StringInput `pulumi:"pylonBlobscanBaseUrl" validate:"required"`
	PylonNetworkStartTimestamp pulumi.StringInput `pulumi:"pylonNetworkStartTimestamp" validate:"required"`
}

// Conversion function to convert public args to internal args
func (args PylonComponentArgs) toInternal() pylonComponentArgsInternal {
	return pylonComponentArgsInternal{
		Namespace:           pulumi.String(args.Namespace),
		Name:                pulumi.String(args.Name),
		DbProjectName:       pulumi.String(args.DbProjectName),
		ExecutionJwt:        pulumi.String(args.ExecutionJwt),
		PylonImage:          pulumi.String(args.PylonImage),
		PylonBlobBucketName: pulumi.String(args.PylonBlobBucketName),
		Env:                 args.Env.toInternal(),
	}
}

// Conversion function to convert public env to internal env
func (e PylonEnv) toInternal() pylonEnvInternal {
	return pylonEnvInternal{
		PylonStartBlock:            pulumi.String(e.PylonStartBlock),
		PylonS3Url:                 pulumi.String(e.PylonS3Url),
		PylonS3Region:              pulumi.String(e.PylonS3Region),
		PylonSenderAddress:         pulumi.String(e.PylonSenderAddress),
		PylonNetworkSlotDuration:   pulumi.String(e.PylonNetworkSlotDuration),
		PylonNetworkSlotOffset:     pulumi.String(e.PylonNetworkSlotOffset),
		PylonRequestsPerSecond:     pulumi.String(e.PylonRequestsPerSecond),
		PylonRustLog:               pulumi.String(e.PylonRustLog),
		PylonPort:                  pulumi.String(e.PylonPort),
		AwsAccessKeyId:             pulumi.String(e.AwsAccessKeyId),
		AwsSecretAccessKey:         pulumi.String(e.AwsSecretAccessKey),
		AwsRegion:                  pulumi.String(e.AwsRegion),
		PylonDbUrl:                 pulumi.String(e.PylonDbUrl),
		PylonConsensusClientUrl:    pulumi.String(e.PylonConsensusClientUrl),
		PylonBlobscanBaseUrl:       pulumi.String(e.PylonBlobscanBaseUrl),
		PylonNetworkStartTimestamp: pulumi.String(e.PylonNetworkStartTimestamp),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal env
func (e pylonEnvInternal) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

type PylonComponent struct {
	pulumi.ResourceState
	EthereumNode      *ethereum.EthereumNodeComponent
	PylonEnvConfigMap *v1.ConfigMap
}
