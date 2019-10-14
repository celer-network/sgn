// Copyright 2018 Celer Network

package testing

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	ethInstance = "http://127.0.0.1:8545"
)

// var (
// 	// TODO: we need to have clear port space/allocation
// 	clientPort       = 20086
// 	clientPortLock   sync.Mutex
// 	pendingNonceLock sync.Mutex
// 	// e2eProfile after deploy contracts, with hardcoded ethgw etc
// 	// serialized json will be saved as outRootDir/profile.json
// 	// tests wish to use a different profile can overrides fields like svrRpc
// 	// and keep contract addresses etc
// 	E2eProfile *ccommon.CProfile
// )

var (
	// to be set by test/e2e
	outRootDir  string
	envDir      = "../../testing/env"
	etherBaseKs = envDir + "/keystore/etherbase.json"
)

func SetEnvDir(dir string) {
	envDir = dir
	etherBaseKs = envDir + "/keystore/etherbase.json"
}

func SetOutRootDir(dir string) {
	outRootDir = dir
}

func StartProcess(name string, args ...string) *os.Process {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return cmd.Process
}

func KillProcess(process *os.Process) {
	process.Kill()
	process.Release()
}

// UpdateShadowStorage updates the shadow directory for database storage
// func UpdateShadowStorage(storeDir, shadowDir string) {
// 	const lock = "LOCK"
// 	if _, err := os.Stat(shadowDir); os.IsNotExist(err) {
// 		os.Mkdir(shadowDir, 0755)
// 		os.Create(path.Join(shadowDir, lock))
// 	}
// 	files, _ := ioutil.ReadDir(storeDir)
// 	for _, f := range files {
// 		name := f.Name()
// 		if name != lock {
// 			os.Symlink(path.Join(storeDir, name), path.Join(shadowDir, name))
// 		}
// 	}
// }

// func prepareEthClient() (
// 	*ethclient.Client, *bind.TransactOpts, context.Context, common.Address, error) {
// 	conn, err := ethclient.Dial(ethInstance)
// 	if err != nil {
// 		return nil, nil, nil, common.Address{}, err
// 	}
// 	fmt.Println("etherBaseKs", etherBaseKs)
// 	etherBaseKsBytes, err := ioutil.ReadFile(etherBaseKs)
// 	if err != nil {
// 		return nil, nil, nil, common.Address{}, err
// 	}
// 	etherBaseAddrStr, err := utils.GetAddressFromKeystore(etherBaseKsBytes)
// 	if err != nil {
// 		return nil, nil, nil, common.Address{}, err
// 	}
// 	etherBaseAddr := ctype.Hex2Addr(etherBaseAddrStr)
// 	auth, err := bind.NewTransactor(strings.NewReader(string(etherBaseKsBytes)), "")
// 	if err != nil {
// 		return nil, nil, nil, common.Address{}, err
// 	}
// 	return conn, auth, context.Background(), etherBaseAddr, nil
// }

// func fundAccount(amount string, recipients []*common.Address) error {
// 	conn, auth, ctx, senderAddr, err := prepareEthClient()
// 	if err != nil {
// 		return err
// 	}
// 	value := big.NewInt(0)
// 	value.SetString(amount, 10)
// 	auth.Value = value
// 	chainID := big.NewInt(883) // Celer Private Testnet
// 	var gasLimit uint64 = 21000
// 	for _, r := range recipients {
// 		pendingNonceLock.Lock()
// 		nonce, err := conn.PendingNonceAt(ctx, senderAddr)
// 		if err != nil {
// 			pendingNonceLock.Unlock()
// 			return err
// 		}
// 		gasPrice, err := conn.SuggestGasPrice(ctx)
// 		if err != nil {
// 			pendingNonceLock.Unlock()
// 			return err
// 		}
// 		tx := types.NewTransaction(nonce, *r, auth.Value, gasLimit, gasPrice, nil)
// 		tx, err = auth.Signer(types.NewEIP155Signer(chainID), senderAddr, tx)
// 		if err != nil {
// 			pendingNonceLock.Unlock()
// 			return err
// 		}
// 		if *r == ctype.ZeroAddr {
// 			log.Info("Advancing block")
// 		} else {
// 			log.Infof(
// 				"Sending %s wei from %s to %s, nonce %d. tx: %x", amount, senderAddr.Hex(), r.Hex(), nonce, tx.Hash())
// 		}

// 		err = conn.SendTransaction(ctx, tx)
// 		if err != nil {
// 			pendingNonceLock.Unlock()
// 			return err
// 		}
// 		pendingNonceLock.Unlock()
// 		receipt, err := utils.WaitMined(ctx, conn, tx, 0)
// 		if err != nil {
// 			log.Error(err)
// 		}
// 		if receipt.Status != 1 {
// 			log.Errorln("tx failed. tx hash:", receipt.TxHash.Hex())
// 		} else {
// 			if *r == ctype.ZeroAddr {
// 				head, _ := conn.HeaderByNumber(ctx, nil)
// 				log.Info("Current block number:", head.Number.String())
// 			} else {
// 				bal, _ := conn.BalanceAt(ctx, *r, nil)
// 				log.Infoln("tx done.", r.String(), "bal:", bal.String())
// 			}
// 		}
// 	}
// 	return nil
// }

// // CreateAccountsWithBalance creates Ethereum accounts with balances.
// func CreateAccountsWithBalance(num int, amount string) ([]string, []string, error) {
// 	paths := make([]string, num)
// 	addrs := make([]*common.Address, num)
// 	var ksDir string
// 	if outRootDir == "" { // in case somehow outRootDir isn't set
// 		ksDir = os.TempDir()
// 	} else {
// 		ksDir = outRootDir + "ksdir"
// 		os.MkdirAll(ksDir, os.ModePerm)
// 	}

// 	for i := 0; i < num; i++ {
// 		ks := keystore.NewKeyStore(ksDir, keystore.LightScryptN, keystore.LightScryptP)
// 		account, err := ks.NewAccount("")
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		addrHex := ctype.Addr2Hex(account.Address)
// 		// geth ks file eg. UTC--2018-09-27T06-20-16.449720124Z--ba756d65a1a03f07d205749f35e2406e4a8522ad
// 		fnamePattern := "UTC*" + addrHex
// 		files, err := filepath.Glob(filepath.Join(ksDir, fnamePattern))
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		if len(files) != 1 {
// 			return nil, nil, errors.New("incorrect numnber of glob keystore files")
// 		}
// 		paths[i] = files[0]
// 		addrs[i] = &account.Address
// 	}
// 	err := fundAccount(amount, addrs)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	addrStrings := make([]string, num)
// 	for i, addr := range addrs {
// 		addrStrings[i] = addr.String()
// 	}
// 	return paths, addrStrings, nil
// }

// func getAuthFor(ksfile string) (*bind.TransactOpts, error) {
// 	ksBytes, err := ioutil.ReadFile(ksfile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	key, err := keystore.DecryptKey(ksBytes, "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Infoln(ksfile, ctype.Bytes2Hex(crypto.FromECDSA(key.PrivateKey)))
// 	ksStr := string(ksBytes)
// 	auth, err := bind.NewTransactor(strings.NewReader(ksStr), "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	price := big.NewInt(0)
// 	price.SetString("2000000000", 10)
// 	auth.GasPrice = price
// 	return auth, nil
// }

// func FundAddr(amt string, recipients []*common.Address) error {
// 	return fundAccount(amt, recipients)
// }

// func FundAccountsWithErc20(erc20Addr string, addrs []string, amount string) error {
// 	conn, auth, ctx, _, err := prepareEthClient()
// 	if err != nil {
// 		return err
// 	}
// 	erc20Contract, err := chain.NewERC20(common.HexToAddress(erc20Addr), conn)
// 	if err != nil {
// 		return err
// 	}
// 	tokenAmt := new(big.Int)
// 	tokenAmt.SetString(amount, 10)
// 	for _, addr := range addrs {
// 		pendingNonceLock.Lock()
// 		tx, err := erc20Contract.Transfer(auth, common.HexToAddress(addr), tokenAmt)
// 		pendingNonceLock.Unlock()
// 		if err != nil {
// 			return err
// 		}
// 		utils.WaitMined(ctx, conn, tx, 0)
// 	}
// 	return nil
// }

// func GetNextClientPort() string {
// 	clientPortLock.Lock()
// 	ret := clientPort
// 	clientPort++
// 	clientPortLock.Unlock()
// 	return strconv.Itoa(ret)
// }

// func CalcNewBalance(balanceBefore string, delta string) string {
// 	balanceAfter := big.NewInt(0)
// 	balanceAfter.SetString(balanceBefore, 10)
// 	deltaInt := big.NewInt(0)
// 	deltaInt.SetString(delta, 10)
// 	return balanceAfter.Add(balanceAfter, deltaInt).String()
// }

// func GetAddressFromKeystore(ksBytes []byte) (common.Address, error) {
// 	type ksStruct struct {
// 		Address string
// 	}
// 	var ks ksStruct
// 	if err := json.Unmarshal(ksBytes, &ks); err != nil {
// 		return common.Address{}, err
// 	}
// 	return common.HexToAddress(ks.Address), nil
// }

// func AdvanceBlock() error {
// 	return fundAccount("0", []*common.Address{&common.Address{}})
// }

// func AdvanceBlocks(blockCount uint64) error {
// 	var i uint64
// 	for i = 0; i < blockCount; i++ {
// 		AdvanceBlock()
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	return nil
// }

// func AdvanceBlocksUntilDone(done chan bool) {
// 	ticker := time.NewTicker(time.Second)
// 	for {
// 		select {
// 		case <-done:
// 			ticker.Stop()
// 			return
// 		case <-ticker.C:
// 			AdvanceBlock()
// 		}
// 	}
// }
