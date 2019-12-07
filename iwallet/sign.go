package iwallet

import (
	"fmt"

	"github.com/spf13/cobra"

	rpcpb "github.com/empow-blockchain/go-empow/rpc/pb"
	"github.com/empow-blockchain/go-empow/sdk"
)

// signCmd represents the command used to sign a transaction.
var signCmd = &cobra.Command{
	Use:   "sign txFile keyFile outputFile",
	Short: "Sign a tx and save the signature",
	Long:  `Sign a transaction loaded from given txFile with keyFile(address json file or private key file) and save the signature as outputFile`,
	Example: `  iwallet sign tx.json ~/.iwallet/EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt.json sign.json
  iwallet sign tx.json ~/.iwallet/EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt_ed25519 sign.json`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "txFile", "keyFile", "outputFile"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		txFile := args[0]
		signKeyFile := args[1]
		outputFile := args[2]

		trx := &rpcpb.TransactionRequest{}
		err := sdk.LoadProtoStructFromJSONFile(txFile, trx)
		if err != nil {
			return fmt.Errorf("failed to load transaction file %v: %v", txFile, err)
		}
		accInfo, err := loadAccountFromFile(signKeyFile, true)
		if err != nil {
			return fmt.Errorf("failed to load addresses from file %v: %v", signKeyFile, err)
		}
		kp, err := accInfo.Keypairs["active"].toKeyPair()
		if err != nil {
			return fmt.Errorf("failed to get key pair from file %v: %v", signKeyFile, err)
		}
		sig := sdk.GetSignatureOfTx(trx, kp, asPublisherSign)
		if verbose {
			fmt.Println("Signature:")
			fmt.Println(sdk.MarshalTextString(sig))
		}
		err = sdk.SaveProtoStructToJSONFile(sig, outputFile)
		if err != nil {
			return fmt.Errorf("failed to save signature as file %v: %v", outputFile, err)
		}
		fmt.Println("Successfully saved signature as:", outputFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(signCmd)
}
