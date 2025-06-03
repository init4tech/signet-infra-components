package consensus

import (
	"testing"
)

func TestConsensusClientArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    ConsensusClientArgs
		wantErr bool
	}{
		{
			name: "valid args",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			args: ConsensusClientArgs{
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			args: ConsensusClientArgs{
				Name:                    "test",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
		{
			name: "missing storage size",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageClass:            "standard",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
		{
			name: "missing storage class",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
		{
			name: "missing image",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 9000,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
		{
			name: "missing execution client endpoint",
			args: ConsensusClientArgs{
				Name:          "test",
				Namespace:     "default",
				StorageSize:   "100Gi",
				StorageClass:  "standard",
				Image:         "test-image",
				P2PPort:       9000,
				MetricsPort:   9090,
				BeaconAPIPort: 5052,
			},
			wantErr: true,
		},
		{
			name: "invalid p2p port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ExecutionClientEndpoint: "http://localhost:8551",
				P2PPort:                 0,
				MetricsPort:             9090,
				BeaconAPIPort:           5052,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsensusClientArgs.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
