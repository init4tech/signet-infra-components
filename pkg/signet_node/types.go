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
	HostZenithAddress             pulumi.StringInput `pulumi:"hostZenithAddress"`
	RuOrdersContractAddress       pulumi.StringInput `pulumi:"ruOrdersContractAddress"`
	HostOrdersContractAddress     pulumi.StringInput `pulumi:"hostOrdersContractAddress"`
	SignetChainId                 pulumi.StringInput `pulumi:"signetChainId"`
	BlobExplorerUrl               pulumi.StringInput `pulumi:"blobExplorerUrl"`
	SignetStaticPath              pulumi.StringInput `pulumi:"signetStaticPath"`
	SignetDatabasePath            pulumi.StringInput `pulumi:"signetDatabasePath"`
	RustLog                       pulumi.StringInput `pulumi:"rustLog"`
	IpcEndpoint                   pulumi.StringInput `pulumi:"ipcEndpoint"`
	RpcPort                       pulumi.StringInput `pulumi:"rpcPort"`
	WsRpcPort                     pulumi.StringInput `pulumi:"wsRpcPort"`
	TxForwardUrl                  pulumi.StringInput `pulumi:"txForwardUrl"`
	GenesisJsonPath               pulumi.StringInput `pulumi:"genesisJsonPath"`
	HostZenithDeployHeight        pulumi.StringInput `pulumi:"hostZenithDeployHeight"`
	BaseFeeRecipient              pulumi.StringInput `pulumi:"baseFeeRecipient"`
	HostPassageContractAddress    pulumi.StringInput `pulumi:"hostPassageContractAddress"`
	HostTransactorContractAddress pulumi.StringInput `pulumi:"hostTransactorContractAddress"`
	RuPassageContractAddress      pulumi.StringInput `pulumi:"ruPassageContractAddress"`
	SignetClUrl                   pulumi.StringInput `pulumi:"signetClUrl"`
	SignetPylonUrl                pulumi.StringInput `pulumi:"signetPylonUrl"`
	HostStartTimestamp            pulumi.StringInput `pulumi:"hostStartTimepstamp"`
	HostSlotOffset                pulumi.StringInput `pulumi:"hostSlotOffset"`
	HostSlotDuration              pulumi.StringInput `pulumi:"hostSlotDuration"`
}

// GetEnvMap implements the utils.EnvProvider interface
// It creates a map of environment variables from the SignetNodeEnv struct
func (e SignetNodeEnv) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}

// SignetNode interface defines methods that the SignetNodeComponent must implement
type SignetNode interface {
}

// Ensure SignetNodeComponent implements SignetNode
var _ SignetNode = &SignetNodeComponent{}
