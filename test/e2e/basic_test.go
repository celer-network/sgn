package e2e

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/celer-network/sgn/testing/log"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	cmd := exec.Command("sgncli", "query", "global", "latest-block")
	cmd.Dir, _ = filepath.Abs("../..")

	out, err := cmd.Output()
	assert.Equal(t, err, nil, "The command should run successfully")
	log.Infof("Latest block number is %s", out)
}
