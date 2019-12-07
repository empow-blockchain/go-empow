package iwallet

import (
	"github.com/spf13/cobra"
)

var gasUser string

var pledgeCmd = &cobra.Command{
	Use:     "gas-pledge amount",
	Aliases: []string{"pledge"},
	Short:   "Pledge EMPOW to obtain gas",
	Long:    `Pledge EMPOW to obtain gas`,
	Example: `  iwallet sys pledge 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet sys pledge 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt --gas_user EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
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
		if gasUser == "" {
			gasUser = accountName
		}
		return saveOrSendAction("gas.empow", "pledge", accountName, gasUser, args[0])
	},
}

var unpledgeCmd = &cobra.Command{
	Use:     "gas-unpledge amount",
	Aliases: []string{"unpledge"},
	Short:   "Undo pledge",
	Long:    `Undo pledge and get back the EMPOW pledged earlier`,
	Example: `  iwallet sys unpledge 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet sys pledge 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt --gas_user EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB`,
	Args: pledgeCmd.Args,
	RunE: func(cmd *cobra.Command, args []string) error {
		if gasUser == "" {
			gasUser = accountName
		}
		return saveOrSendAction("gas.empow", "unpledge", accountName, gasUser, args[0])
	},
}

func init() {
	systemCmd.AddCommand(pledgeCmd)
	pledgeCmd.Flags().StringVarP(&gasUser, "gas_user", "", "", "gas user that pledge EMPOW for (default is pledger himself/herself)")
	systemCmd.AddCommand(unpledgeCmd)
	unpledgeCmd.Flags().StringVarP(&gasUser, "gas_user", "", "", "gas user that earlier pledge for (default is pledger himself/herself)")
}
