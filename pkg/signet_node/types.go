package signet_node

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	crd "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AppLabels represents the Kubernetes labels to be applied to the signet node resources.
type AppLabels struct {
	Labels pulumi.StringMap
}

// Public-facing structs with base Go types
type SignetNodeComponentArgs struct {
	Name                        string
	Namespace                   string
	Env                         SignetNodeEnv
	ExecutionJwt                string
	ExecutionPvcSize            string
	LighthousePvcSize           string
	RollupPvcSize               string
	ExecutionClientImage        string
	ConsensusClientImage        string
	ExecutionClientStartCommand []string
	ConsensusClientStartCommand []string
	AppLabels                   AppLabels
	SignetNodeDataMountPath     string // Optional: defaults to "/root/.local/share/reth"
	RollupDataMountPath         string // Optional: defaults to "/root/.local/share/exex"
	ExecutionJwtMountPath       string // Optional: defaults to "/etc/reth/execution-jwt"
}

// Internal structs with Pulumi types for use within the component
type signetNodeComponentArgsInternal struct {
	Name                        string
	Namespace                   pulumi.StringInput
	Env                         signetNodeEnvInternal
	ExecutionJwt                pulumi.StringInput
	ExecutionPvcSize            pulumi.StringInput
	LighthousePvcSize           pulumi.StringInput
	RollupPvcSize               pulumi.StringInput
	ExecutionClientImage        pulumi.StringInput
	ConsensusClientImage        pulumi.StringInput
	ExecutionClientStartCommand pulumi.StringArrayInput
	ConsensusClientStartCommand pulumi.StringArrayInput
	AppLabels                   AppLabels
	SignetNodeDataMountPath     pulumi.StringInput
	RollupDataMountPath         pulumi.StringInput
	ExecutionJwtMountPath       pulumi.StringInput
}

type SignetNodeComponent struct {
	pulumi.ResourceState

	SignetNodeComponentArgs  SignetNodeComponentArgs
	SignetNodeStatefulSet    *appsv1.StatefulSet
	LighthouseStatefulSet    *appsv1.StatefulSet
	SignetNodeService        *corev1.Service
	LighthouseService        *corev1.Service
	SignetNodeConfigMap      *corev1.ConfigMap
	LighthouseConfigMap      *corev1.ConfigMap
	JwtSecret                *corev1.Secret
	HostDatabasePvc          *corev1.PersistentVolumeClaim
	RollupDatabasePvc        *corev1.PersistentVolumeClaim
	LighthousePvc            *corev1.PersistentVolumeClaim
	SignetNodeVirtualService *crd.CustomResource
}

// Public-facing environment struct with base Go types
type SignetNodeEnv struct {
	ChainName          string `pulumi:"chainName" validate:"required"`
	IpcEndpoint        string `pulumi:"ipcEndpoint" validate:"required"`
	RpcPort            int    `pulumi:"rpcPort" validate:"required"`
	RustLog            string `pulumi:"rustLog"` // optional, defaults to "info"
	SignetChainId      int    `pulumi:"signetChainId" validate:"required"`
	SignetClUrl        string `pulumi:"signetClUrl" validate:"required"`
	SignetDatabasePath string `pulumi:"signetDatabasePath" validate:"required"`
	SignetPylonUrl     string `pulumi:"signetPylonUrl" validate:"required"`
	SignetStaticPath   string `pulumi:"signetStaticPath" validate:"required"`
	TxForwardUrl       string `pulumi:"txForwardUrl" validate:"required"`
	WsRpcPort          int    `pulumi:"wsRpcPort" validate:"required"`
}

// Internal environment struct with Pulumi types
type signetNodeEnvInternal struct {
	ChainName          pulumi.StringInput `pulumi:"chainName" validate:"required"`
	IpcEndpoint        pulumi.StringInput `pulumi:"ipcEndpoint" validate:"required"`
	RpcPort            pulumi.IntInput    `pulumi:"rpcPort" validate:"required"`
	RustLog            pulumi.StringInput `pulumi:"rustLog"` // now StringInput, not StringPtrInput
	SignetChainId      pulumi.IntInput    `pulumi:"signetChainId" validate:"required"`
	SignetClUrl        pulumi.StringInput `pulumi:"signetClUrl" validate:"required"`
	SignetDatabasePath pulumi.StringInput `pulumi:"signetDatabasePath" validate:"required"`
	SignetPylonUrl     pulumi.StringInput `pulumi:"signetPylonUrl" validate:"required"`
	SignetStaticPath   pulumi.StringInput `pulumi:"signetStaticPath" validate:"required"`
	TxForwardUrl       pulumi.StringInput `pulumi:"txForwardUrl" validate:"required"`
	WsRpcPort          pulumi.IntInput    `pulumi:"wsRpcPort" validate:"required"`
}

// Conversion function to convert public args to internal args
func (args SignetNodeComponentArgs) toInternal() signetNodeComponentArgsInternal {
	return signetNodeComponentArgsInternal{
		Name:                        args.Name,
		Namespace:                   pulumi.String(args.Namespace),
		Env:                         args.Env.toInternal(),
		ExecutionJwt:                pulumi.String(args.ExecutionJwt),
		ExecutionPvcSize:            pulumi.String(args.ExecutionPvcSize),
		LighthousePvcSize:           pulumi.String(args.LighthousePvcSize),
		RollupPvcSize:               pulumi.String(args.RollupPvcSize),
		ExecutionClientImage:        pulumi.String(args.ExecutionClientImage),
		ConsensusClientImage:        pulumi.String(args.ConsensusClientImage),
		ExecutionClientStartCommand: pulumi.ToStringArray(args.ExecutionClientStartCommand),
		ConsensusClientStartCommand: pulumi.ToStringArray(args.ConsensusClientStartCommand),
		AppLabels:                   args.AppLabels,
		SignetNodeDataMountPath:     pulumi.String(args.SignetNodeDataMountPath),
		RollupDataMountPath:         pulumi.String(args.RollupDataMountPath),
		ExecutionJwtMountPath:       pulumi.String(args.ExecutionJwtMountPath),
	}
}

// Conversion function to convert public env to internal env
func (e SignetNodeEnv) toInternal() signetNodeEnvInternal {
	return signetNodeEnvInternal{
		ChainName:          pulumi.String(e.ChainName),
		IpcEndpoint:        pulumi.String(e.IpcEndpoint),
		RpcPort:            pulumi.Int(e.RpcPort),
		RustLog:            pulumi.String(e.RustLog),
		SignetChainId:      pulumi.Int(e.SignetChainId),
		SignetClUrl:        pulumi.String(e.SignetClUrl),
		SignetDatabasePath: pulumi.String(e.SignetDatabasePath),
		SignetPylonUrl:     pulumi.String(e.SignetPylonUrl),
		SignetStaticPath:   pulumi.String(e.SignetStaticPath),
		TxForwardUrl:       pulumi.String(e.TxForwardUrl),
		WsRpcPort:          pulumi.Int(e.WsRpcPort),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal env
func (e signetNodeEnvInternal) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// ConsensusEnv contains environment variables for the consensus client
type ConsensusEnv struct {
	Example string `pulumi:"example"`
}

// Internal consensus env with Pulumi types
type consensusEnvInternal struct {
	Example pulumi.StringInput `pulumi:"example"`
}

// Conversion function for consensus env
func (e ConsensusEnv) toInternal() consensusEnvInternal {
	return consensusEnvInternal{
		Example: pulumi.String(e.Example),
	}
}

// GetEnvMap implements the utils.EnvProvider interface for internal consensus env
func (e consensusEnvInternal) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// SignetNode interface defines methods that the SignetNodeComponent must implement
type SignetNode interface {
}

// Ensure SignetNodeComponent implements SignetNode
var _ SignetNode = &SignetNodeComponent{}
