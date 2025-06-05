package consensus

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: false,
		},
		{
			name: "missing name",
			args: ConsensusClientArgs{
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing storage size",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing storage class",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing image",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing image pull policy",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing p2p port",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				BeaconAPIPort:           pulumi.Int(5052),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing beacon api port",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				MetricsPort:             pulumi.Int(9090),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing metrics port",
			args: ConsensusClientArgs{
				Name:                    pulumi.String("test"),
				Namespace:               pulumi.String("default"),
				StorageSize:             pulumi.String("100Gi"),
				StorageClass:            pulumi.String("standard"),
				Image:                   pulumi.String("test-image"),
				ImagePullPolicy:         pulumi.String("IfNotPresent"),
				JWTSecret:               pulumi.String("test-secret"),
				P2PPort:                 pulumi.Int(30303),
				BeaconAPIPort:           pulumi.Int(5052),
				ExecutionClientEndpoint: pulumi.String("http://execution:8551"),
			},
			wantErr: true,
		},
		{
			name: "missing execution client endpoint",
			args: ConsensusClientArgs{
				Name:            pulumi.String("test"),
				Namespace:       pulumi.String("default"),
				StorageSize:     pulumi.String("100Gi"),
				StorageClass:    pulumi.String("standard"),
				Image:           pulumi.String("test-image"),
				ImagePullPolicy: pulumi.String("IfNotPresent"),
				JWTSecret:       pulumi.String("test-secret"),
				P2PPort:         pulumi.Int(30303),
				BeaconAPIPort:   pulumi.Int(5052),
				MetricsPort:     pulumi.Int(9090),
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
