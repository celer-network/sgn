package gateway

import (
	"fmt"
	"net/http"

	"github.com/celer-network/sgn/x/global"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func (rs *RestServer) registerQueryRoutes() {
	rs.Mux.HandleFunc(
		"/global/latestBlock",
		latestBlockHandlerFn(rs),
	).Methods("GET")
}

// http request handler to query latest block
func latestBlockHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", global.ModuleName, global.QueryLatestBlock)
		res, _, err := rs.transactor.CliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, rs.transactor.CliCtx, res)
	}
}
