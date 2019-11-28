// Copyright 2018 Celer Network

package e2e

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"context"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tf "github.com/celer-network/sgn/testing"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func sleepBlocksWithLog(count time.Duration, waitFor string) {
	sleepWithLog(count * sgnBlockInterval, waitFor)
}

func parseGatewayQueryResponse(resp *http.Response, cdc *codec.Codec) (json.RawMessage, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseWithHeight rest.ResponseWithHeight
	err = cdc.UnmarshalJSON(body, &responseWithHeight)
	if err != nil {
		return nil, err
	}

	return responseWithHeight.Result, nil
}

func initializeCandidate(auth *bind.TransactOpts, sgnAddr sdk.AccAddress) error {
	conn := tf.EthClient.Client
	guardContract := tf.EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	log.Info("Call initializeCandidate on guard contract using the validator eth address...")
	tx, err := guardContract.InitializeCandidate(auth, big.NewInt(1), sgnAddr.Bytes())
	if err != nil {
		return err
	}

	tf.WaitMinedWithChk(ctx, conn, tx, blockDelay, "InitializeCandidate")
	sleepBlocksWithLog(5, "sgn syncing InitializeCandidate event on mainchain")
	return nil
}

func delegateStake(fromAuth *bind.TransactOpts, toEthAddress mainchain.Addr, amt *big.Int) error {
	conn := tf.EthClient.Client
	guardContract := tf.EthClient.Guard
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	log.Info("Call delegate on guard contract to delegate stake to the validator eth address...")
	tx, err := celrContract.Approve(fromAuth, guardAddr, amt)
	if err != nil {
		return err
	}
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve CELR to Guard contract")
	tx, err = guardContract.Delegate(fromAuth, toEthAddress, amt)
	if err != nil {
		return err
	}

	tf.WaitMinedWithChk(ctx, conn, tx, 3 * blockDelay, "Delegate to validator")
	sleepWithLog(10, "sgn syncing Delegate event on mainchain")
	return nil
}