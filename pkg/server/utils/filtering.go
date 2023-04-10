package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/grid3-go/deployer"
	proxyTypes "github.com/threefoldtech/grid_proxy_server/pkg/types"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
)

var (
	Status  = "up"
	TrueVal = true
)

const (
	FarmerBotVersionAction  = "farmerbot.farmmanager.version"
	FarmerBotFindNodeAction = "farmerbot.nodemanager.findnode"
	FarmerBotRMBFunction    = "execute_job"
)

type PlannedReservation struct {
	WorkloadName string
	NodeID       uint32
	FarmID       uint32
	MRU          uint64
	SRU          uint64
	HRU          uint64
	CRU          uint64
}

type Args struct {
	RequiredHRU  *uint64  `json:"required_hru,omitempty"`
	RequiredSRU  *uint64  `json:"required_sru,omitempty"`
	RequiredCRU  *uint64  `json:"required_cru,omitempty"`
	RequiredMRU  *uint64  `json:"required_mru,omitempty"`
	NodeExclude  []uint32 `json:"node_exclude,omitempty"`
	Dedicated    *bool    `json:"dedicated,omitempty"`
	PublicConfig *bool    `json:"public_config,omitempty"`
	PublicIPs    *uint32  `json:"public_ips"`
	Certified    *bool    `json:"certified,omitempty"`
}

type Params struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type FarmerBotArgs struct {
	Args   []Args   `json:"args"`
	Params []Params `json:"params"`
}

type FarmerBotAction struct {
	Guid         string        `json:"guid"`
	TwinID       uint32        `json:"twinid"`
	Action       string        `json:"action"`
	Args         FarmerBotArgs `json:"args"`
	Result       FarmerBotArgs `json:"result"`
	State        string        `json:"state"`
	Start        uint64        `json:"start"`
	End          uint64        `json:"end"`
	GracePeriod  uint32        `json:"grace_period"`
	Error        string        `json:"error"`
	Timeout      uint32        `json:"timeout"`
	SourceTwinID uint32        `json:"src_twinid"`
	SourceAction string        `json:"src_action"`
	Dependencies []string      `json:"dependencies"`
}

func BuildFarmerBotParams(options types.FilterOptions) []Params {
	params := []Params{}
	if options.HRU != 0 {
		params = append(params, Params{Key: "required_hru", Value: options.HRU})
	}

	if options.SRU != 0 {
		params = append(params, Params{Key: "required_sru", Value: options.SRU})
	}

	if options.MRU != 0 {
		params = append(params, Params{Key: "required_mru", Value: options.MRU})
	}

	if options.Dedicated {
		params = append(params, Params{Key: "dedicated", Value: options.Dedicated})
	}

	if options.PublicConfig {
		params = append(params, Params{Key: "public_config", Value: options.PublicConfig})
	}

	if options.PublicIpsCount > 0 {
		params = append(params, Params{Key: "public_ips", Value: options.PublicIpsCount})
	}

	return params
}

func BuildFarmerBotAction(farmerTwinID uint32, sourceTwinID uint32, args []Args, params []Params, action string) FarmerBotAction {
	return FarmerBotAction{
		Guid:   uuid.NewString(),
		TwinID: farmerTwinID,
		Action: action,
		Args: FarmerBotArgs{
			Args:   args,
			Params: params,
		},
		Result: FarmerBotArgs{
			Args:   []Args{},
			Params: []Params{},
		},
		State:        "init",
		Start:        uint64(time.Now().Unix()),
		End:          0,
		GracePeriod:  0,
		Error:        "",
		Timeout:      6000,
		SourceTwinID: sourceTwinID,
		Dependencies: []string{},
	}
}

func GetFarmerTwinIDByFarmID(client *deployer.TFPluginClient, farmID uint32) (uint32, error) {
	farmid := uint64(farmID)
	farms, _, err := client.GridProxyClient.Farms(proxyTypes.FarmFilter{
		FarmID: &farmid,
	}, proxyTypes.Limit{
		Size: 1,
		Page: 1,
	})

	if err != nil || len(farms) == 0 {
		return 0, errors.Wrapf(err, "Couldn't get the FarmerTwinID for FarmID: %+v", farmID)
	}

	return uint32(farms[0].TwinID), nil
}

func GetFarmerBotResult(action FarmerBotAction, key string) (string, error) {
	if len(action.Result.Params) > 0 {
		for _, param := range action.Result.Params {
			if param.Key == key {
				return fmt.Sprint(param.Value), nil
			}
		}

	}

	return "", fmt.Errorf("Couldn't found a result for the same key: %s", key)
}

func FilterNodesWithFarmerBot(ctx context.Context, options types.FilterOptions, client *deployer.TFPluginClient) (types.FilterResult, error) {

	// construct farmerbot request
	params := BuildFarmerBotParams(options)

	// make farmerbot request
	farmerTwinID, err := GetFarmerTwinIDByFarmID(client, options.FarmID)
	if err != nil {
		return types.FilterResult{}, errors.Wrapf(err, "Failed to get TwinID for FarmID %+v", options.FarmID)
	}

	// TODO: Fix this by upgrade go-client
	sourceTwinID := uint32(220)

	data := BuildFarmerBotAction(farmerTwinID, sourceTwinID, []Args{}, params, FarmerBotFindNodeAction)

	var output FarmerBotAction

	err = client.RMB.Call(ctx, farmerTwinID, FarmerBotRMBFunction, data, &output)
	if err != nil {
		return types.FilterResult{}, errors.Wrapf(err, "Failed calling farmerbot on farm %d", options.FarmID)
	}

	// build the result
	nodeIdStr, err := GetFarmerBotResult(output, "nodeid")
	if err != nil {
		return types.FilterResult{}, err
	}

	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 32)
	if err != nil {
		return types.FilterResult{}, fmt.Errorf("can't parse node id")
	}

	result := types.FilterResult{
		FilterOption:   options,
		AvailableNodes: []uint32{uint32(nodeId)},
	}

	return result, nil
}

func FilterNodesWithGridProxy(ctx context.Context, options types.FilterOptions, client *deployer.TFPluginClient) (types.FilterResult, error) {
	proxyFilters := proxyTypes.NodeFilter{
		Status:  &Status,
		FreeMRU: &options.MRU,
		FreeSRU: &options.SRU,
		FreeHRU: &options.HRU,
		// TODO: add the others filters
	}

	nodes, err := deployer.FilterNodes(client.GridProxyClient, proxyFilters)
	if err != nil || len(nodes) == 0 {
		return types.FilterResult{}, errors.Wrapf(err, "Couldn't find node for the provided filters: %+v", options)
	}

	nodesIDs := GetNodesIDs(nodes)

	result := types.FilterResult{
		FilterOption:   options,
		AvailableNodes: nodesIDs,
	}

	return result, nil
}

func GetNodesIDs(nodes []proxyTypes.Node) []uint32 {
	ids := []uint32{}

	for _, node := range nodes {
		ids = append(ids, uint32(node.NodeID))
	}

	return ids
}

func HasFarmerBot(ctx context.Context, client *deployer.TFPluginClient, farmID uint32) bool {
	args := []Args{}
	params := []Params{}

	farmerTwinID, err := GetFarmerTwinIDByFarmID(client, farmID)

	// TODO: Fix this by upgrade go-client
	sourceTwinID := uint32(220)

	data := BuildFarmerBotAction(farmerTwinID, sourceTwinID, args, params, FarmerBotVersionAction)

	var output FarmerBotAction

	err = client.RMB.Call(ctx, farmerTwinID, FarmerBotRMBFunction, data, &output)

	return err == nil
}

func checkNodeAvailability(client *deployer.TFPluginClient, nodeId uint32, workload PlannedReservation, reservedCapacity map[uint32]PlannedReservation) bool {
	// get node info
	node, err := client.GridProxyClient.Node(nodeId)
	if err != nil {
		return false
	}

	// get free resources
	free := proxyTypes.Capacity{
		CRU: node.Capacity.Total.CRU - node.Capacity.Used.CRU,
		SRU: node.Capacity.Total.SRU - node.Capacity.Used.SRU,
		HRU: node.Capacity.Total.HRU - node.Capacity.Used.HRU,
		MRU: node.Capacity.Total.MRU - node.Capacity.Used.MRU,
	}

	// check if the free resource greater than the previous reserved capacity plus the current workload capacity
	if uint64(free.MRU) >= reservedCapacity[uint32(node.NodeID)].MRU+workload.MRU &&
		uint64(free.SRU) >= reservedCapacity[uint32(node.NodeID)].SRU+workload.SRU {
		return true
	}
	return false
}

// Searching for node for each workload considering the reserved capacity by workloads in the same deployment.
// Assign the NodeID if found one or return it with NodeID: 0
func AssignNodes(ctx context.Context, client *deployer.TFPluginClient, workloads []*PlannedReservation) error {
	reservedCapacity := make(map[uint32]PlannedReservation)

	for _, workload := range workloads {
		options := types.FilterOptions{
			FarmID: workload.FarmID,
			// HRU: workload.NeededStorage,
			SRU: workload.SRU,
			MRU: workload.MRU,
		}

		var res types.FilterResult
		var err error
		ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		// TODO: store the result to reduce the number of calls
		hasFarmerBot := HasFarmerBot(ctx2, client, options.FarmID)

		if options.FarmID != 0 && hasFarmerBot {
			log.Info().Msg("Calling farmerbot")
			res, err = FilterNodesWithFarmerBot(ctx, options, client)

		} else {
			log.Info().Msg("Calling gridproxy")
			res, err = FilterNodesWithGridProxy(ctx, options, client)
		}

		if err != nil || len(res.AvailableNodes) == 0 {
			return errors.Errorf("Failed to find node on farm %+v", options.FarmID)
		}

		nodes := res.AvailableNodes

		selectedNodeId := uint32(0)

		for _, nodeId := range nodes {
			valid := checkNodeAvailability(client, nodeId, *workload, reservedCapacity)
			if valid {
				selectedNodeId = nodeId
				break
			}
		}

		reservedCapacity[selectedNodeId] = PlannedReservation{
			MRU: reservedCapacity[selectedNodeId].MRU + workload.MRU,
			SRU: reservedCapacity[selectedNodeId].SRU + workload.SRU,
		}

		workload.NodeID = selectedNodeId
	}

	return nil
}
