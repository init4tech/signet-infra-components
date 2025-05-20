package signet_node

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreatePersistentVolumeClaim creates a new PVC with the given name and size
func CreatePersistentVolumeClaim(
	ctx *pulumi.Context,
	name string,
	namespace pulumi.StringInput,
	storageSize pulumi.StringInput,
	storageClass string,
	component pulumi.Resource,
) (*corev1.PersistentVolumeClaim, error) {
	if storageSize == nil {
		storageSize = pulumi.String(DefaultStorageSize)
	}

	if storageClass == "" {
		storageClass = DefaultStorageClass
	}

	pvc, err := corev1.NewPersistentVolumeClaim(ctx, name, &corev1.PersistentVolumeClaimArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.String(name),
				"app.kubernetes.io/part-of": pulumi.String(name),
			},
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

	if err != nil {
		return nil, fmt.Errorf("failed to create PVC %s: %w", name, err)
	}

	return pvc, nil
}

// CreateResourceLabels creates a consistent set of Kubernetes labels for resources
func CreateResourceLabels(name string) pulumi.StringMap {
	return pulumi.StringMap{
		"app":                       pulumi.String(name),
		"app.kubernetes.io/name":    pulumi.String(name),
		"app.kubernetes.io/part-of": pulumi.String(name),
	}
}

// CreateEnvironmentVars creates environment variables from the SignetNodeEnv struct
func CreateEnvironmentVars(env SignetNodeEnv) corev1.EnvVarArray {
	result := corev1.EnvVarArray{}

	// Process all string inputs from the struct's tags
	envVarMap := GetEnvironmentVarsFromStruct(env)
	for name, value := range envVarMap {
		result = append(result, &corev1.EnvVarArgs{
			Name:  pulumi.String(name),
			Value: value.(pulumi.StringInput),
		})
	}

	return result
}

// GetEnvironmentVarsFromStruct uses reflection to extract environment variables from struct tags
func GetEnvironmentVarsFromStruct(env SignetNodeEnv) map[string]pulumi.Input {
	result := make(map[string]pulumi.Input)

	t := reflect.TypeOf(env)
	v := reflect.ValueOf(env)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the field value
		fieldValue := v.Field(i).Interface()

		// Skip nil values
		if fieldValue == nil {
			continue
		}

		// Convert camelCase to SNAKE_CASE for env var name
		envName := CamelToSnake(field.Name)

		// Add to map
		result[envName] = fieldValue.(pulumi.Input)
	}

	return result
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

// GetResourceName returns a consistent name for a resource
func GetResourceName(baseName string, resourceType string) string {
	switch resourceType {
	case "service":
		return fmt.Sprintf("%s%s", baseName, ServiceSuffix)
	case "statefulset":
		return fmt.Sprintf("%s%s", baseName, StatefulSetSuffix)
	case "configmap":
		return fmt.Sprintf("%s%s", baseName, ConfigMapSuffix)
	case "pvc":
		return fmt.Sprintf("%s%s", baseName, PvcSuffix)
	case "secret":
		return fmt.Sprintf("%s%s", baseName, SecretSuffix)
	case "virtualservice":
		return fmt.Sprintf("%s%s", baseName, VirtualServiceSuffix)
	default:
		return baseName
	}
}

// ResourceRequirements returns consistent resource requirements for pods
func NewResourceRequirements(cpuLimit, memoryLimit, cpuRequest, memoryRequest string) *corev1.ResourceRequirementsArgs {
	if cpuLimit == "" {
		cpuLimit = DefaultCPULimit
	}
	if memoryLimit == "" {
		memoryLimit = DefaultMemoryLimit
	}
	if cpuRequest == "" {
		cpuRequest = DefaultCPURequest
	}
	if memoryRequest == "" {
		memoryRequest = DefaultMemoryRequest
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
