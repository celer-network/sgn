// Copyright 2018 Celer Network

package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/celer-network/goCeler/ctype"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	s2Keystore = "../../testing/env/server2.json"
)

var (
	sStoreDir = path.Join(sStoreDirPrefix, sEthAddr)
)

// used by setup_onchain and tests
var (
	etherBaseAddr = ctype.Hex2Addr(etherBaseAddrStr)
	clientAddr    = ctype.Hex2Addr(clientAddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir     string
	envDir         = "../../testing/env"
	noProxyProfile string // full file path to profile.json
	// erc20 token addr hex
	// map from app type to deployed addr, updated by SetupMainchain
	appAddrMap     = make(map[string]ctype.Addr)
	tokenAddrErc20 string // set by SetupMainchain deploy erc20 contract
)

// toBuild map package subpath to binary file name eg. cmd/sgn -> sgn means build sgn/cmd/sgn and output sgn
var toBuild = map[string]string{
	"cmd/sgn":    "sgn",
	"cmd/sgncli": "sgncli",
}

// start process to handle eth rpc, and fund etherbase and server account
func StartMainchain() (*os.Process, error) {
	log.Println("outRootDir", outRootDir, "envDir", envDir)
	chainDataDir := outRootDir + "chaindata"
	logFname := outRootDir + "chain.log"
	if err := os.MkdirAll(chainDataDir, os.ModePerm); err != nil {
		return nil, err
	}

	// geth init
	cmdInit := exec.Command("geth", "--datadir", chainDataDir, "init", envDir+"/mainchain_genesis.json")
	// set cmd.Dir because relative files are under testing/env
	cmdInit.Dir, _ = filepath.Abs(envDir)
	if err := cmdInit.Run(); err != nil {
		return nil, err
	}
	// actually run geth, blocking. set syncmode full to avoid bloom mem cache by fast sync
	cmd := exec.Command("geth", "--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", "--maxpeers", "0",
		"--netrestrict", "127.0.0.1/8", "--datadir", chainDataDir, "--keystore", "keystore", "--targetgaslimit", "8000000",
		"--mine", "--allow-insecure-unlock", "--unlock", "0", "--password", "empty_password.txt", "--rpc", "--rpccorsdomain", "*",
		"--rpcapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3")
	cmd.Dir = cmdInit.Dir

	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	fmt.Println("geth pid:", cmd.Process.Pid)
	// in case geth exits with non-zero, exit test early
	// if geth is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("geth process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, nil
}

// todo: remove addr arg
func getEthClient(addr string) (*ethclient.Client, error) {
	ws, err := ethrpc.Dial(ethGateway)
	if err != nil {
		return nil, err
	}
	conn := ethclient.NewClient(ws)
	return conn, nil
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func StartSidechainDefault(rootDir string) (*os.Process, *exec.Cmd, error) {
	cmd := exec.Command("make", "copy-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	removeCmd := exec.Command("rm", "-rf", "~/.sgn", "~/.sgncli")

	cmd = exec.Command("sgn", "start")
	cmd.Dir, _ = filepath.Abs("../..")
	logFname := rootDir + "sgn.log"
	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	if err := cmd.Start(); err != nil {
		return nil, removeCmd, err
	}

	fmt.Println("sgn pid:", cmd.Process.Pid)
	// in case sgn exits with non-zero, exit test early
	// if sgn is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("sgn process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, removeCmd, nil
}

// save json as file path
// func saveProfile(p *common.CProfile, fpath string) {
// 	b, _ := json.Marshal(p)
// 	ioutil.WriteFile(fpath, b, 0644)
// }

// func SaveProfile(p *common.CProfile, fpath string) {
// 	saveProfile(p, fpath)
// }

// func buildBins(rootDir string) error {
// 	sgnRepo := "github.com/celer-network/sgn/"
// 	for pkg, bin := range toBuild {
// 		fmt.Println("Building", pkg, "->", bin)
// 		cmd := exec.Command("go", "build", "-o", rootDir+bin, sgnRepo+pkg)
// 		cmd.Stderr, _ = os.OpenFile(rootDir+"build.err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 		err := cmd.Run()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func installBins() error {
	cmd := exec.Command("make", "install")
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func chkErr(e error, msg string) {
	if e != nil {
		fmt.Println("Err:", msg, e)
		os.Exit(1)
	}
}

func CheckError(e error, msg string) {
	chkErr(e, msg)
}

func SetEnvDir(dir string) {
	envDir = dir
}

func SetOutRootDir(dir string) {
	outRootDir = dir
}
