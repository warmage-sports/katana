package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/stafihub/stafihub/x/mining/types"
)

var _ = strconv.Itoa(0)

func CmdAddStakePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-stake-pool [stake-token-denom] [max-reward-pools] [min-total-reward-amount]",
		Short: "Add stake pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStakeTokenDenom := args[0]
			argMaxRewardPools, err := sdk.ParseUint(args[1])
			if err != nil {
				return err
			}
			argMinTotalRewardAmount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("minTotalRewardAmount format err")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddStakePool(
				clientCtx.GetFromAddress().String(),
				argStakeTokenDenom,
				uint32(argMaxRewardPools.Uint64()),
				argMinTotalRewardAmount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
