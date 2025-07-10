package ethereum

import (
	"testing"

	"github.com/init4tech/signet-infra-components/pkg/ethereum/consensus"
	"github.com/init4tech/signet-infra-components/pkg/ethereum/execution"
	"github.com/stretchr/testify/assert"
)

func TestEthereumNodeArgsValidate(t *testing.T) {
	// Test with valid args
	validArgs := EthereumNodeArgs{
		Name:      "test-node",
		Namespace: "default",
		ExecutionClient: &execution.ExecutionClientArgs{
			Name:            "test-execution",
			Namespace:       "default",
			StorageSize:     "100Gi",
			StorageClass:    "standard",
			Image:           "test-execution-image",
			ImagePullPolicy: "Always",
			JWTSecret:       "test-jwt-secret",
			P2PPort:         30303,
			RPCPort:         8545,
			WSPort:          8546,
			MetricsPort:     9090,
			AuthRPCPort:     8551,
			DiscoveryPort:   30303,
		},
		ConsensusClient: &consensus.ConsensusClientArgs{
			Name:                    "test-consensus",
			Namespace:               "default",
			StorageSize:             "100Gi",
			StorageClass:            "standard",
			Image:                   "test-consensus-image",
			ImagePullPolicy:         "Always",
			JWTSecret:               "test-jwt-secret",
			P2PPort:                 30303,
			BeaconAPIPort:           5052,
			MetricsPort:             9090,
			ExecutionClientEndpoint: "http://execution:8551",
		},
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test with missing name
	invalidArgs1 := EthereumNodeArgs{
		Namespace:       "default",
		ExecutionClient: validArgs.ExecutionClient,
		ConsensusClient: validArgs.ConsensusClient,
	}

	err = invalidArgs1.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test with missing namespace
	invalidArgs2 := EthereumNodeArgs{
		Name:            "test-node",
		ExecutionClient: validArgs.ExecutionClient,
		ConsensusClient: validArgs.ConsensusClient,
	}

	err = invalidArgs2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")

	// Test with missing execution client
	invalidArgs3 := EthereumNodeArgs{
		Name:            "test-node",
		Namespace:       "default",
		ConsensusClient: validArgs.ConsensusClient,
	}

	err = invalidArgs3.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "executionClient is required")

	// Test with missing consensus client
	invalidArgs4 := EthereumNodeArgs{
		Name:            "test-node",
		Namespace:       "default",
		ExecutionClient: validArgs.ExecutionClient,
	}

	err = invalidArgs4.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "consensusClient is required")

	// Test with invalid execution client
	invalidExecutionClient := &execution.ExecutionClientArgs{
		// Missing required fields
	}
	invalidArgs5 := EthereumNodeArgs{
		Name:            "test-node",
		Namespace:       "default",
		ExecutionClient: invalidExecutionClient,
		ConsensusClient: validArgs.ConsensusClient,
	}

	err = invalidArgs5.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution client validation failed")

	// Test with invalid consensus client
	invalidConsensusClient := &consensus.ConsensusClientArgs{
		// Missing required fields
	}
	invalidArgs6 := EthereumNodeArgs{
		Name:            "test-node",
		Namespace:       "default",
		ExecutionClient: validArgs.ExecutionClient,
		ConsensusClient: invalidConsensusClient,
	}

	err = invalidArgs6.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "consensus client validation failed")
}
