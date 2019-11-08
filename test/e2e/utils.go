// Copyright 2018 Celer Network

package e2e

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func parseGatewayQueryResponse(resp *http.Response, cdc *codec.Codec) json.RawMessage {
	body, err := ioutil.ReadAll(resp.Body)
	tf.ChkErr(err, "failed to read http response")

	var responseWithHeight rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &responseWithHeight)
	return responseWithHeight.Result
}
