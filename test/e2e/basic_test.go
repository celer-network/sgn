package e2e

import (
	"os/exec"
	"path/filepath"
	"testing"

	log "github.com/celer-network/goCeler/clog"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	cmd := exec.Command("sgncli", "query", "global", "latest-block")
	cmd.Dir, _ = filepath.Abs("../..")

	out, err := cmd.Output()
	assert.Equal(t, err, nil, "The command should run successfully")
	log.Infof("Latest block number is %s", out)
}
