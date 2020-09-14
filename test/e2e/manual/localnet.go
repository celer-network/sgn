package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/transactor"
	sgnval "github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/viper"
)

var (
	start   = flag.Bool("start", false, "start local testnet")
	auto    = flag.Bool("auto", false, "auto-add all validators")
	down    = flag.Bool("down", false, "shutdown local testnet")
	up      = flag.Int("up", -1, "start a testnet node")
	stop    = flag.Int("stop", -1, "stop a testnet node")
	upall   = flag.Bool("upall", false, "start all nodes")
	stopall = flag.Bool("stopall", false, "stop all nodes")
	rebuild = flag.Bool("rebuild", false, "rebuild sgn node docker image")
)

func main() {
	flag.Parse()
	repoRoot, _ := filepath.Abs("../../..")
	if *start {
		tc.SetupMainchain()
		tc.SetupSidechain()
		p := &tc.SGNParams{
			CelrAddr:               tc.E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1000000000000000000),
			GovernVoteTimeout:      big.NewInt(30),
			SlashTimeout:           big.NewInt(15),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(5),
			MinStakingPool:         big.NewInt(5000000000000000000), // 5 CELR
			AdvanceNoticePeriod:    big.NewInt(30),
			SidechainGoLiveTimeout: big.NewInt(0),
		}
		tc.SetupNewSGNEnv(p, true)

		log.Infoln("install sgncli and sgnops in host machine")
		cmd := exec.Command("make", "install-tools")
		cmd.Dir = repoRoot
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "WITH_CLEVELDB=yes")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}

		log.Infoln("copy config files")
		cmd2 := exec.Command("make", "copy-manual-test-data")
		cmd2.Dir = repoRoot
		if err := cmd2.Run(); err != nil {
			log.Error(err)
		}
		log.Infoln("update config files")
		for i := 0; i < 3; i++ {
			configPath := fmt.Sprintf("./data/node%d/config.json", i)
			configFileViper := viper.New()
			configFileViper.SetConfigFile(configPath)
			if err := configFileViper.ReadInConfig(); err != nil {
				log.Error(err)
			}
			ksPath, _ := filepath.Abs(fmt.Sprintf("./data/node%d/keys/vethks%d.json", i, i))
			configFileViper.Set(common.FlagEthKeystore, ksPath)
			configFileViper.Set(common.FlagEthGateway, tc.LocalGeth)
			configFileViper.Set(common.FlagSgnNodeURI, tc.SgnNodeURIs[i])
			if err := configFileViper.WriteConfig(); err != nil {
				log.Error(err)
			}
		}
		if *auto {
			time.Sleep(10 * time.Second)
			addValidators()
		}
	} else if *down {
		log.Infoln("Tearing down all containers...")
		cmd := exec.Command("make", "localnet-down")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	} else if *up != -1 {
		log.Infoln("Start node", *up)
		cmd := exec.Command("docker-compose", "up", "-d", fmt.Sprintf("sgnnode%d", *up))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
	} else if *stop != -1 {
		log.Infoln("Stop node", *stop)
		cmd := exec.Command("docker-compose", "stop", fmt.Sprintf("sgnnode%d", *stop))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
	} else if *upall {
		log.Infoln("Start all nodes ...")
		cmd := exec.Command("make", "localnet-up-nodes")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	} else if *stopall {
		log.Infoln("Stop all nodes ...")
		cmd := exec.Command("make", "localnet-down-nodes")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	} else if *rebuild {
		log.Infoln("Rebuild sgn node docker image ...")
		cmd := exec.Command("make", "build-node")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	}
}

func addValidators() {
	cdc := app.MakeCodec()
	txr, err := transactor.NewTransactor(
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
		cdc, nil,
	)
	if err != nil {
		log.Error(err)
		return
	}

	amts := []*big.Int{
		new(big.Int).Mul(big.NewInt(10000), big.NewInt(common.TokenDec)),
		new(big.Int).Mul(big.NewInt(20000), big.NewInt(common.TokenDec)),
		new(big.Int).Mul(big.NewInt(10000), big.NewInt(common.TokenDec)),
	}
	minAmts := []*big.Int{
		new(big.Int).Mul(big.NewInt(1000), big.NewInt(common.TokenDec)),
		new(big.Int).Mul(big.NewInt(2000), big.NewInt(common.TokenDec)),
		new(big.Int).Mul(big.NewInt(3000), big.NewInt(common.TokenDec)),
	}
	commissions := []*big.Int{big.NewInt(150), big.NewInt(200), big.NewInt(120)}

	for i := 0; i < 3; i++ {
		log.Infoln("Adding validator", i)
		ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[i])
		if err != nil {
			log.Error(err)
			return
		}
		addCandidateWithStake(txr, ethAddr, auth, tc.ValAccounts[i], amts[i], minAmts[i], commissions[i], big.NewInt(300))
	}
}

func addCandidateWithStake(
	txr *transactor.Transactor,
	ethAddr mainchain.Addr,
	auth *bind.TransactOpts,
	valacct string,
	amt *big.Int, minAmt *big.Int,
	commissionRate *big.Int,
	rateLockEndTime *big.Int) {

	// get sgnAddr
	sgnAddr, err := sdk.AccAddressFromBech32(valacct)
	if err != nil {
		log.Error(err)
		return
	}

	// add candidate
	err = tc.InitializeCandidate(auth, sgnAddr, minAmt, commissionRate, rateLockEndTime)
	if err != nil {
		log.Error(err)
		return
	}
	for retry := 0; retry < 10; retry++ {
		c, err2 := sgnval.CLIQueryCandidate(txr.CliCtx, sgnval.RouterKey, ethAddr.Hex())
		if err2 == nil {
			log.Infof("query candidate success: %s", c)
			break
		}
		time.Sleep(time.Second)
	}

	// self delegate stake
	err = tc.DelegateStake(auth, ethAddr, amt)
	if err != nil {
		log.Error(err)
		return
	}
	for retry := 0; retry < 10; retry++ {
		_, err = sgnval.CLIQueryValidator(txr.CliCtx, staking.RouterKey, valacct)
		if err == nil {
			log.Infof("query validator success: %s", valacct)
			break
		}
		time.Sleep(time.Second)
	}
}
