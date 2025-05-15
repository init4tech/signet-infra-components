package builder

import (
	"reflect"
	"strings"
	"unicode"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateKMSPolicy creates a KMS policy for the builder service.
// Exported for testing.
func CreateKMSPolicy(key pulumi.StringInput) pulumi.StringOutput {
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

// createKMSPolicy is the internal version kept for backward compatibility
func createKMSPolicy(key pulumi.StringInput) pulumi.StringOutput {
	return CreateKMSPolicy(key)
}

// CreateEnvVars creates environment variables by automatically mapping
// struct field names to environment variable names.
// Exported for testing.
func CreateEnvVars(env BuilderEnv) corev1.EnvVarArray {
	result := corev1.EnvVarArray{}

	// Special case for BuilderPort as it needs string conversion
	result = append(result, &corev1.EnvVarArgs{
		Name:  pulumi.String("BUILDER_PORT"),
		Value: pulumi.Sprintf("%d", env.BuilderPort),
	})

	// Process all string inputs from the struct's tags
	envVarMap := GetEnvironmentVarsFromStruct(env)
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

// createEnvVars is kept for backward compatibility
func createEnvVars(env BuilderEnv) corev1.EnvVarArray {
	return CreateEnvVars(env)
}

// GetEnvironmentVarsFromStruct uses reflection to extract environment variables from struct tags
// Exported for testing.
func GetEnvironmentVarsFromStruct(env BuilderEnv) map[string]pulumi.Input {
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
		envName := CamelToSnake(field.Name)

		// Add to map
		result[envName] = fieldValue.(pulumi.Input)
	}

	return result
}

// getEnvironmentVarsFromStruct is kept for backward compatibility
func getEnvironmentVarsFromStruct(env BuilderEnv) map[string]pulumi.Input {
	return GetEnvironmentVarsFromStruct(env)
}

// CamelToSnake converts a camelCase string to SNAKE_CASE
// Exported for testing purposes
func CamelToSnake(s string) string {
	var result strings.Builder

	// Handle consecutive uppercase characters (acronyms)
	for i, r := range s {
		if unicode.IsUpper(r) {
			// Add underscore if not the first character and either:
			// 1. Previous character is lowercase, or
			// 2. Not the last character and next character is lowercase (end of acronym)
			needsUnderscore := i > 0 && (unicode.IsLower(rune(s[i-1])) ||
				(i < len(s)-1 && unicode.IsLower(rune(s[i+1])) && i > 1 && unicode.IsUpper(rune(s[i-1]))))

			if needsUnderscore {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}
