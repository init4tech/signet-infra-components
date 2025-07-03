package quincey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuinceyComponentArgsValidate(t *testing.T) {
	// Test validation with missing fields
	testCases := []struct {
		name    string
		args    QuinceyComponentArgs
		wantErr string
	}{
		{
			name:    "missing namespace",
			args:    QuinceyComponentArgs{},
			wantErr: "namespace is required",
		},
		{
			name: "missing image",
			args: QuinceyComponentArgs{
				Namespace: "test-namespace",
			},
			wantErr: "image is required",
		},
		{
			name: "missing virtual service hosts",
			args: QuinceyComponentArgs{
				Namespace: "test-namespace",
				Image:     "test-image:latest",
				Port:      8080,
				Env: QuinceyEnv{
					QuinceyPort:        8080,
					QuinceyKeyId:       "test-key-id",
					AwsAccessKeyId:     "test-access-key",
					AwsSecretAccessKey: "test-secret-key",
					AwsDefaultRegion:   "us-west-2",
					BlockQueryStart:    1000,
					BlockQueryCutoff:   2000,
					ChainOffset:        10,
					HostRpcUrl:         "http://host-rpc",
					OauthIssuer:        "https://issuer",
					OauthJwksUri:       "https://jwks",
					QuinceyBuilders:    "builder1,builder2",
				},
			},
			wantErr: "virtual service hosts is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.args.Validate()
			assert.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}

	// Test with valid args but empty QuinceyEnv
	t.Run("valid args but empty QuinceyEnv", func(t *testing.T) {
		args := QuinceyComponentArgs{
			Namespace:           "test-namespace",
			Image:               "test-image:latest",
			Port:                8080,
			VirtualServiceHosts: []string{"example.com"},
		}
		err := args.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "env is invalid")
	})

	// Test with valid args and complete QuinceyEnv
	t.Run("valid args with complete QuinceyEnv", func(t *testing.T) {
		args := QuinceyComponentArgs{
			Namespace:           "test-namespace",
			Image:               "test-image:latest",
			Port:                8080,
			VirtualServiceHosts: []string{"example.com"},
			Env: QuinceyEnv{
				QuinceyPort:        8080,
				QuinceyKeyId:       "test-key-id",
				AwsAccessKeyId:     "test-access-key",
				AwsSecretAccessKey: "test-secret-key",
				AwsDefaultRegion:   "us-west-2",
				BlockQueryStart:    1000,
				BlockQueryCutoff:   2000,
				ChainOffset:        10,
				HostRpcUrl:         "http://host-rpc",
				OauthIssuer:        "https://issuer",
				OauthJwksUri:       "https://jwks",
				QuinceyBuilders:    "builder1,builder2",
				// Optional fields
				OtelExporterOtlpEndpoint: "http://otel",
				OtelExporterOtlpProtocol: "grpc",
				RustLog:                  "info",
			},
		}
		err := args.Validate()
		assert.NoError(t, err)
	})
}
