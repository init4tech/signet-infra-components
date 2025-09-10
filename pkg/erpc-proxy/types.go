package erpcproxy

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErpcProxyComponent represents a Pulumi component that deploys an eRPC proxy service
type ErpcProxyComponent struct {
	pulumi.ResourceState
	ServiceAccount *corev1.ServiceAccount
	ConfigMap      *corev1.ConfigMap
	Secret         *corev1.Secret
	Deployment     *appsv1.Deployment
	Service        *corev1.Service
}

// ErpcProxyComponentArgs represents the public-facing arguments for the eRPC proxy component
type ErpcProxyComponentArgs struct {
	Namespace string             `pulumi:"namespace" validate:"required"`
	Name      string             `pulumi:"name" validate:"required"`
	Image     string             `pulumi:"image" validate:"required"`
	Config    ErpcProxyConfig    `pulumi:"config" validate:"required"`
	ApiKeys   map[string]string  `pulumi:"apiKeys"`
	Resources ErpcProxyResources `pulumi:"resources"`
	Replicas  int                `pulumi:"replicas"`
}

// erpcProxyComponentArgsInternal represents the internal arguments with Pulumi types
type erpcProxyComponentArgsInternal struct {
	Namespace pulumi.StringInput
	Name      pulumi.StringInput
	Image     pulumi.StringInput
	Config    erpcProxyConfigInternal
	ApiKeys   pulumi.StringMap
	Resources erpcProxyResourcesInternal
	Replicas  pulumi.IntInput
}

// ErpcProxyConfig represents the eRPC proxy configuration
type ErpcProxyConfig struct {
	LogLevel string                   `pulumi:"logLevel"`
	Database ErpcProxyDatabaseConfig  `pulumi:"database"`
	Server   ErpcProxyServerConfig    `pulumi:"server"`
	Projects []ErpcProxyProjectConfig `pulumi:"projects" validate:"required,min=1"`
}

// erpcProxyConfigInternal represents the internal config with Pulumi types
type erpcProxyConfigInternal struct {
	LogLevel pulumi.StringInput
	Database erpcProxyDatabaseConfigInternal
	Server   erpcProxyServerConfigInternal
	Projects pulumi.StringInput // Will be marshaled to YAML
}

// ErpcProxyDatabaseConfig represents database configuration
type ErpcProxyDatabaseConfig struct {
	Type          string `pulumi:"type"`
	ConnectionUrl string `pulumi:"connectionUrl"`
}

// erpcProxyDatabaseConfigInternal represents internal database config
type erpcProxyDatabaseConfigInternal struct {
	Type          pulumi.StringInput
	ConnectionUrl pulumi.StringInput
}

// ErpcProxyServerConfig represents server configuration
type ErpcProxyServerConfig struct {
	HttpHostV4 string `pulumi:"httpHostV4"`
	HttpPortV4 int    `pulumi:"httpPortV4"`
	MaxTimeout string `pulumi:"maxTimeout"`
}

// erpcProxyServerConfigInternal represents internal server config
type erpcProxyServerConfigInternal struct {
	HttpHostV4 pulumi.StringInput
	HttpPortV4 pulumi.IntInput
	MaxTimeout pulumi.StringInput
}

// ErpcProxyProjectConfig represents a project configuration
type ErpcProxyProjectConfig struct {
	Id              string                    `pulumi:"id" validate:"required"`
	Networks        []ErpcProxyNetworkConfig  `pulumi:"networks" validate:"required,min=1"`
	Upstreams       []ErpcProxyUpstreamConfig `pulumi:"upstreams" validate:"required,min=1"`
	RateLimitBudget string                    `pulumi:"rateLimitBudget"`
	Cors            *ErpcProxyCorsConfig      `pulumi:"cors"`
}

// ErpcProxyCorsConfig represents CORS configuration for a project
type ErpcProxyCorsConfig struct {
	AllowedOrigins   []string `pulumi:"allowedOrigins"`
	AllowedMethods   []string `pulumi:"allowedMethods"`
	AllowedHeaders   []string `pulumi:"allowedHeaders"`
	ExposedHeaders   []string `pulumi:"exposedHeaders"`
	AllowCredentials bool     `pulumi:"allowCredentials"`
	MaxAge           int      `pulumi:"maxAge" validate:"min=0"`
}

// ErpcProxyNetworkConfig represents a network configuration
type ErpcProxyNetworkConfig struct {
	ChainId      int                     `pulumi:"chainId" validate:"required"`
	Architecture string                  `pulumi:"architecture" validate:"required"`
	Failover     ErpcProxyFailoverConfig `pulumi:"failover"`
}

// ErpcProxyFailoverConfig represents failover configuration
type ErpcProxyFailoverConfig struct {
	MaxRetries    int    `pulumi:"maxRetries"`
	BackoffMs     int    `pulumi:"backoffMs"`
	BackoffMaxMs  int    `pulumi:"backoffMaxMs"`
	BackoffFactor int    `pulumi:"backoffFactor"`
	Duration      string `pulumi:"duration"`
}

// ErpcProxyUpstreamConfig represents an upstream configuration
type ErpcProxyUpstreamConfig struct {
	Id              string `pulumi:"id" validate:"required"`
	Type            string `pulumi:"type" validate:"required"`
	Endpoint        string `pulumi:"endpoint" validate:"required"`
	RateLimitBudget string `pulumi:"rateLimitBudget"`
	MaxRetries      int    `pulumi:"maxRetries"`
	Timeout         string `pulumi:"timeout"`
}

// ErpcProxyResources represents resource requirements
type ErpcProxyResources struct {
	MemoryRequest string `pulumi:"memoryRequest"`
	MemoryLimit   string `pulumi:"memoryLimit"`
	CpuRequest    string `pulumi:"cpuRequest"`
	CpuLimit      string `pulumi:"cpuLimit"`
}

// erpcProxyResourcesInternal represents internal resource requirements
type erpcProxyResourcesInternal struct {
	MemoryRequest pulumi.StringInput
	MemoryLimit   pulumi.StringInput
	CpuRequest    pulumi.StringInput
	CpuLimit      pulumi.StringInput
}

// ErpcProxyEnv represents environment variables for the eRPC proxy
type ErpcProxyEnv struct {
	GoGC       string `pulumi:"goGC"`
	GoMemLimit string `pulumi:"goMemLimit"`
}

// erpcProxyEnvInternal represents internal environment variables
type erpcProxyEnvInternal struct {
	GoGC       pulumi.StringInput
	GoMemLimit pulumi.StringInput
}

// toInternal converts public args to internal args
func (args ErpcProxyComponentArgs) toInternal() erpcProxyComponentArgsInternal {
	// Convert ApiKeys map to pulumi.StringMap
	apiKeysMap := make(pulumi.StringMap)
	for k, v := range args.ApiKeys {
		apiKeysMap[k] = pulumi.String(v)
	}

	return erpcProxyComponentArgsInternal{
		Namespace: pulumi.String(args.Namespace),
		Name:      pulumi.String(args.Name),
		Image:     pulumi.String(args.Image),
		Config:    args.Config.toInternal(),
		ApiKeys:   apiKeysMap,
		Resources: args.Resources.toInternal(),
		Replicas:  pulumi.Int(args.Replicas),
	}
}

// toInternal converts public config to internal config
func (c ErpcProxyConfig) toInternal() erpcProxyConfigInternal {
	// Will be implemented to marshal projects to YAML
	return erpcProxyConfigInternal{
		LogLevel: pulumi.String(c.LogLevel),
		Database: c.Database.toInternal(),
		Server:   c.Server.toInternal(),
		Projects: pulumi.String(""), // Will be properly marshaled in the component
	}
}

// toInternal converts public database config to internal
func (d ErpcProxyDatabaseConfig) toInternal() erpcProxyDatabaseConfigInternal {
	return erpcProxyDatabaseConfigInternal{
		Type:          pulumi.String(d.Type),
		ConnectionUrl: pulumi.String(d.ConnectionUrl),
	}
}

// toInternal converts public server config to internal
func (s ErpcProxyServerConfig) toInternal() erpcProxyServerConfigInternal {
	return erpcProxyServerConfigInternal{
		HttpHostV4: pulumi.String(s.HttpHostV4),
		HttpPortV4: pulumi.Int(s.HttpPortV4),
		MaxTimeout: pulumi.String(s.MaxTimeout),
	}
}

// toInternal converts public resources to internal
func (r ErpcProxyResources) toInternal() erpcProxyResourcesInternal {
	return erpcProxyResourcesInternal{
		MemoryRequest: pulumi.String(r.MemoryRequest),
		MemoryLimit:   pulumi.String(r.MemoryLimit),
		CpuRequest:    pulumi.String(r.CpuRequest),
		CpuLimit:      pulumi.String(r.CpuLimit),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal env
func (e erpcProxyEnvInternal) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}
