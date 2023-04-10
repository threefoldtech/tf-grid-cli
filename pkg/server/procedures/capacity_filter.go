package procedure

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/types"
	"github.com/threefoldtech/tf-grid-cli/pkg/server/utils"
)

func FilterNodes(ctx context.Context, options types.FilterOptions, client *deployer.TFPluginClient) (types.FilterResult, error) {
	var res types.FilterResult
	var err error

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hasFarmerBot := utils.HasFarmerBot(ctx2, client, options.FarmID)

	if options.FarmID != 0 && hasFarmerBot {
		log.Info().Msg("Calling farmerbot")
		res, err = utils.FilterNodesWithFarmerBot(ctx, options, client)
	} else {
		log.Info().Msg("Calling gridproxy")
		res, err = utils.FilterNodesWithGridProxy(ctx, options, client)
	}

	return res, err
}
