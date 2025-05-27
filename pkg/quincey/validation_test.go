package quincey

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
				Namespace: pulumi.String("test-namespace"),
			},
			wantErr: "image is required",
		},
		{
			name: "missing virtual service hosts",
			args: QuinceyComponentArgs{
				Namespace: pulumi.String("test-namespace"),
				Image:     pulumi.String("test-image:latest"),
				Env: QuinceyEnv{
					QuinceyPort:        pulumi.String("8080"),
					QuinceyKeyId:       pulumi.String("test-key-id"),
					AwsAccessKeyId:     pulumi.String("test-access-key"),
					AwsSecretAccessKey: pulumi.String("test-secret-key"),
					AwsDefaultRegion:   pulumi.String("us-west-2"),
					BlockQueryStart:    pulumi.String("1000"),
					BlockQueryCutoff:   pulumi.String("2000"),
					ChainOffset:        pulumi.String("10"),
					HostRpcUrl:         pulumi.String("http://host-rpc"),
					OauthIssuer:        pulumi.String("https://issuer"),
					OauthJwksUri:       pulumi.String("https://jwks"),
					QuinceyBuilders:    pulumi.String("builder1,builder2"),
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
			Namespace:           pulumi.String("test-namespace"),
			Image:               pulumi.String("test-image:latest"),
			VirtualServiceHosts: pulumi.StringArray{pulumi.String("example.com")},
		}
		err := args.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "env is invalid")
	})

	// Test with valid args and complete QuinceyEnv
	t.Run("valid args with complete QuinceyEnv", func(t *testing.T) {
		args := QuinceyComponentArgs{
			Namespace:           pulumi.String("test-namespace"),
			Image:               pulumi.String("test-image:latest"),
			VirtualServiceHosts: pulumi.StringArray{pulumi.String("example.com")},
			Env: QuinceyEnv{
				QuinceyPort:        pulumi.String("8080"),
				QuinceyKeyId:       pulumi.String("test-key-id"),
				AwsAccessKeyId:     pulumi.String("test-access-key"),
				AwsSecretAccessKey: pulumi.String("test-secret-key"),
				AwsDefaultRegion:   pulumi.String("us-west-2"),
				BlockQueryStart:    pulumi.String("1000"),
				BlockQueryCutoff:   pulumi.String("2000"),
				ChainOffset:        pulumi.String("10"),
				HostRpcUrl:         pulumi.String("http://host-rpc"),
				OauthIssuer:        pulumi.String("https://issuer"),
				OauthJwksUri:       pulumi.String("https://jwks"),
				QuinceyBuilders:    pulumi.String("builder1,builder2"),
				// Optional fields
				OtelExporterOtlpEndpoint: pulumi.String("http://otel"),
				OtelExporterOtlpProtocol: pulumi.String("grpc"),
				RustLog:                  pulumi.String("info"),
			},
		}
		err := args.Validate()
		assert.NoError(t, err)
	})
}
