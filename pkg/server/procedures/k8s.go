package procedure

import (
	"context"

	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

func K8sDeploy(ctx context.Context, cluster types.K8sCluster, client deployer.TFPluginClient) (types.K8sCluster, error) {
	// validate project name is unique
	// map to workloads.k8sCluster
	// use k8sDeployer to deploy cluster
	// construct result
	k8s := workloads.K8sCluster{}

	//master name must be: projectName.master
	//worker names must be: projectName.workerX
	client.K8sDeployer.Deploy(ctx)
}

func K8sDelete(ctx context.Context, clusterName string) error {
	// get all contracts by project name
	// client.SubstrateConn.CancelContract(client.Identity, contractID)
}

func K8sAddNode(ctx context.Context, cluserName string, node types.K8sNode) (types.K8sCluster, error)

func K8sRemoveNode(ctx context.Context, clusterName string, nodeName string) (types.K8sCluster, error)

func K8sGet(ctx context.Context, clusterName string, client deployer.TFPluginClient) (types.K8sCluster, error) {
	// get all contracts by project name
	masterMap := map[uint32]string{}
	workerMap := map[uint32][]string{}
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(clusterName)
	nodeClient, err := client.NcPool.GetNodeClient(client.SubstrateConn, contracts.NodeContracts[0].NodeID)
	dl, err := nodeClient.DeploymentGet(ctx, contractID)
	for _, w := range dl.Workloads {
		// if name indicates master, assign to master map
		// if name indicates worker, append to worker map
	}

	//
	client.State.LoadK8sFromGrid()
}
