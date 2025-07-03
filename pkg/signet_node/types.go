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
	LighthousePvcSize           string
	RollupPvcSize               string
	ExecutionClientImage        string
	ConsensusClientImage        string
	ExecutionClientStartCommand []string
	ConsensusClientStartCommand []string
	AppLabels                   AppLabels
}

// Internal structs with Pulumi types for use within the component
type signetNodeComponentArgsInternal struct {
	Name                        string
	Namespace                   pulumi.StringInput
	Env                         signetNodeEnvInternal
	ExecutionJwt                pulumi.StringInput
	LighthousePvcSize           pulumi.StringInput
	RollupPvcSize               pulumi.StringInput
	ExecutionClientImage        pulumi.StringInput
	ConsensusClientImage        pulumi.StringInput
	ExecutionClientStartCommand pulumi.StringArrayInput
	ConsensusClientStartCommand pulumi.StringArrayInput
	AppLabels                   AppLabels
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
	HostZenithAddress             string `pulumi:"hostZenithAddress" validate:"required"`
	RuOrdersContractAddress       string `pulumi:"ruOrdersContractAddress" validate:"required"`
	HostOrdersContractAddress     string `pulumi:"hostOrdersContractAddress" validate:"required"`
	SignetChainId                 int    `pulumi:"signetChainId" validate:"required"`
	BlobExplorerUrl               string `pulumi:"blobExplorerUrl" validate:"required"`
	SignetStaticPath              string `pulumi:"signetStaticPath" validate:"required"`
	SignetDatabasePath            string `pulumi:"signetDatabasePath" validate:"required"`
	RustLog                       string `pulumi:"rustLog"`
	IpcEndpoint                   string `pulumi:"ipcEndpoint" validate:"required"`
	RpcPort                       int    `pulumi:"rpcPort" validate:"required"`
	WsRpcPort                     int    `pulumi:"wsRpcPort" validate:"required"`
	TxForwardUrl                  string `pulumi:"txForwardUrl" validate:"required"`
	GenesisJsonPath               string `pulumi:"genesisJsonPath" validate:"required"`
	HostZenithDeployHeight        string `pulumi:"hostZenithDeployHeight" validate:"required"`
	BaseFeeRecipient              string `pulumi:"baseFeeRecipient" validate:"required"`
	HostPassageContractAddress    string `pulumi:"hostPassageContractAddress" validate:"required"`
	HostTransactorContractAddress string `pulumi:"hostTransactorContractAddress" validate:"required"`
	RuPassageContractAddress      string `pulumi:"ruPassageContractAddress" validate:"required"`
	SignetClUrl                   string `pulumi:"signetClUrl" validate:"required"`
	SignetPylonUrl                string `pulumi:"signetPylonUrl" validate:"required"`
	HostStartTimestamp            int    `pulumi:"hostStartTimestamp" validate:"required"`
	HostSlotOffset                int    `pulumi:"hostSlotOffset" validate:"required"`
	HostSlotDuration              int    `pulumi:"hostSlotDuration" validate:"required"`
}

// Internal environment struct with Pulumi types
type signetNodeEnvInternal struct {
	HostZenithAddress             pulumi.StringInput `pulumi:"hostZenithAddress" validate:"required"`
	RuOrdersContractAddress       pulumi.StringInput `pulumi:"ruOrdersContractAddress" validate:"required"`
	HostOrdersContractAddress     pulumi.StringInput `pulumi:"hostOrdersContractAddress" validate:"required"`
	SignetChainId                 pulumi.StringInput `pulumi:"signetChainId" validate:"required"`
	BlobExplorerUrl               pulumi.StringInput `pulumi:"blobExplorerUrl" validate:"required"`
	SignetStaticPath              pulumi.StringInput `pulumi:"signetStaticPath" validate:"required"`
	SignetDatabasePath            pulumi.StringInput `pulumi:"signetDatabasePath" validate:"required"`
	RustLog                       pulumi.StringInput `pulumi:"rustLog"`
	IpcEndpoint                   pulumi.StringInput `pulumi:"ipcEndpoint" validate:"required"`
	RpcPort                       pulumi.StringInput `pulumi:"rpcPort" validate:"required"`
	WsRpcPort                     pulumi.StringInput `pulumi:"wsRpcPort" validate:"required"`
	TxForwardUrl                  pulumi.StringInput `pulumi:"txForwardUrl" validate:"required"`
	GenesisJsonPath               pulumi.StringInput `pulumi:"genesisJsonPath" validate:"required"`
	HostZenithDeployHeight        pulumi.StringInput `pulumi:"hostZenithDeployHeight" validate:"required"`
	BaseFeeRecipient              pulumi.StringInput `pulumi:"baseFeeRecipient" validate:"required"`
	HostPassageContractAddress    pulumi.StringInput `pulumi:"hostPassageContractAddress" validate:"required"`
	HostTransactorContractAddress pulumi.StringInput `pulumi:"hostTransactorContractAddress" validate:"required"`
	RuPassageContractAddress      pulumi.StringInput `pulumi:"ruPassageContractAddress" validate:"required"`
	SignetClUrl                   pulumi.StringInput `pulumi:"signetClUrl" validate:"required"`
	SignetPylonUrl                pulumi.StringInput `pulumi:"signetPylonUrl" validate:"required"`
	HostStartTimestamp            pulumi.StringInput `pulumi:"hostStartTimestamp" validate:"required"`
	HostSlotOffset                pulumi.StringInput `pulumi:"hostSlotOffset" validate:"required"`
	HostSlotDuration              pulumi.StringInput `pulumi:"hostSlotDuration" validate:"required"`
}

// Conversion function to convert public args to internal args
func (args SignetNodeComponentArgs) toInternal() signetNodeComponentArgsInternal {
	return signetNodeComponentArgsInternal{
		Name:                        args.Name,
		Namespace:                   pulumi.String(args.Namespace),
		Env:                         args.Env.toInternal(),
		ExecutionJwt:                pulumi.String(args.ExecutionJwt),
		LighthousePvcSize:           pulumi.String(args.LighthousePvcSize),
		RollupPvcSize:               pulumi.String(args.RollupPvcSize),
		ExecutionClientImage:        pulumi.String(args.ExecutionClientImage),
		ConsensusClientImage:        pulumi.String(args.ConsensusClientImage),
		ExecutionClientStartCommand: pulumi.ToStringArray(args.ExecutionClientStartCommand),
		ConsensusClientStartCommand: pulumi.ToStringArray(args.ConsensusClientStartCommand),
		AppLabels:                   args.AppLabels,
	}
}

// Conversion function to convert public env to internal env
func (e SignetNodeEnv) toInternal() signetNodeEnvInternal {
	return signetNodeEnvInternal{
		HostZenithAddress:             pulumi.String(e.HostZenithAddress),
		RuOrdersContractAddress:       pulumi.String(e.RuOrdersContractAddress),
		HostOrdersContractAddress:     pulumi.String(e.HostOrdersContractAddress),
		SignetChainId:                 pulumi.String(strconv.Itoa(e.SignetChainId)),
		BlobExplorerUrl:               pulumi.String(e.BlobExplorerUrl),
		SignetStaticPath:              pulumi.String(e.SignetStaticPath),
		SignetDatabasePath:            pulumi.String(e.SignetDatabasePath),
		RustLog:                       pulumi.String(e.RustLog),
		IpcEndpoint:                   pulumi.String(e.IpcEndpoint),
		RpcPort:                       pulumi.String(strconv.Itoa(e.RpcPort)),
		WsRpcPort:                     pulumi.String(strconv.Itoa(e.WsRpcPort)),
		TxForwardUrl:                  pulumi.String(e.TxForwardUrl),
		GenesisJsonPath:               pulumi.String(e.GenesisJsonPath),
		HostZenithDeployHeight:        pulumi.String(e.HostZenithDeployHeight),
		BaseFeeRecipient:              pulumi.String(e.BaseFeeRecipient),
		HostPassageContractAddress:    pulumi.String(e.HostPassageContractAddress),
		HostTransactorContractAddress: pulumi.String(e.HostTransactorContractAddress),
		RuPassageContractAddress:      pulumi.String(e.RuPassageContractAddress),
		SignetClUrl:                   pulumi.String(e.SignetClUrl),
		SignetPylonUrl:                pulumi.String(e.SignetPylonUrl),
		HostStartTimestamp:            pulumi.String(strconv.Itoa(e.HostStartTimestamp)),
		HostSlotOffset:                pulumi.String(strconv.Itoa(e.HostSlotOffset)),
		HostSlotDuration:              pulumi.String(strconv.Itoa(e.HostSlotDuration)),
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
