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

// getVMCmd represents the get vm command
var getVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "Get deployed vm",
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

		vm, err := command.GetVM(t, args[0])
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		s, err := json.MarshalIndent(vm, "", "\t")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Msg("vm:\n" + string(s))

	},
}

func init() {
	getCmd.AddCommand(getVMCmd)
}
