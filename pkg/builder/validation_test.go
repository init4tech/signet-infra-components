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
		assert.Equal(t, "builder port is required", err.Error())
	})

	// Test with valid args and minimal valid BuilderEnv
	t.Run("valid args with minimal BuilderEnv", func(t *testing.T) {
		args := BuilderComponentArgs{
			Namespace: "test-namespace",
			Name:      "test-builder",
			Image:     "test-image:latest",
			AppLabels: AppLabels{Labels: pulumi.StringMap{"app": pulumi.String("test")}},
			BuilderEnv: BuilderEnv{
				BuilderPort:              pulumi.String("8080"),
				BuilderKey:               pulumi.String("test-key"),
				HostRpcUrl:               pulumi.String("http://host-rpc"),
				RollupRpcUrl:             pulumi.String("http://rollup-rpc"),
				ZenithAddress:            pulumi.String("0x123456"),
				QuinceyUrl:               pulumi.String("http://quincey"),
				OtelExporterOtlpEndpoint: pulumi.String("http://otel"),
				OauthAudience:            pulumi.String("audience"),
				OauthAuthenticateUrl:     pulumi.String("http://auth"),
				OAuthClientId:            pulumi.String("client-id"),
				OauthClientSecret:        pulumi.String("secret"),
				OauthIssuer:              pulumi.String("issuer"),
				OauthTokenUrl:            pulumi.String("http://token"),
				RustLog:                  pulumi.String("info"),
				SlotOffset:               pulumi.String("10"),
				StartTimestamp:           pulumi.String("123456789"),
				SubmitViaCallData:        pulumi.String("true"),
				TargetSlotTime:           pulumi.String("20"),
				TxBroadcastUrls:          pulumi.String("http://broadcast"),
				TxPoolCacheDuration:      pulumi.String("60"),
				TxPoolUrl:                pulumi.String("http://txpool"),
			},
		}
		err := args.Validate()
		assert.NoError(t, err)
	})
}

func TestBuilderEnvValidate(t *testing.T) {
	// Test just one case for BuilderEnv validation
	t.Run("missing builder port", func(t *testing.T) {
		env := BuilderEnv{}
		err := env.Validate()
		assert.Error(t, err)
		assert.Equal(t, "builder port is required", err.Error())
	})

	// Test with just BuilderPort set
	t.Run("missing builder key", func(t *testing.T) {
		env := BuilderEnv{
			BuilderPort: pulumi.String("8080"),
		}
		err := env.Validate()
		assert.Error(t, err)
		assert.Equal(t, "builder key is required", err.Error())
	})
}
