package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oldfurya/furya/x/ledger/types"
	sudotypes "github.com/oldfurya/furya/x/sudo/types"
)

func (k msgServer) SetWithdrawalAddr(goCtx context.Context, msg *types.MsgSetWithdrawalAddr) (*types.MsgSetWithdrawalAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.sudoKeeper.IsAdmin(ctx, msg.Creator) {
		return nil, sudotypes.ErrCreatorNotAdmin
	}

	icaPoolDetail, found := k.GetIcaPoolByDelegationAddr(ctx, msg.DelegationAddr)
	if !found {
		return nil, types.ErrIcaPoolNotFound
	}

	if icaPoolDetail.Status != types.IcaPoolStatusCreateTwo {
		return nil, types.ErrIcaPoolStatusUnmatch
	}

	err := k.Keeper.SetWithdrawAddressOnHost(
		ctx,
		icaPoolDetail.DelegationAccount.Owner,
		icaPoolDetail.DelegationAccount.CtrlConnectionId,
		icaPoolDetail.DelegationAccount.Address,
		icaPoolDetail.WithdrawalAccount.Address)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetWithdrawalAddrResponse{}, nil
}
