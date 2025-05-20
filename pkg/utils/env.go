package utils

import (
	"reflect"
	"strings"
	"unicode"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// EnvProvider is an interface for structs that provide environment variables
type EnvProvider interface {
	// GetEnvMap converts the struct's fields to a map of environment variables
	GetEnvMap() pulumi.StringMap
}

// CreateConfigMap creates a Kubernetes ConfigMap from an environment variables struct
// It automatically converts field names to environment variable format (UPPER_SNAKE_CASE)
func CreateConfigMap(
	ctx *pulumi.Context,
	name string,
	namespace pulumi.StringInput,
	labels pulumi.StringMap,
	env interface{},
) (*corev1.ConfigMap, error) {
	// Create metadata for ConfigMap
	metadata := &metav1.ObjectMetaArgs{
		Name:      pulumi.String(name),
		Namespace: namespace,
		Labels:    labels,
	}

	// Get environment variables as a map
	data := pulumi.StringMap{}

	// If the env object implements EnvProvider, use its GetEnvMap method
	if provider, ok := env.(EnvProvider); ok {
		data = provider.GetEnvMap()
	} else {
		// Otherwise use reflection to extract fields
		data = CreateEnvMap(env)
	}

	// Create and return ConfigMap
	return corev1.NewConfigMap(ctx, name, &corev1.ConfigMapArgs{
		Metadata: metadata,
		Data:     data,
	})
}

// CreateEnvMap converts a struct to a map of environment variables
// Field names are converted from camelCase to UPPER_SNAKE_CASE
func CreateEnvMap(env interface{}) pulumi.StringMap {
	result := pulumi.StringMap{}
	t := reflect.TypeOf(env)
	v := reflect.ValueOf(env)

	// If not a struct, return empty map
	if t.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i).Interface()

		// Skip nil values
		if fieldValue == nil {
			continue
		}

		// Convert camelCase to SNAKE_CASE for env var name
		envName := CamelToSnake(field.Name)

		// Convert value to string using appropriate method
		if stringValue := getStringValue(fieldValue); stringValue != nil {
			result[envName] = stringValue
		}
	}

	return result
}

// getStringValue converts any pulumi Input to a StringInput
func getStringValue(value interface{}) pulumi.StringInput {
	switch v := value.(type) {
	case pulumi.StringInput:
		return v
	case pulumi.IntInput:
		return pulumi.Sprintf("%d", v)
	case pulumi.BoolInput:
		return pulumi.Sprintf("%t", v)
	case pulumi.Float64Input:
		return pulumi.Sprintf("%f", v)
	default:
		// Skip values we don't know how to convert
		return nil
	}
}

// CamelToSnake converts a camelCase string to SNAKE_CASE
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
