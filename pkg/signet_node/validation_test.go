package signet_node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignetNodeComponentArgsValidate(t *testing.T) {
	// Test with valid args (without ChainName, so conditional fields required)
	validArgs := SignetNodeComponentArgs{
		Name:      "test-node",
		Namespace: "default",
		Env: SignetNodeEnv{
			// Always required
			BlobExplorerUrl:    "http://explorer",
			SignetStaticPath:   "/static",
			SignetDatabasePath: "/db",
			// Optional with defaults
			RpcPort:   5959,
			WsRpcPort: 5960,
			// Optional (no defaults)
			RustLog: "info",
			// Conditional (required since ChainName is not set)
			RollupGenesisJsonPath: "/rollup-genesis.json",
			HostGenesisJsonPath:   "/host-genesis.json",
			StartTimestamp:        123456789,
			SlotOffset:            1,
			SlotDuration:          12,
		},
		ExecutionJwt:                "jwt-token",
		ExecutionPvcSize:            "150Gi",
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
		ExecutionPvcSize:            validArgs.ExecutionPvcSize,
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
		ExecutionPvcSize:            validArgs.ExecutionPvcSize,
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
		ExecutionPvcSize:            validArgs.ExecutionPvcSize,
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
		ExecutionPvcSize:            validArgs.ExecutionPvcSize,
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

	// Test with missing execution pvc size
	invalidArgs5 := SignetNodeComponentArgs{
		Name:                        "test-node",
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

	err = invalidArgs5.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution pvc size is required")
}

func TestSignetNodeEnvValidate(t *testing.T) {
	// Test with ChainName set - conditional fields not required
	t.Run("valid with ChainName set", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:    "http://explorer",
			SignetStaticPath:   "/static",
			SignetDatabasePath: "/db",
			ChainName:          "Pecorino",
		}
		err := env.Validate()
		assert.NoError(t, err)
	})

	// Test without ChainName - conditional fields required
	t.Run("valid without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:       "http://explorer",
			SignetStaticPath:      "/static",
			SignetDatabasePath:    "/db",
			RollupGenesisJsonPath: "/rollup-genesis.json",
			HostGenesisJsonPath:   "/host-genesis.json",
			StartTimestamp:        123456789,
			SlotOffset:            0,
			SlotDuration:          12,
		}
		err := env.Validate()
		assert.NoError(t, err)
	})

	// Test missing always-required fields
	t.Run("missing blob explorer url", func(t *testing.T) {
		env := SignetNodeEnv{
			SignetStaticPath:   "/static",
			SignetDatabasePath: "/db",
			ChainName:          "Pecorino",
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "blob explorer url is required")
	})

	t.Run("missing signet static path", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:    "http://explorer",
			SignetDatabasePath: "/db",
			ChainName:          "Pecorino",
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signet static path is required")
	})

	t.Run("missing signet database path", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:  "http://explorer",
			SignetStaticPath: "/static",
			ChainName:        "Pecorino",
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signet database path is required")
	})

	// Test missing conditional fields when ChainName is not set
	t.Run("missing rollup genesis json path without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:     "http://explorer",
			SignetStaticPath:    "/static",
			SignetDatabasePath:  "/db",
			HostGenesisJsonPath: "/host-genesis.json",
			StartTimestamp:      123456789,
			SlotOffset:          0,
			SlotDuration:        12,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rollup genesis json path is required when chain name is not set")
	})

	t.Run("missing host genesis json path without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:       "http://explorer",
			SignetStaticPath:      "/static",
			SignetDatabasePath:    "/db",
			RollupGenesisJsonPath: "/rollup-genesis.json",
			StartTimestamp:        123456789,
			SlotOffset:            0,
			SlotDuration:          12,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "host genesis json path is required when chain name is not set")
	})

	t.Run("missing start timestamp without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:       "http://explorer",
			SignetStaticPath:      "/static",
			SignetDatabasePath:    "/db",
			RollupGenesisJsonPath: "/rollup-genesis.json",
			HostGenesisJsonPath:   "/host-genesis.json",
			SlotOffset:            0,
			SlotDuration:          12,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start timestamp must be a positive integer when chain name is not set")
	})

	t.Run("invalid slot offset without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:       "http://explorer",
			SignetStaticPath:      "/static",
			SignetDatabasePath:    "/db",
			RollupGenesisJsonPath: "/rollup-genesis.json",
			HostGenesisJsonPath:   "/host-genesis.json",
			StartTimestamp:        123456789,
			SlotOffset:            -1,
			SlotDuration:          12,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "slot offset must be a non-negative integer when chain name is not set")
	})

	t.Run("missing slot duration without ChainName", func(t *testing.T) {
		env := SignetNodeEnv{
			BlobExplorerUrl:       "http://explorer",
			SignetStaticPath:      "/static",
			SignetDatabasePath:    "/db",
			RollupGenesisJsonPath: "/rollup-genesis.json",
			HostGenesisJsonPath:   "/host-genesis.json",
			StartTimestamp:        123456789,
			SlotOffset:            0,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "slot duration must be a positive integer when chain name is not set")
	})
}

func TestApplyDefaults(t *testing.T) {
	// Test with no mount paths specified
	args := SignetNodeComponentArgs{}
	args.ApplyDefaults()

	assert.Equal(t, DefaultSignetNodeDataMountPath, args.SignetNodeDataMountPath)
	assert.Equal(t, DefaultRollupDataMountPath, args.RollupDataMountPath)
	assert.Equal(t, DefaultExecutionJwtMountPath, args.ExecutionJwtMountPath)

	// Test with custom mount paths
	customArgs := SignetNodeComponentArgs{
		SignetNodeDataMountPath: "/custom/signet/data",
		RollupDataMountPath:     "/custom/rollup/data",
		ExecutionJwtMountPath:   "/custom/jwt",
	}
	customArgs.ApplyDefaults()

	assert.Equal(t, "/custom/signet/data", customArgs.SignetNodeDataMountPath)
	assert.Equal(t, "/custom/rollup/data", customArgs.RollupDataMountPath)
	assert.Equal(t, "/custom/jwt", customArgs.ExecutionJwtMountPath)

	// Test with partial custom mount paths
	partialArgs := SignetNodeComponentArgs{
		SignetNodeDataMountPath: "/custom/signet/data",
	}
	partialArgs.ApplyDefaults()

	assert.Equal(t, "/custom/signet/data", partialArgs.SignetNodeDataMountPath)
	assert.Equal(t, DefaultRollupDataMountPath, partialArgs.RollupDataMountPath)
	assert.Equal(t, DefaultExecutionJwtMountPath, partialArgs.ExecutionJwtMountPath)
}

func TestSignetNodeEnvApplyDefaults(t *testing.T) {
	// Test with no ports specified - should apply defaults
	t.Run("applies default ports", func(t *testing.T) {
		env := SignetNodeEnv{}
		env.ApplyDefaults()

		assert.Equal(t, DefaultSignetRpcPort, env.RpcPort)
		assert.Equal(t, DefaultSignetWsRpcPort, env.WsRpcPort)
	})

	// Test with custom ports - should preserve them
	t.Run("preserves custom ports", func(t *testing.T) {
		env := SignetNodeEnv{
			RpcPort:   8545,
			WsRpcPort: 8546,
		}
		env.ApplyDefaults()

		assert.Equal(t, 8545, env.RpcPort)
		assert.Equal(t, 8546, env.WsRpcPort)
	})

	// Test with partial custom ports
	t.Run("applies defaults only for unset ports", func(t *testing.T) {
		env := SignetNodeEnv{
			RpcPort: 9000,
		}
		env.ApplyDefaults()

		assert.Equal(t, 9000, env.RpcPort)
		assert.Equal(t, DefaultSignetWsRpcPort, env.WsRpcPort)
	})
}

func TestSignetNodeComponentArgsApplyDefaultsIncludesEnv(t *testing.T) {
	// Test that ApplyDefaults on SignetNodeComponentArgs also applies env defaults
	args := SignetNodeComponentArgs{}
	args.ApplyDefaults()

	assert.Equal(t, DefaultSignetRpcPort, args.Env.RpcPort)
	assert.Equal(t, DefaultSignetWsRpcPort, args.Env.WsRpcPort)
}
