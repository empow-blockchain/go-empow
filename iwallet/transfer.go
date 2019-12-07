package iwallet

import (
	"github.com/spf13/cobra"
)

var memo string

var transferCmd = &cobra.Command{
	Use:     "transfer receiver amount",
	Aliases: []string{"trans"},
	Short:   "Transfer EMPOW",
	Long:    `Transfer EMPOW`,
	Example: `  iwallet transfer test1 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet transfer test1 100 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt --memo "just for test :D\n中文测试\n😏"`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "receiver", "amount"); err != nil {
			return err
		}
		if err := checkFloat(cmd, args[1], "amount"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveOrSendAction("token.empow", "transfer", "em", accountName, args[0], args[1], memo)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().StringVarP(&memo, "memo", "", "", "memo of transfer")
}
