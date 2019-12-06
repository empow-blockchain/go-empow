package iwallet

import (
	"fmt"

	"github.com/empow-blockchain/go-empow/sdk"

	"github.com/spf13/cobra"
)

// accountInfoCmd represents the balance command.
var accountInfoCmd = &cobra.Command{
	Use:     "balance address",
	Short:   "Check the information of a specified address",
	Long:    `Check the information of a specified address`,
	Example: `  iwallet balance EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "address"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		info, err := iwalletSDK.GetAccountInfo(id)
		if err != nil {
			return err
		}
		fmt.Println(sdk.MarshalTextString(info))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountInfoCmd)
}
