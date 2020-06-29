package gov

import (
	govtypes "github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

// NewUpgradeProposalHandler creates a governance handler to manage new proposal types.
// It enables UpgradeProposal to propose an Upgrade, and CancelUpgradeProposal
// to abort a previously voted upgrade.
func NewUpgradeProposalHandler(k upgrade.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case govtypes.UpgradeProposal:
			return handleUpgradeProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized software upgrade proposal content type: %T", c)
		}
	}
}

func handleUpgradeProposal(ctx sdk.Context, k upgrade.Keeper, p govtypes.UpgradeProposal) error {
	return k.ScheduleUpgrade(ctx, p.Plan)
}
