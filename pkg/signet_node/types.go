package signet_node

import (
	"strconv"

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
	// Always required
	BlobExplorerUrl    string `pulumi:"blobExplorerUrl" validate:"required"`
	SignetStaticPath   string `pulumi:"signetStaticPath" validate:"required"`
	SignetDatabasePath string `pulumi:"signetDatabasePath" validate:"required"`

	// Optional with defaults (RpcPort: 5959, WsRpcPort: 5960)
	RpcPort   int `pulumi:"rpcPort"`
	WsRpcPort int `pulumi:"wsRpcPort"`

	// Optional (no defaults)
	ChainName      string `pulumi:"chainName"`
	SignetClUrl    string `pulumi:"signetClUrl"`
	SignetPylonUrl string `pulumi:"signetPylonUrl"`
	TxForwardUrl   string `pulumi:"txForwardUrl"`
	IpcEndpoint    string `pulumi:"ipcEndpoint"`
	RustLog        string `pulumi:"rustLog"`

	// Conditional: required if ChainName is not set
	RollupGenesisJsonPath string `pulumi:"rollupGenesisJsonPath"`
	HostGenesisJsonPath   string `pulumi:"hostGenesisJsonPath"`
	StartTimestamp        int    `pulumi:"startTimestamp"`
	SlotOffset            int    `pulumi:"slotOffset"`
	SlotDuration          int    `pulumi:"slotDuration"`
}

// Internal environment struct with Pulumi types
type signetNodeEnvInternal struct {
	// Always required
	BlobExplorerUrl    pulumi.StringInput `pulumi:"blobExplorerUrl" validate:"required"`
	SignetStaticPath   pulumi.StringInput `pulumi:"signetStaticPath" validate:"required"`
	SignetDatabasePath pulumi.StringInput `pulumi:"signetDatabasePath" validate:"required"`

	// Optional with defaults
	RpcPort   pulumi.StringInput `pulumi:"rpcPort"`
	WsRpcPort pulumi.StringInput `pulumi:"wsRpcPort"`

	// Optional (no defaults)
	ChainName      pulumi.StringInput `pulumi:"chainName"`
	SignetClUrl    pulumi.StringInput `pulumi:"signetClUrl"`
	SignetPylonUrl pulumi.StringInput `pulumi:"signetPylonUrl"`
	TxForwardUrl   pulumi.StringInput `pulumi:"txForwardUrl"`
	IpcEndpoint    pulumi.StringInput `pulumi:"ipcEndpoint"`
	RustLog        pulumi.StringInput `pulumi:"rustLog"`

	// Conditional: required if ChainName is not set
	RollupGenesisJsonPath pulumi.StringInput `pulumi:"rollupGenesisJsonPath"`
	HostGenesisJsonPath   pulumi.StringInput `pulumi:"hostGenesisJsonPath"`
	StartTimestamp        pulumi.StringInput `pulumi:"startTimestamp"`
	SlotOffset            pulumi.StringInput `pulumi:"slotOffset"`
	SlotDuration          pulumi.StringInput `pulumi:"slotDuration"`
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
		// Always required
		BlobExplorerUrl:    pulumi.String(e.BlobExplorerUrl),
		SignetStaticPath:   pulumi.String(e.SignetStaticPath),
		SignetDatabasePath: pulumi.String(e.SignetDatabasePath),

		// Optional with defaults
		RpcPort:   pulumi.String(strconv.Itoa(e.RpcPort)),
		WsRpcPort: pulumi.String(strconv.Itoa(e.WsRpcPort)),

		// Optional (no defaults)
		ChainName:      pulumi.String(e.ChainName),
		SignetClUrl:    pulumi.String(e.SignetClUrl),
		SignetPylonUrl: pulumi.String(e.SignetPylonUrl),
		TxForwardUrl:   pulumi.String(e.TxForwardUrl),
		IpcEndpoint:    pulumi.String(e.IpcEndpoint),
		RustLog:        pulumi.String(e.RustLog),

		// Conditional fields
		RollupGenesisJsonPath: pulumi.String(e.RollupGenesisJsonPath),
		HostGenesisJsonPath:   pulumi.String(e.HostGenesisJsonPath),
		StartTimestamp:        pulumi.String(strconv.Itoa(e.StartTimestamp)),
		SlotOffset:            pulumi.String(strconv.Itoa(e.SlotOffset)),
		SlotDuration:          pulumi.String(strconv.Itoa(e.SlotDuration)),
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
