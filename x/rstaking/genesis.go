package rstaking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stafihub/stafihub/x/rstaking/keeper"
	"github.com/stafihub/stafihub/x/rstaking/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, mintKeeper types.MintKeeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetInflationBase(ctx, genState.InflationBase)
	moduleAddress := authTypes.NewModuleAddress(types.ModuleName)
	params := mintKeeper.GetParams(ctx)
	if params.MintDenom != genState.CoinToBeBurned.Denom {
		panic("mint denom not equal coinToBeBurned denom")
	}
	if len(genState.GetValidatorWhitelist()) == 0 {
		panic("val_address_white_list empty")
	}

	for _, addr := range genState.GetValidatorWhitelist() {
		valAddr, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(fmt.Sprintf("valAddress format err, %s", err))
		}
		k.AddValAddressToWhitelist(ctx, valAddr)
	}

	balance := k.GetBankKeeper().GetBalance(ctx, moduleAddress, params.MintDenom)

	if balance.Amount.GT(sdk.ZeroInt()) {
		err := k.GetBankKeeper().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(balance))
		if err != nil {
			panic(err)
		}
	}
	err := k.GetBankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(genState.CoinToBeBurned))
	if err != nil {
		panic(err)
	}
	k.SetWhitelistSwitch(ctx, genState.WhitelistSwitch)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper, mintKeeper types.MintKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	inflationBase, found := k.GetInflationBase(ctx)
	if !found {
		inflationBase = sdk.ZeroInt()
	}
	genesis.InflationBase = inflationBase
	params := mintKeeper.GetParams(ctx)
	moduleAddress := authTypes.NewModuleAddress(types.ModuleName)
	balance := k.GetBankKeeper().GetBalance(ctx, moduleAddress, params.MintDenom)
	genesis.CoinToBeBurned = balance
	genesis.ValidatorWhitelist = k.GetValAddressWhitelist(ctx)
	genesis.WhitelistSwitch = k.GetWhitelistSwitch(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}