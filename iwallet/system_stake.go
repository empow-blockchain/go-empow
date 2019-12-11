package iwallet

import (
	"strconv"

	"github.com/spf13/cobra"
)

var stakeCmd = &cobra.Command{
	Use:     "stake amount",
	Aliases: []string{"stake"},
	Short:   "Stake EM",
	Long:    `Stake EM`,
	Example: `  iwallet sys stake 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "amount"); err != nil {
			return err
		}
		if err := checkFloat(cmd, args[0], "amount"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveOrSendAction("stake.empow", "stake", accountName, args[0])
	},
}

var stakeWithdrawCmd = &cobra.Command{
	Use:     "stake-withdraw packageID",
	Aliases: []string{"stake-withdraw"},
	Short:   "Withdraw stake with packageID",
	Long:    `Withdraw stake with packageID`,
	Example: `  iwallet sys stake-withdraw 0 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "packageID"); err != nil {
			return err
		}
		if err := checkFloat(cmd, args[0], "packageID"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		packageID, _ := strconv.ParseInt(args[0], 0, 64)
		return saveOrSendAction("stake.empow", "withdraw", accountName, packageID)
	},
}

var unstakeCmd = &cobra.Command{
	Use:     "unstake packageID",
	Aliases: []string{"unstake"},
	Short:   "Unstake with packageID",
	Long:    `Unstake with packageID`,
	Example: `  iwallet sys unstake 0 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "packageID"); err != nil {
			return err
		}
		if err := checkFloat(cmd, args[0], "packageID"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		packageID, _ := strconv.ParseInt(args[0], 0, 64)
		return saveOrSendAction("stake.empow", "unstake", accountName, packageID)
	},
}

var withdrawAllCmd = &cobra.Command{
	Use:     "stake-withdraw-all",
	Aliases: []string{"stake-withdraw-all"},
	Short:   "Withdraw all stake package",
	Long:    `Withdraw all stake package`,
	Example: `  iwallet sys stake-withdraw-all --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveOrSendAction("stake.empow", "withdrawAll", accountName)
	},
}

func init() {
	systemCmd.AddCommand(stakeCmd)
	systemCmd.AddCommand(stakeWithdrawCmd)
	systemCmd.AddCommand(unstakeCmd)
	systemCmd.AddCommand(withdrawAllCmd)
}
