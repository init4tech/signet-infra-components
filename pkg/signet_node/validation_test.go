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
			HostZenithAddress:             pulumi.String("0x123"),
			RuOrdersContractAddress:       pulumi.String("0x789"),
			HostOrdersContractAddress:     pulumi.String("0x456"),
			SignetChainId:                 pulumi.String("123"),
			BlobExplorerUrl:               pulumi.String("http://explorer"),
			SignetStaticPath:              pulumi.String("/static"),
			SignetDatabasePath:            pulumi.String("/db"),
			RustLog:                       pulumi.String("info"),
			IpcEndpoint:                   pulumi.String("/tmp/reth.ipc"),
			RpcPort:                       pulumi.String("8545"),
			WsRpcPort:                     pulumi.String("8546"),
			TxForwardUrl:                  pulumi.String("http://forward"),
			GenesisJsonPath:               pulumi.String("/genesis.json"),
			HostZenithDeployHeight:        pulumi.String("1000"),
			BaseFeeRecipient:              pulumi.String("0xabc"),
			HostPassageContractAddress:    pulumi.String("0xdef"),
			HostTransactorContractAddress: pulumi.String("0x321"),
			RuPassageContractAddress:      pulumi.String("0x654"),
			SignetClUrl:                   pulumi.String("http://cl"),
			SignetPylonUrl:                pulumi.String("http://pylon"),
			HostStartTimestamp:            pulumi.String("123456789"),
			HostSlotOffset:                pulumi.String("0"),
			HostSlotDuration:              pulumi.String("12"),
		},
		ExecutionJwt:                pulumi.String("jwt-token"),
		LighthousePvcSize:           pulumi.String("100Gi"),
		RollupPvcSize:               pulumi.String("50Gi"),
		ExecutionClientImage:        pulumi.String("execution:latest"),
		ConsensusClientImage:        pulumi.String("consensus:latest"),
		ExecutionClientStartCommand: pulumi.StringArray{pulumi.String("./start-execution")},
		ConsensusClientStartCommand: pulumi.StringArray{pulumi.String("./start-consensus")},
		AppLabels:                   AppLabels{Labels: pulumi.StringMap{"app": pulumi.String("test")}},
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test with missing name
	invalidArgs1 := SignetNodeComponentArgs{
		Namespace:                   pulumi.String("default"),
		Env:                         validArgs.Env,
		ExecutionJwt:                validArgs.ExecutionJwt,
		LighthousePvcSize:           validArgs.LighthousePvcSize,
		RollupPvcSize:               validArgs.RollupPvcSize,
		ExecutionClientImage:        validArgs.ExecutionClientImage,
		ConsensusClientImage:        validArgs.ConsensusClientImage,
		ExecutionClientStartCommand: validArgs.ExecutionClientStartCommand,
		ConsensusClientStartCommand: validArgs.ConsensusClientStartCommand,
	}

	err = invalidArgs1.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test with missing namespace
	invalidArgs2 := SignetNodeComponentArgs{
		Name:                        "test-node",
		Env:                         validArgs.Env,
		ExecutionJwt:                validArgs.ExecutionJwt,
		LighthousePvcSize:           validArgs.LighthousePvcSize,
		RollupPvcSize:               validArgs.RollupPvcSize,
		ExecutionClientImage:        validArgs.ExecutionClientImage,
		ConsensusClientImage:        validArgs.ConsensusClientImage,
		ExecutionClientStartCommand: validArgs.ExecutionClientStartCommand,
		ConsensusClientStartCommand: validArgs.ConsensusClientStartCommand,
	}

	err = invalidArgs2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")

	// Test with missing execution jwt
	invalidArgs3 := SignetNodeComponentArgs{
		Name:                        "test-node",
		Namespace:                   pulumi.String("default"),
		Env:                         validArgs.Env,
		LighthousePvcSize:           validArgs.LighthousePvcSize,
		RollupPvcSize:               validArgs.RollupPvcSize,
		ExecutionClientImage:        validArgs.ExecutionClientImage,
		ConsensusClientImage:        validArgs.ConsensusClientImage,
		ExecutionClientStartCommand: validArgs.ExecutionClientStartCommand,
		ConsensusClientStartCommand: validArgs.ConsensusClientStartCommand,
	}

	err = invalidArgs3.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution jwt is required")

	// Test with invalid env
	invalidArgs4 := SignetNodeComponentArgs{
		Name:                        "test-node",
		Namespace:                   pulumi.String("default"),
		ExecutionJwt:                validArgs.ExecutionJwt,
		LighthousePvcSize:           validArgs.LighthousePvcSize,
		RollupPvcSize:               validArgs.RollupPvcSize,
		ExecutionClientImage:        validArgs.ExecutionClientImage,
		ConsensusClientImage:        validArgs.ConsensusClientImage,
		ExecutionClientStartCommand: validArgs.ExecutionClientStartCommand,
		ConsensusClientStartCommand: validArgs.ConsensusClientStartCommand,
		Env:                         SignetNodeEnv{
			// Missing required fields
		},
	}

	err = invalidArgs4.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signet node env")
}
