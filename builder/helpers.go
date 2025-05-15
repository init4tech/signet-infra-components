package builder

import (
	"reflect"
	"strings"
	"unicode"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createKMSPolicy creates a KMS policy for the builder service.
func createKMSPolicy(key pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"kms:Sign",
					"kms:GetPublicKey"
				],
				"Resource": %s
			}
		]
	}`, key)
}

// createEnvVars creates environment variables by automatically mapping
// struct field names to environment variable names.
func createEnvVars(env BuilderEnv) corev1.EnvVarArray {
	result := corev1.EnvVarArray{}

	// Special case for BuilderPort as it needs string conversion
	result = append(result, &corev1.EnvVarArgs{
		Name:  pulumi.String("BUILDER_PORT"),
		Value: pulumi.Sprintf("%d", env.BuilderPort),
	})

	// Process all string inputs from the struct's tags
	envVarMap := getEnvironmentVarsFromStruct(env)
	for name, value := range envVarMap {
		if name != "BUILDER_PORT" { // Skip the one we already handled
			result = append(result, &corev1.EnvVarArgs{
				Name:  pulumi.String(name),
				Value: value.(pulumi.StringInput),
			})
		}
	}

	return result
}

// getEnvironmentVarsFromStruct uses reflection to extract environment variables from struct tags
func getEnvironmentVarsFromStruct(env BuilderEnv) map[string]pulumi.Input {
	result := make(map[string]pulumi.Input)

	t := reflect.TypeOf(env)
	v := reflect.ValueOf(env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the field value
		fieldValue := v.Field(i).Interface()

		// Skip nil values and BuilderPort (handled specially)
		if fieldValue == nil || field.Name == "BuilderPort" {
			continue
		}

		// Convert camelCase to SNAKE_CASE for env var name
		envName := camelToSnake(field.Name)

		// Add to map
		result[envName] = fieldValue.(pulumi.Input)
	}

	return result
}

// camelToSnake converts a camelCase string to SNAKE_CASE
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}
