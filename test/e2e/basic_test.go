package e2e

import (
	"fmt"
	"testing"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	route := fmt.Sprintf("custom/%s/%s", global.ModuleName, global.QueryLatestBlock)
	block, _, err := tf.Transactor.CliCtx.Query(route)
	assert.Equal(t, err, nil, "The command should run successfully")
	log.Infof("Latest block number is %s", block)
}
