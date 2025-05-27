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
				AuthTokenRefreshInterval: pulumi.String("300"),
				AwsAccountId:             pulumi.String("123456789012"),
				AwsAccessKeyId:           pulumi.String("test-access-key"),
				AwsRegion:                pulumi.String("us-west-2"),
				AwsSecretAccessKey:       pulumi.String("test-secret-key"),
				BlockConfirmationBuffer:  pulumi.String("10"),
				BlockQueryCutoff:         pulumi.String("2000"),
				BlockQueryStart:          pulumi.String("1000"),
				BuilderHelperAddress:     pulumi.String("0x123456"),
				BuilderKey:               pulumi.String("test-key"),
				BuilderPort:              pulumi.String("8080"),
				BuilderRewardsAddress:    pulumi.String("0x789abc"),
				ChainOffset:              pulumi.String("10"),
				ConcurrentLimit:          pulumi.String("100"),
				HostChainId:              pulumi.String("1"),
				HostRpcUrl:               pulumi.String("http://host-rpc"),
				OauthAudience:            pulumi.String("audience"),
				OauthAuthenticateUrl:     pulumi.String("http://auth"),
				OAuthClientId:            pulumi.String("client-id"),
				OauthClientSecret:        pulumi.String("secret"),
				OauthIssuer:              pulumi.String("issuer"),
				OauthTokenUrl:            pulumi.String("http://token"),
				OtelExporterOtlpEndpoint: pulumi.String("http://otel"),
				QuinceyUrl:               pulumi.String("http://quincey"),
				RollupBlockGasLimit:      pulumi.String("30000000"),
				RollupChainId:            pulumi.String("2"),
				RollupRpcUrl:             pulumi.String("http://rollup-rpc"),
				RustLog:                  pulumi.String("info"),
				SlotOffset:               pulumi.String("10"),
				StartTimestamp:           pulumi.String("123456789"),
				SubmitViaCallData:        pulumi.String("true"),
				TargetSlotTime:           pulumi.String("20"),
				TxBroadcastUrls:          pulumi.String("http://broadcast"),
				TxPoolCacheDuration:      pulumi.String("60"),
				TxPoolUrl:                pulumi.String("http://txpool"),
				ZenithAddress:            pulumi.String("0xdef456"),
			},
		}
		err := args.Validate()
		assert.NoError(t, err)
	})
}
