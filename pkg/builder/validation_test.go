package builder

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestBuilderComponentArgsValidate(t *testing.T) {
	// Test validation with missing fields
	testCases := []struct {
		name    string
		args    BuilderComponentArgs
		wantErr string
	}{
		{
			name:    "missing namespace",
			args:    BuilderComponentArgs{},
			wantErr: "namespace is required",
		},
		{
			name: "missing name",
			args: BuilderComponentArgs{
				Namespace: "test-namespace",
			},
			wantErr: "name is required",
		},
		{
			name: "missing image",
			args: BuilderComponentArgs{
				Namespace: "test-namespace",
				Name:      "test-builder",
			},
			wantErr: "image is required",
		},
		{
			name: "missing app labels",
			args: BuilderComponentArgs{
				Namespace: "test-namespace",
				Name:      "test-builder",
				Image:     "test-image:latest",
			},
			wantErr: "app labels are required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.args.Validate()
			assert.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}

	// Test with valid args but empty BuilderEnv
	t.Run("valid args but empty BuilderEnv", func(t *testing.T) {
		args := BuilderComponentArgs{
			Namespace: "test-namespace",
			Name:      "test-builder",
			Image:     "test-image:latest",
			AppLabels: AppLabels{Labels: pulumi.StringMap{"app": pulumi.String("test")}},
		}
		err := args.Validate()
		assert.Error(t, err)
		assert.Equal(t, "auth token refresh interval is required", err.Error())
	})

	// Test with valid args and minimal valid BuilderEnv
	t.Run("valid args with minimal BuilderEnv", func(t *testing.T) {
		args := BuilderComponentArgs{
			Namespace: "test-namespace",
			Name:      "test-builder",
			Image:     "test-image:latest",
			AppLabels: AppLabels{Labels: pulumi.StringMap{"app": pulumi.String("test")}},
			BuilderEnv: BuilderEnv{
				AuthTokenRefreshInterval: "300",
				AwsAccountId:             "123456789012",
				AwsAccessKeyId:           "test-access-key",
				AwsRegion:                "us-west-2",
				AwsSecretAccessKey:       "test-secret-key",
				BlockConfirmationBuffer:  10,
				BlockQueryCutoff:         2000,
				BlockQueryStart:          1000,
				BuilderHelperAddress:     "0x123456",
				BuilderKey:               "test-key",
				BuilderPort:              8080,
				BuilderRewardsAddress:    "0x789abc",
				ChainOffset:              10,
				ConcurrentLimit:          100,
				HostChainId:              1,
				HostRpcUrl:               "http://host-rpc",
				OauthAudience:            "audience",
				OauthAuthenticateUrl:     "http://auth",
				OAuthClientId:            "client-id",
				OauthClientSecret:        "secret",
				OauthIssuer:              "issuer",
				OauthTokenUrl:            "http://token",
				OtelExporterOtlpEndpoint: "http://otel",
				QuinceyUrl:               "http://quincey",
				RollupBlockGasLimit:      30000000,
				RollupChainId:            2,
				RollupRpcUrl:             "http://rollup-rpc",
				RustLog:                  "info",
				SlotOffset:               10,
				StartTimestamp:           123456789,
				SubmitViaCallData:        "true",
				TargetSlotTime:           20,
				TxBroadcastUrls:          "http://broadcast",
				TxPoolCacheDuration:      60,
				TxPoolUrl:                "http://txpool",
				ZenithAddress:            "0xdef456",
			},
		}
		err := args.Validate()
		assert.NoError(t, err)
	})
}
