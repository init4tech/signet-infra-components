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

type SignetNodeComponentArgs struct {
	Name                        string
	Namespace                   pulumi.StringInput
	Env                         SignetNodeEnv
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

type SignetNodeEnv struct {
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

// GetEnvMap implements the utils.EnvProvider interface
// It creates a map of environment variables from the SignetNodeEnv struct
func (e SignetNodeEnv) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// ConsensusEnv contains environment variables for the consensus client
type ConsensusEnv struct {
	Example pulumi.StringInput `pulumi:"example"`
}

// GetEnvMap implements the utils.EnvProvider interface
func (e ConsensusEnv) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// SignetNode interface defines methods that the SignetNodeComponent must implement
type SignetNode interface {
}

// Ensure SignetNodeComponent implements SignetNode
var _ SignetNode = &SignetNodeComponent{}
