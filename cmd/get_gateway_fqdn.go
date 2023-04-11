// Package cmd for parsing command line arguments
package cmd

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/tf-grid-cli/internal/config"
)

// getGatewayFQDNCmd represents the get gateway fqdn command
var getGatewayFQDNCmd = &cobra.Command{
	Use:   "fqdn",
	Short: "Get deployed gateway FQDN",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetUserConfig()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		t, err := deployer.NewTFPluginClient(cfg.Mnemonics, "sr25519", cfg.Network, "", "", "", 100, true, false)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		contracts, err := t.ContractsGetter.ListContractsOfProjectName(args[0])
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		var nodeID uint32
		var contractID uint64
		for _, contract := range contracts.NodeContracts {
			var deploymentData workloads.DeploymentData
			err := json.Unmarshal([]byte(contract.DeploymentData), &deploymentData)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			if deploymentData.Type != "Gateway Fqdn" || deploymentData.Name != args[0] {
				continue
			}
			nodeID = contract.NodeID
			contractID, err = strconv.ParseUint(contract.ContractID, 0, 64)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			t.State.CurrentNodeDeployments[nodeID] = []uint64{contractID}
			break
		}
		if nodeID == 0 {
			log.Info().Msgf("no gateway fqdn with name %s found", args[0])
			os.Exit(0)
		}
		gateway, err := t.State.LoadGatewayFQDNFromGrid(nodeID, args[0], args[0])
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		s, err := json.MarshalIndent(gateway, "", "\t")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Msg("gateway fqdn:\n" + string(s))
	},
}

func init() {
	getGatewayCmd.AddCommand(getGatewayFQDNCmd)
}
