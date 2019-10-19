package e2e

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	conn, err := ethclient.Dial(tf.EthInstance)
	if err != nil {
		os.Exit(1)
	}

	route := fmt.Sprintf("custom/%s/%s", global.ModuleName, global.QueryLatestBlock)
	resJson, _, err := tf.Transactor.CliCtx.Query(route)
	var result map[string]interface{}
	json.Unmarshal(resJson, &result)
	blkNumSGN := new(big.Int)
	blkNumSGN, _ = blkNumSGN.SetString(result["number"].(string), 10)
	log.Infof("Latest block number on SGN is %d", blkNumSGN)

	blkNumMain, err := tf.GetLatestBlkNum(conn)
	if err != nil {
		os.Exit(1)
	}
	log.Infof("Latest block number on mainchain is %d", blkNumMain)

	assert.Equal(t, err, nil, "The command should run successfully")
	diff := new(big.Int).Sub(blkNumMain, blkNumSGN)
	assert.GreaterOrEqual(t, big.NewInt(5).Cmp(diff), 0, "blkNumMain should be greater than or equal to blkNumSGN")
}
