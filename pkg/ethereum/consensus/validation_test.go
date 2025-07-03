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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing image pull policy",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing p2p port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing beacon api port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing metrics port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "missing execution client endpoint",
			args: ConsensusClientArgs{
				Name:            "test",
				Namespace:       "default",
				StorageSize:     "100Gi",
				StorageClass:    "standard",
				Image:           "test-image",
				ImagePullPolicy: "IfNotPresent",
				JWTSecret:       "test-secret",
				P2PPort:         30303,
				BeaconAPIPort:   5052,
				MetricsPort:     9090,
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
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 0,
				BeaconAPIPort:           5052,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "invalid beacon api port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           0,
				MetricsPort:             9090,
				ExecutionClientEndpoint: "http://execution:8551",
			},
			wantErr: true,
		},
		{
			name: "invalid metrics port",
			args: ConsensusClientArgs{
				Name:                    "test",
				Namespace:               "default",
				StorageSize:             "100Gi",
				StorageClass:            "standard",
				Image:                   "test-image",
				ImagePullPolicy:         "IfNotPresent",
				JWTSecret:               "test-secret",
				P2PPort:                 30303,
				BeaconAPIPort:           5052,
				MetricsPort:             0,
				ExecutionClientEndpoint: "http://execution:8551",
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
