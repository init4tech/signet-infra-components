// Package erpcproxy provides a Pulumi component for deploying an eRPC proxy service to Kubernetes.
package erpcproxy

import (
	"fmt"

	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

// NewErpcProxy creates a new eRPC proxy component with the given configuration.
func NewErpcProxy(ctx *pulumi.Context, args ErpcProxyComponentArgs, opts ...pulumi.ResourceOption) (*ErpcProxyComponent, error) {
	// Apply defaults
	if args.Image == "" {
		// Log that using default image
		ctx.Log.Info(fmt.Sprintf("Using default image: %s", DefaultImage), nil)
		args.Image = DefaultImage
	}
	if args.Config.LogLevel == "" {
		ctx.Log.Info(fmt.Sprintf("Using default log level: %s", DefaultLogLevel), nil)
		args.Config.LogLevel = DefaultLogLevel
	}
	if args.Config.Server.HttpPortV4 == 0 {
		ctx.Log.Info(fmt.Sprintf("Using default HTTP port: %d", DefaultHttpPort), nil)
		args.Config.Server.HttpPortV4 = DefaultHttpPort
	}
	if args.Config.Server.MaxTimeout == "" {
		ctx.Log.Info(fmt.Sprintf("Using default max timeout: %s", DefaultMaxTimeoutMs), nil)
		args.Config.Server.MaxTimeout = DefaultMaxTimeoutMs
	}
	if args.Resources.MemoryRequest == "" {
		ctx.Log.Info(fmt.Sprintf("Using default memory request: %s", DefaultMemoryRequest), nil)
		args.Resources.MemoryRequest = DefaultMemoryRequest
	}
	if args.Resources.MemoryLimit == "" {
		ctx.Log.Info(fmt.Sprintf("Using default memory limit: %s", DefaultMemoryLimit), nil)
		args.Resources.MemoryLimit = DefaultMemoryLimit
	}
	if args.Resources.CpuRequest == "" {
		ctx.Log.Info(fmt.Sprintf("Using default CPU request: %s", DefaultCpuRequest), nil)
		args.Resources.CpuRequest = DefaultCpuRequest
	}
	if args.Resources.CpuLimit == "" {
		ctx.Log.Info(fmt.Sprintf("Using default CPU limit: %s", DefaultCpuLimit), nil)
		args.Resources.CpuLimit = DefaultCpuLimit
	}
	if args.Replicas == 0 {
		ctx.Log.Info(fmt.Sprintf("Using default replicas: %d", DefaultReplicas), nil)
		args.Replicas = DefaultReplicas
	}

	// Validate arguments
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid erpc proxy component args: %w", err)
	}

	// Convert public args to internal args for use with Pulumi
	internalArgs := args.toInternal()

	component := &ErpcProxyComponent{}
	err := ctx.RegisterComponentResource(ComponentKind, args.Name, component)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %w", err)
	}

	// Create service account
	serviceAccountName := fmt.Sprintf("%s%s", args.Name, ServiceAccountSuffix)
	sa, err := corev1.NewServiceAccount(ctx, serviceAccountName, &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceAccountName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, serviceAccountName, args.Name, nil),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}
	component.ServiceAccount = sa

	// Create ConfigMap for eRPC configuration
	configMapName := fmt.Sprintf("%s%s", args.Name, ConfigMapSuffix)
	configYaml, err := marshalErpcConfig(args.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal eRPC config: %w", err)
	}

	configMap, err := corev1.NewConfigMap(ctx, configMapName, &corev1.ConfigMapArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(configMapName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, configMapName, args.Name, nil),
		},
		Data: pulumi.StringMap{
			ConfigFileName: pulumi.String(configYaml),
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create config map: %w", err)
	}
	component.ConfigMap = configMap

	// Create Secret for API keys if provided
	if len(args.ApiKeys) > 0 {
		secretName := fmt.Sprintf("%s%s", args.Name, SecretSuffix)
		secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(secretName),
				Namespace: internalArgs.Namespace,
				Labels:    utils.CreateResourceLabels(args.Name, secretName, args.Name, nil),
			},
			StringData: internalArgs.ApiKeys,
		}, pulumi.Parent(component))
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %w", err)
		}
		component.Secret = secret
	}

	// Create pod labels
	podLabels := utils.CreateResourceLabels(args.Name, args.Name, args.Name, nil)
	podLabels["app"] = pulumi.String(args.Name)

	// Build environment variables
	envVars := corev1.EnvVarArray{
		&corev1.EnvVarArgs{
			Name:  pulumi.String(EnvGoGC),
			Value: pulumi.String(DefaultGoGC),
		},
		&corev1.EnvVarArgs{
			Name:  pulumi.String(EnvGoMemLimit),
			Value: pulumi.String(DefaultGoMemLimit),
		},
		&corev1.EnvVarArgs{
			Name:  pulumi.String("CONFIG_FILE"),
			Value: pulumi.Sprintf("%s/%s", ConfigMountPath, ConfigFileName),
		},
	}

	// Build envFrom sources
	envFromSources := corev1.EnvFromSourceArray{}
	if component.Secret != nil {
		envFromSources = append(envFromSources, &corev1.EnvFromSourceArgs{
			SecretRef: &corev1.SecretEnvSourceArgs{
				Name: component.Secret.Metadata.Name(),
			},
		})
	}

	// Create deployment
	deploymentName := fmt.Sprintf("%s%s", args.Name, DeploymentSuffix)
	deployment, err := appsv1.NewDeployment(ctx, deploymentName, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(deploymentName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, deploymentName, args.Name, nil),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: internalArgs.Replicas,
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: podLabels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: podLabels,
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String(serviceAccountName),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:    pulumi.String(args.Name),
							Image:   internalArgs.Image,
							Args:    pulumi.StringArray{pulumi.Sprintf("%s/%s", ConfigMountPath, ConfigFileName)},
							Command: pulumi.StringArray{pulumi.String("/erpc-server")},
							Env:     envVars,
							EnvFrom: envFromSources,
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("http"),
									ContainerPort: pulumi.Int(args.Config.Server.HttpPortV4),
									Protocol:      pulumi.String("TCP"),
								},
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("metrics"),
									ContainerPort: pulumi.Int(DefaultMetricsPort),
									Protocol:      pulumi.String("TCP"),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"memory": internalArgs.Resources.MemoryRequest,
									"cpu":    internalArgs.Resources.CpuRequest,
								},
								Limits: pulumi.StringMap{
									"memory": internalArgs.Resources.MemoryLimit,
									"cpu":    internalArgs.Resources.CpuLimit,
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("config"),
									MountPath: pulumi.String(ConfigMountPath),
									ReadOnly:  pulumi.Bool(false),
								},
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: pulumi.Int(args.Config.Server.HttpPortV4),
								},
								InitialDelaySeconds: pulumi.Int(30),
								PeriodSeconds:       pulumi.Int(10),
								TimeoutSeconds:      pulumi.Int(5),
								FailureThreshold:    pulumi.Int(3),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: pulumi.Int(args.Config.Server.HttpPortV4),
								},
								InitialDelaySeconds: pulumi.Int(10),
								PeriodSeconds:       pulumi.Int(5),
								TimeoutSeconds:      pulumi.Int(3),
								FailureThreshold:    pulumi.Int(3),
							},
							StartupProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthcheck"),
									Port: pulumi.Int(args.Config.Server.HttpPortV4),
								},
								InitialDelaySeconds: pulumi.Int(0),
								PeriodSeconds:       pulumi.Int(10),
								TimeoutSeconds:      pulumi.Int(5),
								FailureThreshold:    pulumi.Int(30),
							},
						},
					},
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name: pulumi.String("config"),
							ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
								Name:        component.ConfigMap.Metadata.Name(),
								DefaultMode: pulumi.Int(0644),
							},
						},
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}
	component.Deployment = deployment

	// Create service
	serviceName := fmt.Sprintf("%s%s", args.Name, ServiceSuffix)
	service, err := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: internalArgs.Namespace,
			Labels:    utils.CreateResourceLabels(args.Name, serviceName, args.Name, nil),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type:     pulumi.String("ClusterIP"),
			Selector: podLabels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					Port:       pulumi.Int(args.Config.Server.HttpPortV4),
					TargetPort: pulumi.String("http"),
					Protocol:   pulumi.String("TCP"),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("metrics"),
					Port:       pulumi.Int(DefaultMetricsPort),
					TargetPort: pulumi.String("metrics"),
					Protocol:   pulumi.String("TCP"),
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	component.Service = service

	return component, nil
}

// marshalErpcConfig converts the eRPC config to YAML format
func marshalErpcConfig(config ErpcProxyConfig) (string, error) {
	// Build the config structure for YAML marshaling
	configMap := map[string]interface{}{
		"logLevel": config.LogLevel,
	}

	// Add server config if provided
	if config.Server.HttpHostV4 != "" || config.Server.HttpPortV4 > 0 || config.Server.MaxTimeout != "" {
		serverMap := map[string]interface{}{}
		if config.Server.HttpHostV4 != "" {
			serverMap["httpHostV4"] = config.Server.HttpHostV4
		}
		if config.Server.HttpPortV4 > 0 {
			serverMap["httpPortV4"] = config.Server.HttpPortV4
		}
		if config.Server.MaxTimeout != "" {
			serverMap["maxTimeout"] = config.Server.MaxTimeout
		}
		if len(serverMap) > 0 {
			configMap["server"] = serverMap
		}
	}

	// Add database config if provided
	if config.Database.Type != "" {
		configMap["database"] = map[string]interface{}{
			"type":          config.Database.Type,
			"connectionUrl": config.Database.ConnectionUrl,
		}
	}

	// Build projects array
	projects := make([]map[string]interface{}, 0, len(config.Projects))
	for _, project := range config.Projects {
		projectMap := map[string]interface{}{
			"id": project.Id,
		}

		if project.RateLimitBudget != "" {
			projectMap["rateLimitBudget"] = project.RateLimitBudget
		}

		// Build upstreams array
		upstreams := make([]map[string]interface{}, 0, len(project.Upstreams))
		for _, upstream := range project.Upstreams {
			upstreamMap := map[string]interface{}{
				"id":       upstream.Id,
				"type":     upstream.Type,
				"endpoint": upstream.Endpoint,
			}
			if upstream.RateLimitBudget != "" {
				upstreamMap["rateLimitBudget"] = upstream.RateLimitBudget
			}
			if upstream.MaxRetries > 0 {
				upstreamMap["maxRetries"] = upstream.MaxRetries
			}
			if upstream.Timeout != "" {
				upstreamMap["timeout"] = upstream.Timeout
			}
			upstreams = append(upstreams, upstreamMap)
		}
		projectMap["upstreams"] = upstreams

		// Build networks array
		networks := make([]map[string]interface{}, 0, len(project.Networks))
		for _, network := range project.Networks {
			networkMap := map[string]interface{}{
				"evm": map[string]interface{}{
					"chainId": network.ChainId,
				},
				"architecture": network.Architecture,
			}

			// Add failover config if provided
			if network.Failover.MaxRetries > 0 || network.Failover.BackoffMs > 0 {
				failover := map[string]interface{}{}
				if network.Failover.MaxRetries > 0 {
					failover["maxRetries"] = network.Failover.MaxRetries
				}
				if network.Failover.BackoffMs > 0 {
					failover["backoffMs"] = network.Failover.BackoffMs
				}
				if network.Failover.BackoffMaxMs > 0 {
					failover["backoffMaxMs"] = network.Failover.BackoffMaxMs
				}
				if network.Failover.BackoffFactor > 0 {
					failover["backoffFactor"] = network.Failover.BackoffFactor
				}
				if network.Failover.Duration != "" {
					failover["duration"] = network.Failover.Duration
				}
				networkMap["failover"] = failover
			}

			networks = append(networks, networkMap)
		}
		projectMap["networks"] = networks

		// Add CORS config if provided
		if project.Cors != nil {
			corsMap := map[string]interface{}{}

			if len(project.Cors.AllowedOrigins) > 0 {
				corsMap["allowedOrigins"] = project.Cors.AllowedOrigins
			}
			if len(project.Cors.AllowedMethods) > 0 {
				corsMap["allowedMethods"] = project.Cors.AllowedMethods
			}
			if len(project.Cors.AllowedHeaders) > 0 {
				corsMap["allowedHeaders"] = project.Cors.AllowedHeaders
			}
			if len(project.Cors.ExposedHeaders) > 0 {
				corsMap["exposedHeaders"] = project.Cors.ExposedHeaders
			}
			if project.Cors.AllowCredentials {
				corsMap["allowCredentials"] = project.Cors.AllowCredentials
			}
			if project.Cors.MaxAge > 0 {
				corsMap["maxAge"] = project.Cors.MaxAge
			}

			// Only add CORS section if there are any configured values
			if len(corsMap) > 0 {
				projectMap["cors"] = corsMap
			}
		}

		projects = append(projects, projectMap)
	}
	configMap["projects"] = projects

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// GetServiceURL returns the service URL for the eRPC proxy
func (c *ErpcProxyComponent) GetServiceURL() pulumi.StringOutput {
	return pulumi.Sprintf("http://%s:%d", c.Service.Metadata.Name(), DefaultHttpPort)
}

// GetMetricsURL returns the metrics URL for the eRPC proxy
func (c *ErpcProxyComponent) GetMetricsURL() pulumi.StringOutput {
	return pulumi.Sprintf("http://%s:%d/metrics", c.Service.Metadata.Name(), DefaultMetricsPort)
}
