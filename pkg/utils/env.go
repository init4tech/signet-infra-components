package utils

import (
	"reflect"
	"strconv"
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
func CreateConfigMap[T EnvProvider](
	ctx *pulumi.Context,
	name string,
	namespace pulumi.StringInput,
	labels pulumi.StringMap,
	env T,
) (*corev1.ConfigMap, error) {
	// Create metadata for ConfigMap
	metadata := &metav1.ObjectMetaArgs{
		Name:      pulumi.String(name),
		Namespace: namespace,
		Labels:    labels,
	}

	// Get environment variables as a map using the EnvProvider interface
	data := env.GetEnvMap()

	// Create and return ConfigMap
	return corev1.NewConfigMap(ctx, name, &corev1.ConfigMapArgs{
		Metadata: metadata,
		Data:     data,
	})
}

// CreateEnvMap converts a struct to a map of environment variables
// Field names are converted from camelCase to UPPER_SNAKE_CASE
func CreateEnvMap[T any](env T) pulumi.StringMap {
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
	if len(s) == 0 {
		return ""
	}

	var result strings.Builder
	runes := []rune(s)

	for i, currentRune := range runes {
		upperRune := unicode.ToUpper(currentRune)

		// Determine if we need to add an underscore before this character
		shouldAddUnderscore := false

		if i > 0 {
			prevRune := runes[i-1]

			// Pattern 1: lowercase followed by uppercase (word boundary)
			// "camelCase" -> "CAMEL_CASE"
			if unicode.IsUpper(currentRune) && unicode.IsLower(prevRune) {
				shouldAddUnderscore = true
			}

			// Pattern 2: End of acronym detection
			// "APIVersion" -> "API_VERSION" (at 'V')
			if unicode.IsUpper(currentRune) && i < len(runes)-1 {
				nextRune := runes[i+1]
				isEndOfAcronym := unicode.IsLower(nextRune) &&
					i > 1 &&
					unicode.IsUpper(prevRune)

				if isEndOfAcronym {
					shouldAddUnderscore = true
				}
			}

			// Pattern 3: number followed by letter (word boundary)
			// "PylonS3Region" -> "PYLON_S3_REGION"
			if unicode.IsLetter(currentRune) && unicode.IsDigit(prevRune) {
				shouldAddUnderscore = true
			}
		}

		if shouldAddUnderscore {
			result.WriteRune('_')
		}
		result.WriteRune(upperRune)
	}

	return result.String()
}

// ParsePortWithDefault converts a port string to an integer with a default fallback
// This provides consistent port parsing across all components
func ParsePortWithDefault(portStr pulumi.StringInput, defaultPort int) pulumi.IntOutput {
	return pulumi.All(portStr).ApplyT(func(inputs []interface{}) int {
		portStr := inputs[0].(string)
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
		// Return default if there's an error parsing the port
		return defaultPort
	}).(pulumi.IntOutput)
}

// CreateResourceRequirements creates consistent resource requirements for pods
// Provides default values if any parameter is empty
func CreateResourceRequirements(cpuLimit, memoryLimit, cpuRequest, memoryRequest string) *corev1.ResourceRequirementsArgs {
	// Default values
	if cpuLimit == "" {
		cpuLimit = "2"
	}
	if memoryLimit == "" {
		memoryLimit = "2Gi"
	}
	if cpuRequest == "" {
		cpuRequest = "1"
	}
	if memoryRequest == "" {
		memoryRequest = "1Gi"
	}

	return &corev1.ResourceRequirementsArgs{
		Limits: pulumi.StringMap{
			"cpu":    pulumi.String(cpuLimit),
			"memory": pulumi.String(memoryLimit),
		},
		Requests: pulumi.StringMap{
			"cpu":    pulumi.String(cpuRequest),
			"memory": pulumi.String(memoryRequest),
		},
	}
}

// CreatePersistentVolumeClaim creates a PVC with consistent labeling and defaults
func CreatePersistentVolumeClaim(
	ctx *pulumi.Context,
	name string,
	namespace pulumi.StringInput,
	storageSize pulumi.StringInput,
	storageClass string,
	labels pulumi.StringMap,
	component pulumi.Resource,
) (*corev1.PersistentVolumeClaim, error) {
	if storageSize == nil {
		storageSize = pulumi.String("150Gi") // Default storage size
	}

	if storageClass == "" {
		storageClass = "aws-gp3" // Default storage class
	}

	if labels == nil {
		labels = CreateResourceLabels(name, name, name, nil)
	}

	return corev1.NewPersistentVolumeClaim(ctx, name, &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Labels:    labels,
			Namespace: namespace,
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{pulumi.String("ReadWriteOnce")},
			Resources: &corev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": storageSize,
				},
			},
			StorageClassName: pulumi.String(storageClass),
		},
	}, pulumi.Parent(component))
}
