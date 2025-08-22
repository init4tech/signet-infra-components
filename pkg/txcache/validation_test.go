package txcache

import (
	"testing"
)

func TestTxCacheComponentArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    TxCacheComponentArgs
		wantErr bool
	}{
		{
			name: "valid args",
			args: TxCacheComponentArgs{
				Namespace:    "test-namespace",
				Name:         "test-name",
				Image:        "test-image",
				Port:         8080,
				OauthIssuer:  "test-issuer",
				OauthJwksUri: "test-jwks-uri",
				Env: TxCacheEnv{
					HttpPort:                  "8080",
					AwsAccessKeyId:            "test-key",
					AwsSecretAccessKey:        "test-secret",
					AwsRegion:                 "us-west-2",
					RustLog:                   "info",
					BlockQueryStart:           "1000",
					BlockQueryCutoff:          "2000",
					SlotOffset:                "0",
					ExpirationTimestampOffset: "3600",
					NetworkName:               "testnet",
					Builders:                  "builder1,builder2",
					SlotDuration:              "12",
					StartTimestamp:            "1640995200",
				},
			},
			wantErr: false,
		},
		{
			name: "missing namespace",
			args: TxCacheComponentArgs{
				Name: "test-name",
				Env: TxCacheEnv{
					HttpPort:                  "8080",
					AwsAccessKeyId:            "test-key",
					AwsSecretAccessKey:        "test-secret",
					AwsRegion:                 "us-west-2",
					RustLog:                   "info",
					BlockQueryStart:           "1000",
					BlockQueryCutoff:          "2000",
					SlotOffset:                "0",
					ExpirationTimestampOffset: "3600",
					NetworkName:               "testnet",
					Builders:                  "builder1,builder2",
					SlotDuration:              "12",
					StartTimestamp:            "1640995200",
				},
			},
			wantErr: true,
		},
		{
			name: "missing name",
			args: TxCacheComponentArgs{
				Namespace: "test-namespace",
				Env: TxCacheEnv{
					HttpPort:                  "8080",
					AwsAccessKeyId:            "test-key",
					AwsSecretAccessKey:        "test-secret",
					AwsRegion:                 "us-west-2",
					RustLog:                   "info",
					BlockQueryStart:           "1000",
					BlockQueryCutoff:          "2000",
					SlotOffset:                "0",
					ExpirationTimestampOffset: "3600",
					NetworkName:               "testnet",
					Builders:                  "builder1,builder2",
					SlotDuration:              "12",
					StartTimestamp:            "1640995200",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TxCacheComponentArgs.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     TxCacheEnv
		wantErr bool
	}{
		{
			name: "valid env",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: false,
		},
		{
			name: "valid env with optional fields",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
				OtelExporterOtlpProtocol:  "grpc",
				OtelExporterOtlpEndpoint:  "localhost:4317",
			},
			wantErr: false,
		},
		{
			name: "missing httpPort",
			env: TxCacheEnv{
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing awsAccessKeyId",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing awsSecretAccessKey",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing awsRegion",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing rustLog",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing blockQueryStart",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing blockQueryCutoff",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing slotOffset",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing expirationTimestampOffset",
			env: TxCacheEnv{
				HttpPort:           "8080",
				AwsAccessKeyId:     "test-key",
				AwsSecretAccessKey: "test-secret",
				AwsRegion:          "us-west-2",
				RustLog:            "info",
				BlockQueryStart:    "1000",
				BlockQueryCutoff:   "2000",
				SlotOffset:         "0",
				NetworkName:        "testnet",
				Builders:           "builder1,builder2",
				SlotDuration:       "12",
				StartTimestamp:     "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing networkName",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing builders",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				SlotDuration:              "12",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing slotDuration",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				StartTimestamp:            "1640995200",
			},
			wantErr: true,
		},
		{
			name: "missing startTimestamp",
			env: TxCacheEnv{
				HttpPort:                  "8080",
				AwsAccessKeyId:            "test-key",
				AwsSecretAccessKey:        "test-secret",
				AwsRegion:                 "us-west-2",
				RustLog:                   "info",
				BlockQueryStart:           "1000",
				BlockQueryCutoff:          "2000",
				SlotOffset:                "0",
				ExpirationTimestampOffset: "3600",
				NetworkName:               "testnet",
				Builders:                  "builder1,builder2",
				SlotDuration:              "12",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnv(tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
