package transactor

import (
	"fmt"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/viper"
)

const (
	blockWindowSize    = 100
	incrementThreshold = 100
)

type GasPriceEstimator struct {
	currentHeight int64
	txCounts      []int
	txCountSum    int
	baseGasPrice  float64
	cliCtx        context.CLIContext
}

func NewGasPriceEstimator(nodeURI string) *GasPriceEstimator {
	cliCtx := context.
		NewCLIContext().
		WithNodeURI(nodeURI)

	gpe := &GasPriceEstimator{
		txCounts:     []int{},
		txCountSum:   0,
		baseGasPrice: viper.GetFloat64(common.FlagSgnBaseGasPrice),
		cliCtx:       cliCtx,
	}

	go gpe.collectStatistics()
	return gpe
}

func (gpe *GasPriceEstimator) collectStatistics() {
	ticker := time.NewTicker(viper.GetDuration(common.FlagSgnTimeoutCommit) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		node, err := gpe.cliCtx.GetNode()
		if err != nil {
			log.Errorln("GetNode error:", err)
			continue
		}

		block, err := node.Block(nil)
		if err != nil {
			log.Errorln("Query block error:", err)
			continue
		}
		if block.Block == nil {
			continue
		}
		if gpe.currentHeight == block.Block.Height {
			continue
		}

		gpe.currentHeight = block.Block.Height
		txCount := len(block.Block.Txs)
		gpe.txCounts = append(gpe.txCounts, txCount)
		gpe.txCountSum += txCount

		if len(gpe.txCounts) > blockWindowSize {
			gpe.txCountSum -= gpe.txCounts[0]
			gpe.txCounts = gpe.txCounts[1:]
		}
	}
}

// Get proper gas price to send a tx
func (gpe *GasPriceEstimator) GetGasPrice() string {
	if gpe.baseGasPrice == 0 {
		return ""
	}

	var priceMultiplier float64 = 0
	if len(gpe.txCounts) != 0 {
		priceMultiplier = float64(gpe.txCountSum / len(gpe.txCounts) / incrementThreshold)
	}
	if priceMultiplier == 0 {
		priceMultiplier = 1
	}

	gasPrice := priceMultiplier * gpe.baseGasPrice
	log.Info(gasPrice)
	return fmt.Sprintf("%f%s", gasPrice, common.QuotaCoinName)
}
