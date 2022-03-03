package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/stafihub/stafihub/x/relayers/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateRelayer())
	cmd.AddCommand(CmdDeleteRelayer())
	cmd.AddCommand(CmdUpdateThreshold())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdCreateRelayer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-relayers [taipe] [denom] [addresses]",
		Short: "Add new relayers",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTaipe := args[0]

			argDenom := args[1]
			if err := sdk.ValidateDenom(argDenom); err != nil {
				return err
			}

			argValidators := strings.Split(args[2], ":")
			for _, v := range argValidators {
				_, err := sdk.AccAddressFromBech32(v)
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateRelayer(
				clientCtx.GetFromAddress(),
				argTaipe,
				argDenom,
				argValidators,
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

func CmdDeleteRelayer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-relayer [taipe] [denom] [address]",
		Short: "Delete a relayer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTaipe := args[0]

			argDenom := args[1]
			if sdk.ValidateDenom(argDenom) != nil {
				return nil
			}

			relAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteRelayer(
				clientCtx.GetFromAddress(),
				argTaipe,
				argDenom,
				relAddr,
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

func CmdUpdateThreshold() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-threshold [taipe] [denom] [value]",
		Short: "Update a threshold",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTaipe := args[0]

			argDenom := args[1]
			if sdk.ValidateDenom(argDenom) != nil {
				return nil
			}

			value, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateThreshold(
				clientCtx.GetFromAddress(),
				argTaipe,
				argDenom,
				uint32(value),
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
