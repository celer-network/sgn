// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/spf13/viper"
)

func shutdownNode(node uint) {
	log.Infoln("Shutdown node", node)
	cmd := exec.Command("docker-compose", "stop", fmt.Sprintf("sgnnode%d", node))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

func restartNode(node uint) {
	log.Infoln("Shutdown node", node)
	cmd := exec.Command("docker-compose", "restart", fmt.Sprintf("sgnnode%d", node))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

func restartWithConfig(node uint, flag string, value interface{}) {
	log.Infof("restart node %d with config %s = %v", node, flag, value)

	configPath := fmt.Sprintf("../../../docker-volumes/node%d/sgncli/config/sgn.toml", node)
	configFileViper := viper.New()
	configFileViper.SetConfigFile(configPath)
	err := configFileViper.ReadInConfig()
	tc.ChkErr(err, "Failed to read config")
	configFileViper.Set(flag, value)
	err = configFileViper.WriteConfig()
	tc.ChkErr(err, "Failed to write config")

	restartNode(node)
}
