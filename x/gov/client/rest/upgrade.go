package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func UpgradeProposalRESTHandler(cliCtx context.CLIContext) ProposalRESTHandler {
	return ProposalRESTHandler{
		SubRoute: "upgrade",
		Handler:  postUpgradeProposalHandlerFn(cliCtx),
	}
}
func postUpgradeProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
