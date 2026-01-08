package signet_node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignetNodeComponentArgsValidate(t *testing.T) {
	validArgs := SignetNodeComponentArgs{
		Name:      "test-node",
		Namespace: "default",
		Env: SignetNodeEnv{
			ChainName:          "test-chain",
			IpcEndpoint:        "/tmp/reth.ipc",
			RpcPort:            8545,
			SignetChainId:      12345,
			SignetClUrl:        "https://cl.example.com",
			SignetDatabasePath: "/db/signet",
			SignetPylonUrl:     "https://pylon.example.com",
			SignetStaticPath:   "/static",
			TxForwardUrl:       "https://tx-forward.example.com",
			WsRpcPort:          8546,
			RustLog:            "info",
		},
		ExecutionJwt:                "jwt-secret-token",
		ExecutionPvcSize:            "150Gi",
		LighthousePvcSize:           "100Gi",
		RollupPvcSize:               "50Gi",
		ExecutionClientImage:        "execution:latest",
		ConsensusClientImage:        "consensus:latest",
		ExecutionClientStartCommand: []string{"./start-execution"},
		ConsensusClientStartCommand: []string{"./start-consensus"},
		AppLabels:                   AppLabels{},
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test missing name
	invalidName := validArgs
	invalidName.Name = ""
	err = invalidName.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test missing namespace
	invalidNamespace := validArgs
	invalidNamespace.Namespace = ""
	err = invalidNamespace.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")

	// Test missing execution jwt
	invalidJwt := validArgs
	invalidJwt.ExecutionJwt = ""
	err = invalidJwt.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution jwt is required")

	// Test missing execution pvc size
	invalidPvc := validArgs
	invalidPvc.ExecutionPvcSize = ""
	err = invalidPvc.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution pvc size is required")

	// Test invalid env (missing required field)
	invalidEnv := validArgs
	invalidEnv.Env.IpcEndpoint = ""
	err = invalidEnv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signet node env")
	assert.Contains(t, err.Error(), "ipcEndpoint is required")
}

func TestSignetNodeEnvValidate(t *testing.T) {
	t.Run("valid with all required fields", func(t *testing.T) {
		env := SignetNodeEnv{
			ChainName:          "Pecorino",
			IpcEndpoint:        "/tmp/reth.ipc",
			RpcPort:            8545,
			SignetChainId:      999,
			SignetClUrl:        "https://cl.example.com",
			SignetDatabasePath: "/data/db",
			SignetPylonUrl:     "https://pylon.example.com",
			SignetStaticPath:   "/static/files",
			TxForwardUrl:       "https://tx.example.com",
			WsRpcPort:          8546,
			RustLog:            "debug",
		}
		err := env.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid with RustLog nil", func(t *testing.T) {
		env := SignetNodeEnv{
			ChainName:          "Pecorino",
			IpcEndpoint:        "/tmp/reth.ipc",
			RpcPort:            8545,
			SignetChainId:      999,
			SignetClUrl:        "https://cl.example.com",
			SignetDatabasePath: "/data/db",
			SignetPylonUrl:     "https://pylon.example.com",
			SignetStaticPath:   "/static/files",
			TxForwardUrl:       "https://tx.example.com",
			WsRpcPort:          8546,
			RustLog:            "info",
		}
		err := env.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing ChainName", func(t *testing.T) {
		env := SignetNodeEnv{IpcEndpoint: "/tmp/reth.ipc"}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chainName is required")
	})

	t.Run("missing IpcEndpoint", func(t *testing.T) {
		env := SignetNodeEnv{ChainName: "Test"}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ipcEndpoint is required")
	})

	t.Run("invalid RpcPort (non-positive)", func(t *testing.T) {
		env := SignetNodeEnv{
			ChainName:          "Test",
			IpcEndpoint:        "/tmp/reth.ipc",
			RpcPort:            0,
			SignetChainId:      1,
			SignetClUrl:        "https://cl",
			SignetDatabasePath: "/db",
			SignetPylonUrl:     "https://pylon",
			SignetStaticPath:   "/static",
			TxForwardUrl:       "https://tx",
			WsRpcPort:          8546,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rpcPort must be a positive integer")
	})

	t.Run("invalid SignetChainId (non-positive)", func(t *testing.T) {
		env := SignetNodeEnv{
			ChainName:          "Test",
			IpcEndpoint:        "/tmp/reth.ipc",
			RpcPort:            8545,
			SignetChainId:      0,
			SignetClUrl:        "https://cl",
			SignetDatabasePath: "/db",
			SignetPylonUrl:     "https://pylon",
			SignetStaticPath:   "/static",
			TxForwardUrl:       "https://tx",
			WsRpcPort:          8546,
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signetChainId must be a positive integer")
	})
}

func TestApplyDefaults(t *testing.T) {
	t.Run("applies all defaults when empty", func(t *testing.T) {
		args := SignetNodeComponentArgs{}
		args.ApplyDefaults()

		assert.Equal(t, DefaultSignetNodeDataMountPath, args.SignetNodeDataMountPath)
		assert.Equal(t, DefaultRollupDataMountPath, args.RollupDataMountPath)
		assert.Equal(t, DefaultExecutionJwtMountPath, args.ExecutionJwtMountPath)
	})

	t.Run("preserves custom values", func(t *testing.T) {
		args := SignetNodeComponentArgs{
			SignetNodeDataMountPath: "/custom/signet",
			RollupDataMountPath:     "/custom/rollup",
			ExecutionJwtMountPath:   "/custom/jwt",
		}
		args.ApplyDefaults()

		assert.Equal(t, "/custom/signet", args.SignetNodeDataMountPath)
		assert.Equal(t, "/custom/rollup", args.RollupDataMountPath)
		assert.Equal(t, "/custom/jwt", args.ExecutionJwtMountPath)
	})

	t.Run("applies defaults only to unset fields", func(t *testing.T) {
		args := SignetNodeComponentArgs{
			SignetNodeDataMountPath: "/custom/signet",
		}
		args.ApplyDefaults()

		assert.Equal(t, "/custom/signet", args.SignetNodeDataMountPath)
		assert.Equal(t, DefaultRollupDataMountPath, args.RollupDataMountPath)
		assert.Equal(t, DefaultExecutionJwtMountPath, args.ExecutionJwtMountPath)
	})
}

// Helper to create *string
func stringPtr(s string) *string {
	return &s
}
