package execution

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
				Namespace:     pulumi.String("default"),
				StorageSize:   pulumi.String("100Gi"),
				StorageClass:  pulumi.String("standard"),
				Image:         pulumi.String("test-image"),
				JWTSecret:     pulumi.String("test-secret"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			args: ExecutionClientArgs{
				Name:          pulumi.String("test"),
				StorageSize:   pulumi.String("100Gi"),
				StorageClass:  pulumi.String("standard"),
				Image:         pulumi.String("test-image"),
				JWTSecret:     pulumi.String("test-secret"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
		},
		{
			name: "missing storage size",
			args: ExecutionClientArgs{
				Name:          pulumi.String("test"),
				Namespace:     pulumi.String("default"),
				StorageClass:  pulumi.String("standard"),
				Image:         pulumi.String("test-image"),
				JWTSecret:     pulumi.String("test-secret"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
		},
		{
			name: "missing storage class",
			args: ExecutionClientArgs{
				Name:          pulumi.String("test"),
				Namespace:     pulumi.String("default"),
				StorageSize:   pulumi.String("100Gi"),
				Image:         pulumi.String("test-image"),
				JWTSecret:     pulumi.String("test-secret"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
		},
		{
			name: "missing image",
			args: ExecutionClientArgs{
				Name:          pulumi.String("test"),
				Namespace:     pulumi.String("default"),
				StorageSize:   pulumi.String("100Gi"),
				StorageClass:  pulumi.String("standard"),
				JWTSecret:     pulumi.String("test-secret"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			args: ExecutionClientArgs{
				Name:          pulumi.String("test"),
				Namespace:     pulumi.String("default"),
				StorageSize:   pulumi.String("100Gi"),
				StorageClass:  pulumi.String("standard"),
				Image:         pulumi.String("test-image"),
				P2PPort:       pulumi.Int(30303),
				RPCPort:       pulumi.Int(8545),
				WSPort:        pulumi.Int(8546),
				MetricsPort:   pulumi.Int(9090),
				AuthRPCPort:   pulumi.Int(8551),
				DiscoveryPort: pulumi.Int(30303),
			},
			wantErr: true,
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
