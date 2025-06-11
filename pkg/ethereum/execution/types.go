package execution

import (
	"github.com/init4tech/signet-infra-components/pkg/utils"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ExecutionClientArgs contains the configuration for an execution client
type ExecutionClientArgs struct {
	// Name is the base name for all resources
	Name pulumi.StringInput `pulumi:"name"`
	// Namespace is the Kubernetes namespace to deploy resources in
	Namespace pulumi.StringInput `pulumi:"namespace"`
	// StorageSize is the size of the persistent volume claim
	StorageSize pulumi.StringInput `pulumi:"storageSize"`
	// StorageClass is the Kubernetes storage class to use
	StorageClass pulumi.StringInput `pulumi:"storageClass"`
	// Image is the container image to use
	Image pulumi.StringInput `pulumi:"image"`
	// ImagePullPolicy is the Kubernetes image pull policy
	ImagePullPolicy pulumi.StringInput `pulumi:"imagePullPolicy"`
	// Resources contains the resource requests and limits
	Resources *corev1.ResourceRequirements `pulumi:"resources,optional"`
	// NodeSelector is the Kubernetes node selector
	NodeSelector pulumi.StringMap `pulumi:"nodeSelector,optional"`
	// Tolerations are the Kubernetes tolerations
	Tolerations corev1.TolerationArray `pulumi:"tolerations,optional"`
	// JWTSecret is the JWT secret for authentication
	JWTSecret pulumi.StringInput `pulumi:"jwtSecret"`
	// P2PPort is the port for P2P communication
	P2PPort pulumi.IntInput `pulumi:"p2pPort"`
	// RPCPort is the port for RPC communication
	RPCPort pulumi.IntInput `pulumi:"rpcPort"`
	// WSPort is the port for WebSocket communication
	WSPort pulumi.IntInput `pulumi:"wsPort"`
	// MetricsPort is the port for metrics
	MetricsPort pulumi.IntInput `pulumi:"metricsPort"`
	// AuthRPCPort is the port for authenticated RPC
	AuthRPCPort pulumi.IntInput `pulumi:"authRpcPort"`
	// DiscoveryPort is the port for node discovery
	DiscoveryPort pulumi.IntInput `pulumi:"discoveryPort"`
	// Bootnodes is a list of bootnode URLs
	Bootnodes pulumi.StringArray `pulumi:"bootnodes,optional"`
	// AdditionalArgs are additional command line arguments
	AdditionalArgs pulumi.StringArray `pulumi:"additionalArgs,optional"`
	// Environment variables
	ExecutionClientEnv ExecutionClientEnv `pulumi:"executionClientEnv,optional"`
}

// ExecutionClientComponent represents an execution client deployment
type ExecutionClientComponent struct {
	pulumi.ResourceState

	// Name is the base name for all resources
	Name string
	// Namespace is the Kubernetes namespace
	Namespace string
	// ConfigMap is the shared config map
	ConfigMap *corev1.ConfigMap
	// PVC is the persistent volume claim
	PVC *corev1.PersistentVolumeClaim
	// JWTSecret is the JWT secret
	JWTSecret *corev1.Secret
	// P2PService is the P2P service
	P2PService *corev1.Service
	// RPCService is the RPC service
	RPCService *corev1.Service
	// StatefulSet is the stateful set
	StatefulSet *appsv1.StatefulSet
}

// ExecutionClientEnv contains environment variables for the execution client
type ExecutionClientEnv struct {
	// HOST_ZENITH_CONTRACT_ADDRESS - The address of the Host Zenith contract
	HostZenithContractAddress pulumi.StringInput `pulumi:"hostZenithContractAddress"`
	// RU_ORDERS_CONTRACT_ADDRESS - The address of the Rollup Orders contract
	RuOrdersContractAddress pulumi.StringInput `pulumi:"ruOrdersContractAddress"`
	// HOST_ORDERS_CONTRACT_ADDRESS - The address of the Host Orders contract
	HostOrdersContractAddress pulumi.StringInput `pulumi:"hostOrdersContractAddress"`
	// SIGNET_CHAIN_ID - The chain ID of the Signet network
	SignetChainID pulumi.StringInput `pulumi:"signetChainID"`
	// BLOB_EXPLORER_URL - The URL of the Blob Explorer
	BlobExplorerUrl pulumi.StringInput `pulumi:"blobExplorerUrl"`
	// SIGNET_STATIC_PATH - The path to the Signet static files
	SignetStaticPath pulumi.StringInput `pulumi:"signetStaticPath"`
	// SIGNET_DATABASE_PATH - The path to the Signet database
	SignetDatabasePath pulumi.StringInput `pulumi:"signetDatabasePath"`
	// RUST_LOG - The log level for the signet node
	RustLog pulumi.StringInput `pulumi:"rustLog"`
	// IPC_ENDPOINT - The IPC endpoint for the Signet client
	IpcEndpoint pulumi.StringInput `pulumi:"ipcEndpoint"`
	// RPC_PORT - The port for the JSON RPC service
	RpcPort pulumi.StringInput `pulumi:"rpcPort"`
	// WS_RPC_PORT - The port for the WebSocket RPC service
	WsRpcPort pulumi.StringInput `pulumi:"wsRpcPort"`
	// TX_FORWARD_URL - The URL for the transaction forwarder to send transactions to
	TxForwardUrl pulumi.StringInput `pulumi:"txForwardUrl"`
	// GENESIS_JSON_PATH - The path to the genesis JSON file
	GenesisJsonPath pulumi.StringInput `pulumi:"genesisJsonPath"`
	// HOST_ZENITH_DEPLOY_HEIGHT - The height of the Host Zenith contract deployment
	HostZenithDeployHeight pulumi.StringInput `pulumi:"hostZenithDeployHeight"`
	// BASE_FEE_RECIPIENT - The address of the base fee recipient
	BaseFeeRecipient pulumi.StringInput `pulumi:"baseFeeRecipient"`
	// HOST_PASSAGE_CONTRACT_ADDRESS - The address of the Host Passage contract
	HostPassageContractAddress pulumi.StringInput `pulumi:"hostPassageContractAddress"`
	// HOST_TRANSACTOR_CONTRACT_ADDRESS - The address of the Host Transactor contract
	HostTransactorContractAddress pulumi.StringInput `pulumi:"hostTransactorContractAddress"`
	// RU_PASSAGE_CONTRACT_ADDRESS - The address of the Rollup Passage contract
	RuPassageContractAddress pulumi.StringInput `pulumi:"ruPassageContractAddress"`
	// SIGNET_CL_URL - The URL of the consensus client http api to fetch blobs
	SignetClUrl pulumi.StringInput `pulumi:"signetClUrl"`
	// SIGNET_PYLON_URL - The URL of the pylon api to fetch blobs
	SignetPylonUrl pulumi.StringInput `pulumi:"signetPylonUrl"`
	// START_TIMESTAMP - The start timestamp for the signet node in epoch time
	StartTimestamp pulumi.StringInput `pulumi:"startTimestamp"`
	// SLOT_OFFSET - The slot offset for signet node in seconds
	SlotOffset pulumi.StringInput `pulumi:"slotOffset"`
	// SLOT_DURATION - The slot duration for the signet node in seconds
	SlotDuration pulumi.StringInput `pulumi:"slotDuration"`
}

// GetEnvMap implements the utils.EnvProvider interface
func (e ExecutionClientEnv) GetEnvMap() pulumi.StringMap {
	return utils.CreateEnvMap(e)
}
