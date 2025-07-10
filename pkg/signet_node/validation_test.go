package signet_node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignetNodeComponentArgsValidate(t *testing.T) {
	// Test with valid args
	validArgs := SignetNodeComponentArgs{
		Name:      "test-node",
		Namespace: "default",
		Env: SignetNodeEnv{
			HostZenithAddress:             "0x123",
			RuOrdersContractAddress:       "0x789",
			HostOrdersContractAddress:     "0x456",
			SignetChainId:                 123,
			BlobExplorerUrl:               "http://explorer",
			SignetStaticPath:              "/static",
			SignetDatabasePath:            "/db",
			RustLog:                       "info",
			IpcEndpoint:                   "/tmp/reth.ipc",
			RpcPort:                       8545,
			WsRpcPort:                     8546,
			TxForwardUrl:                  "http://forward",
			GenesisJsonPath:               "/genesis.json",
			HostZenithDeployHeight:        "1000",
			BaseFeeRecipient:              "0xabc",
			HostPassageContractAddress:    "0xdef",
			HostTransactorContractAddress: "0x321",
			RuPassageContractAddress:      "0x654",
			SignetClUrl:                   "http://cl",
			SignetPylonUrl:                "http://pylon",
			HostStartTimestamp:            123456789,
			HostSlotOffset:                1,
			HostSlotDuration:              12,
		},
		ExecutionJwt:                "jwt-token",
		LighthousePvcSize:           "100Gi",
		RollupPvcSize:               "50Gi",
		ExecutionClientImage:        "execution:latest",
		ConsensusClientImage:        "consensus:latest",
		ExecutionClientStartCommand: []string{"./start-execution"},
		ConsensusClientStartCommand: []string{"./start-consensus"},
		AppLabels:                   AppLabels{}, // Labels not needed for validation
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test with missing name
	invalidArgs1 := SignetNodeComponentArgs{
		Namespace:                   "default",
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
		Namespace:                   "default",
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
		Namespace:                   "default",
		ExecutionJwt:                validArgs.ExecutionJwt,
		LighthousePvcSize:           validArgs.LighthousePvcSize,
		RollupPvcSize:               validArgs.RollupPvcSize,
		ExecutionClientImage:        validArgs.ExecutionClientImage,
		ConsensusClientImage:        validArgs.ConsensusClientImage,
		ExecutionClientStartCommand: validArgs.ExecutionClientStartCommand,
		ConsensusClientStartCommand: validArgs.ConsensusClientStartCommand,
		Env:                         SignetNodeEnv{}, // Missing required fields
	}

	err = invalidArgs4.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signet node env")
}
