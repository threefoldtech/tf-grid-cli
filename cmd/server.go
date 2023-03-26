package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	command "github.com/threefoldtech/tf-grid-cli/internal/cmd"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run an rpc server listening for incoming commands to the tfgrid client",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := &daemon.Context{
			LogFilePerm: 0640,
			WorkDir:     "/",
			Umask:       027,
			// LogFileName: "/var/run/.tfgridclient.log",
		}

		d, err := ctx.Reborn()
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to run grid client server")
		}
		if d != nil {
			return
		}
		defer func() {
			_ = ctx.Release()
		}()

		err = command.RPCServer()
		if err != nil {
			log.Fatal().Err(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
