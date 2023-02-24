package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oldfurya/furya/x/mining/types"
	sudotypes "github.com/oldfurya/furya/x/sudo/types"
)

func (k msgServer) AddStakeToken(goCtx context.Context, msg *types.MsgAddStakeToken) (*types.MsgAddStakeTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.sudoKeeper.IsAdmin(ctx, msg.Creator) {
		return nil, sudotypes.ErrCreatorNotAdmin
	}

	k.Keeper.AddStakeToken(ctx, msg.Denom)

	return &types.MsgAddStakeTokenResponse{}, nil
}
