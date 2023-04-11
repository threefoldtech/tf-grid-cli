// Package cmd for parsing command line arguments
package cmd

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/grid3-go/deployer"
	command "github.com/threefoldtech/tf-grid-cli/internal/cmd"
	"github.com/threefoldtech/tf-grid-cli/internal/config"
)

// getKubernetesCmd represents the get kubernetes command
var getKubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Get deployed kubernetes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetUserConfig()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		t, err := deployer.NewTFPluginClient(cfg.Mnemonics, "sr25519", cfg.Network, "", "", "", 100, true, false)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		cluster, err := command.GetK8sCluster(t, args[0])
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		s, err := json.MarshalIndent(cluster, "", "\t")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Msg("k8s cluster:\n" + string(s))

	},
}

func init() {
	getCmd.AddCommand(getKubernetesCmd)
}
