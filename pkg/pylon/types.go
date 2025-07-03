package pylon

import (
	"github.com/init4tech/signet-infra-components/pkg/ethereum"
	"github.com/init4tech/signet-infra-components/pkg/utils"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PylonComponentArgs struct {
	Namespace           string
	Name                string
	DbProjectName       string
	ExecutionJwt        pulumi.StringInput
	PylonImage          pulumi.StringInput
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
	PylonDbUrl                 pulumi.StringInput
	PylonConsensusClientUrl    pulumi.StringInput
	PylonBlobscanBaseUrl       pulumi.StringInput
	PylonNetworkStartTimestamp pulumi.StringInput
}

// GetEnvMap implements the utils.EnvProvider interface
// It creates a map of environment variables from the Env struct
func (e Env) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

type PylonComponent struct {
	pulumi.ResourceState
	EthereumNode      *ethereum.EthereumNodeComponent
	PylonEnvConfigMap *v1.ConfigMap
}
