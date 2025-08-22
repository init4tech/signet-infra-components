package signet_node

import (
	"fmt"

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
