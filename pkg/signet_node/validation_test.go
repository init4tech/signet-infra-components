package signet_node

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestSignetNodeComponentArgsValidate(t *testing.T) {
	// Test with valid args
	validArgs := SignetNodeComponentArgs{
		Name:      "test-node",
		Namespace: pulumi.String("default"),
		Env: SignetNodeEnv{
			HostZenithAddress:         pulumi.String("0x123"),
			HostOrdersContractAddress: pulumi.String("0x456"),
			SignetChainId:             pulumi.String("123"),
			BlobExplorerUrl:           pulumi.String("http://explorer"),
			SignetStaticPath:          pulumi.String("/static"),
			SignetDatabasePath:        pulumi.String("/db"),
			RpcPort:                   pulumi.String("8545"),
			WsRpcPort:                 pulumi.String("8546"),
			TxForwardUrl:              pulumi.String("http://forward"),
			GenesisJsonPath:           pulumi.String("/genesis.json"),
			HostStartTimestamp:        pulumi.String("123456789"),
			HostSlotOffset:            pulumi.String("0"),
			HostSlotDuration:          pulumi.String("12"),
		},
		ExecutionJwt: pulumi.String("jwt-token"),
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test with missing name
	invalidArgs1 := SignetNodeComponentArgs{
		Namespace: pulumi.String("default"),
		Env:       validArgs.Env,
	}

	err = invalidArgs1.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test with missing namespace
	invalidArgs2 := SignetNodeComponentArgs{
		Name: "test-node",
		Env:  validArgs.Env,
	}

	err = invalidArgs2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")

	// Test with invalid env
	invalidArgs3 := SignetNodeComponentArgs{
		Name:      "test-node",
		Namespace: pulumi.String("default"),
		Env:       SignetNodeEnv{
			// Missing required fields
		},
	}

	err = invalidArgs3.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signet node env")
}

func TestSignetNodeEnvValidate(t *testing.T) {
	// Test with valid env
	validEnv := SignetNodeEnv{
		HostZenithAddress:         pulumi.String("0x123"),
		HostOrdersContractAddress: pulumi.String("0x456"),
		SignetChainId:             pulumi.String("123"),
		BlobExplorerUrl:           pulumi.String("http://explorer"),
		SignetStaticPath:          pulumi.String("/static"),
		SignetDatabasePath:        pulumi.String("/db"),
		RpcPort:                   pulumi.String("8545"),
		WsRpcPort:                 pulumi.String("8546"),
		TxForwardUrl:              pulumi.String("http://forward"),
		GenesisJsonPath:           pulumi.String("/genesis.json"),
		HostStartTimestamp:        pulumi.String("123456789"),
		HostSlotOffset:            pulumi.String("0"),
		HostSlotDuration:          pulumi.String("12"),
	}

	err := validEnv.Validate()
	assert.NoError(t, err)

	// Test with missing fields
	tests := []struct {
		name          string
		modifyEnv     func(*SignetNodeEnv)
		expectedError string
	}{
		{
			name: "missing HostZenithAddress",
			modifyEnv: func(env *SignetNodeEnv) {
				env.HostZenithAddress = nil
			},
			expectedError: "host zenith address is required",
		},
		{
			name: "missing HostOrdersContractAddress",
			modifyEnv: func(env *SignetNodeEnv) {
				env.HostOrdersContractAddress = nil
			},
			expectedError: "host orders contract address is required",
		},
		{
			name: "missing SignetChainId",
			modifyEnv: func(env *SignetNodeEnv) {
				env.SignetChainId = nil
			},
			expectedError: "signet chain id is required",
		},
		{
			name: "missing RpcPort",
			modifyEnv: func(env *SignetNodeEnv) {
				env.RpcPort = nil
			},
			expectedError: "rpc port is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a copy of the valid env
			envCopy := validEnv
			// Modify it for the test case
			test.modifyEnv(&envCopy)

			err := envCopy.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		})
	}
}
