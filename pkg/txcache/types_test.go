package txcache

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestTxCacheEnvInternal_GetEnvMap(t *testing.T) {
	env := TxCacheEnvInternal{
		HttpPort:                  pulumi.String("8080"),
		AwsAccessKeyId:            pulumi.String("test-key"),
		AwsSecretAccessKey:        pulumi.String("test-secret"),
		AwsRegion:                 pulumi.String("us-east-1"),
		RustLog:                   pulumi.String("info"),
		BlockQueryStart:           pulumi.String("0"),
		BlockQueryCutoff:          pulumi.String("100"),
		SlotOffset:                pulumi.String("10"),
		ExpirationTimestampOffset: pulumi.String("3600"),
		NetworkName:               pulumi.String("mainnet"),
		Builders:                  pulumi.String("builder1,builder2"),
		SlotDuration:              pulumi.String("12"),
		StartTimestamp:            pulumi.String("1234567890"),
		OtelExporterOtlpProtocol:  pulumi.String("grpc"),
		OtelExporterOtlpEndpoint:  pulumi.String("localhost:4317"),
	}

	envMap := env.GetEnvMap()

	// Verify all required fields are mapped correctly
	assert.Equal(t, pulumi.String("8080"), envMap["HTTP_PORT"])
	assert.Equal(t, pulumi.String("test-key"), envMap["AWS_ACCESS_KEY_ID"])
	assert.Equal(t, pulumi.String("test-secret"), envMap["AWS_SECRET_ACCESS_KEY"])
	assert.Equal(t, pulumi.String("us-east-1"), envMap["AWS_REGION"])
	assert.Equal(t, pulumi.String("info"), envMap["RUST_LOG"])
	assert.Equal(t, pulumi.String("0"), envMap["BLOCK_QUERY_START"])
	assert.Equal(t, pulumi.String("100"), envMap["BLOCK_QUERY_CUTOFF"])
	assert.Equal(t, pulumi.String("10"), envMap["SLOT_OFFSET"])
	assert.Equal(t, pulumi.String("3600"), envMap["EXPIRATION_TIMESTAMP_OFFSET"])
	assert.Equal(t, pulumi.String("mainnet"), envMap["NETWORK_NAME"])
	assert.Equal(t, pulumi.String("builder1,builder2"), envMap["BUILDERS"])
	assert.Equal(t, pulumi.String("12"), envMap["SLOT_DURATION"])
	assert.Equal(t, pulumi.String("1234567890"), envMap["START_TIMESTAMP"])
	assert.Equal(t, pulumi.String("grpc"), envMap["OTEL_EXPORTER_OTLP_PROTOCOL"])
	assert.Equal(t, pulumi.String("localhost:4317"), envMap["OTEL_EXPORTER_OTLP_ENDPOINT"])
}

func TestTxCacheEnvInternal_GetEnvMap_FieldNamesConversion(t *testing.T) {
	// Test that camelCase field names are properly converted to UPPER_SNAKE_CASE
	env := TxCacheEnvInternal{
		HttpPort:                  pulumi.String("8080"),
		AwsAccessKeyId:            pulumi.String("key"),
		AwsSecretAccessKey:        pulumi.String("secret"),
		AwsRegion:                 pulumi.String("region"),
		RustLog:                   pulumi.String("debug"),
		BlockQueryStart:           pulumi.String("0"),
		BlockQueryCutoff:          pulumi.String("10"),
		SlotOffset:                pulumi.String("5"),
		ExpirationTimestampOffset: pulumi.String("60"),
		NetworkName:               pulumi.String("test"),
		Builders:                  pulumi.String("b1"),
		SlotDuration:              pulumi.String("6"),
		StartTimestamp:            pulumi.String("123"),
	}

	envMap := env.GetEnvMap()

	// Ensure all expected env var names exist
	expectedKeys := []string{
		"HTTP_PORT",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
		"RUST_LOG",
		"BLOCK_QUERY_START",
		"BLOCK_QUERY_CUTOFF",
		"SLOT_OFFSET",
		"EXPIRATION_TIMESTAMP_OFFSET",
		"NETWORK_NAME",
		"BUILDERS",
		"SLOT_DURATION",
		"START_TIMESTAMP",
	}

	for _, key := range expectedKeys {
		_, exists := envMap[key]
		assert.True(t, exists, "Expected environment variable %s to exist", key)
	}
}