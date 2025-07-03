package execution

import (
	"testing"
)

func TestExecutionClientArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    ExecutionClientArgs
		wantErr bool
	}{
		{
			name: "missing name",
			args: ExecutionClientArgs{
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			args: ExecutionClientArgs{
				Name:            "test",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "missing storage size",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "missing storage class",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "missing image",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "missing image pull policy",
			args: ExecutionClientArgs{
				Name:          "test",
				Namespace:     "default",
				StorageSize:   "100Gi",
				StorageClass:  "standard",
				Image:         "test-image",
				JWTSecret:     "test-secret",
				P2PPort:       30303,
				RPCPort:       8545,
				WSPort:        8546,
				MetricsPort:   9090,
				AuthRPCPort:   8551,
				DiscoveryPort: 30303,
			},
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "invalid p2p port",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         0,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: true,
		},
		{
			name: "valid args",
			args: ExecutionClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "Always",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				RPCPort:         8545,
				WSPort:          8546,
				MetricsPort:     9090,
				AuthRPCPort:     8551,
				DiscoveryPort:   30303,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecutionClientArgs.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
