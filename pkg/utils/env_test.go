package utils

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type TestEnv struct {
	StringField    pulumi.StringInput `pulumi:"stringField"`
	IntField       pulumi.IntInput    `pulumi:"intField"`
	BoolField      pulumi.BoolInput   `pulumi:"boolField"`
	CamelCaseField pulumi.StringInput `pulumi:"camelCaseField"`
	HTTPSProxy     pulumi.StringInput `pulumi:"httpsProxy"`
}

// GetEnvMap implements the EnvProvider interface for TestEnv
func (e TestEnv) GetEnvMap() pulumi.StringMap {
	// Custom implementation can provide special handling
	result := pulumi.StringMap{
		"CUSTOM_FIELD": pulumi.String("custom-value"),
	}

	// Add the basic fields
	if e.StringField != nil {
		result["STRING_FIELD"] = e.StringField
	}
	if e.IntField != nil {
		result["INT_FIELD"] = pulumi.Sprintf("%d", e.IntField)
	}
	if e.BoolField != nil {
		result["BOOL_FIELD"] = pulumi.Sprintf("%t", e.BoolField)
	}
	if e.CamelCaseField != nil {
		result["CAMEL_CASE_FIELD"] = e.CamelCaseField
	}
	if e.HTTPSProxy != nil {
		result["HTTPS_PROXY"] = e.HTTPSProxy
	}

	return result
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simpleField", "SIMPLE_FIELD"},
		{"camelCaseField", "CAMEL_CASE_FIELD"},
		{"HTTPSProxy", "HTTPS_PROXY"},
		{"API", "API"},
		{"APIVersion", "API_VERSION"},
		{"SimpleHTTPServer", "SIMPLE_HTTP_SERVER"},
		{"PylonS3Region", "PYLON_S3_REGION"},
		{"EC2Instance", "EC2_INSTANCE"},
		{"Part1Part2", "PART1_PART2"},
		{"Version123Beta", "VERSION123_BETA"},
	}

	for _, test := range tests {
		result := CamelToSnake(test.input)
		if result != test.expected {
			t.Errorf("CamelToSnake(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
		valid    bool
	}{
		{pulumi.String("test"), "test", true},
		{pulumi.Int(42), "42", true},
		{pulumi.Bool(true), "true", true},
		{pulumi.Float64(3.14), "3.140000", true},
		{nil, "", false},
	}

	for _, test := range tests {
		result := getStringValue(test.input)
		if test.valid {
			if result == nil {
				t.Errorf("getStringValue(%v) returned nil, expected StringInput", test.input)
			}
		} else {
			if result != nil {
				t.Errorf("getStringValue(%v) returned non-nil, expected nil", test.input)
			}
		}
	}
}

func TestCreateEnvMap(t *testing.T) {
	env := TestEnv{
		StringField:    pulumi.String("string-value"),
		IntField:       pulumi.Int(42),
		BoolField:      pulumi.Bool(true),
		CamelCaseField: pulumi.String("camel-value"),
		HTTPSProxy:     pulumi.String("https://proxy.example.com"),
	}

	// Test automatic creation through reflection
	result := CreateEnvMap(env)

	// Should have 5 fields
	if len(result) != 5 {
		t.Errorf("Expected 5 environment variables, got %d", len(result))
	}

	// Check values
	stringValue, ok := result["STRING_FIELD"]
	if !ok || stringValue == nil {
		t.Errorf("STRING_FIELD not found or nil")
	}
}

func TestEnvProviderInterface(t *testing.T) {
	env := TestEnv{
		StringField:    pulumi.String("string-value"),
		IntField:       pulumi.Int(42),
		BoolField:      pulumi.Bool(true),
		CamelCaseField: pulumi.String("camel-value"),
		HTTPSProxy:     pulumi.String("https://proxy.example.com"),
	}

	// Test the interface implementation
	var provider EnvProvider = env
	result := provider.GetEnvMap()

	// Should have 6 fields (5 from struct + 1 custom)
	if len(result) != 6 {
		t.Errorf("Expected 6 environment variables via interface, got %d", len(result))
	}

	// Check for the custom field
	customValue, ok := result["CUSTOM_FIELD"]
	if !ok || customValue == nil {
		t.Errorf("CUSTOM_FIELD not found or nil")
	}
}
