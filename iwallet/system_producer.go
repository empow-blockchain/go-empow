package iwallet

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	rpcpb "github.com/empow-blockchain/go-empow/rpc/pb"
	"github.com/empow-blockchain/go-empow/sdk"
)

var location string
var url string
var networkID string
var isPartner bool
var publicKey string
var target string

var voteCmd = &cobra.Command{
	Use:     "producer-vote producerID amount",
	Aliases: []string{"vote"},
	Short:   "Vote a producer",
	Long:    `Vote a producer by given amount of EMPOWs`,
	Example: `  iwallet sys vote EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB 1000000 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "producerID", "amount"); err != nil {
			return err
		}
		if err := checkFloat(cmd, args[1], "amount"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveOrSendAction("vote_producer.empow", "vote", accountName, args[0], args[1])
	},
}
var unvoteCmd = &cobra.Command{
	Use:     "producer-unvote producerID amount",
	Aliases: []string{"unvote"},
	Short:   "Unvote a producer",
	Long:    `Unvote a producer by given amount of EMPOWs`,
	Example: `  iwallet sys unvote EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB 1000000 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args:    voteCmd.Args,
	RunE: func(cmd *cobra.Command, args []string) error {
		return saveOrSendAction("vote_producer.empow", "unvote", accountName, args[0], args[1])
	},
}

var registerCmd = &cobra.Command{
	Use:     "producer-register publicKey",
	Aliases: []string{"register", "reg"},
	Short:   "Register as producer",
	Long:    `Register as producer`,
	Example: `  iwallet sys register XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet sys register XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt --location PEK --url iost.io --net_id 123 --partner`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "publicKey"); err != nil {
			return err
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "applyRegister", target, args[0], location, url, networkID, !isPartner)
	},
}
var unregisterCmd = &cobra.Command{
	Use:     "producer-unregister",
	Aliases: []string{"unregister", "unreg"},
	Short:   "Unregister from a producer",
	Long:    `Unregister from a producer`,
	Example: `  iwallet sys unregister --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "applyUnregister", target)
	},
}
var pcleanCmd = &cobra.Command{
	Use:     "producer-clean",
	Aliases: []string{"pclean"},
	Short:   "Clean producer info",
	Long:    `Clean producer info`,
	Example: `  iwallet sys pclean --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "unregister", target)
	},
}

var ploginCmd = &cobra.Command{
	Use:     "producer-login",
	Aliases: []string{"plogin"},
	Short:   "Producer login as online state",
	Long:    `Producer login as online state`,
	Example: `  iwallet sys plogin --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "logInProducer", target)
	},
}
var plogoutCmd = &cobra.Command{
	Use:     "producer-logout",
	Aliases: []string{"plogout"},
	Short:   "Producer logout as offline state",
	Long:    `Producer logout as offline state`,
	Example: `  iwallet sys plogout --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "logOutProducer", target)
	},
}

func getProducerVoteInfo(account string) (*rpcpb.GetProducerVoteInfoResponse, error) {
	return iwalletSDK.GetProducerVoteInfo(&rpcpb.GetProducerVoteInfoRequest{
		Account:        account,
		ByLongestChain: useLongestChain,
	})
}

var pinfoCmd = &cobra.Command{
	Use:     "producer-info producerID",
	Aliases: []string{"pinfo"},
	Short:   "Show producer info",
	Long:    `Show producer info`,
	Example: `  iwallet sys pinfo EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "producerID"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := getProducerVoteInfo(args[0])
		if err != nil {
			return err
		}
		fmt.Println(sdk.MarshalTextString(info))
		return nil
	},
}

var plistCmd = &cobra.Command{
	Use:     "producer-list",
	Aliases: []string{"plist"},
	Short:   "Show current/pending producer list",
	Long:    `Show current/pending producer list`,
	Example: `  iwallet sys plist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := iwalletSDK.Connect(); err != nil {
			return err
		}
		defer iwalletSDK.CloseConn()
		chainInfo, err := iwalletSDK.GetChainInfo()
		if err != nil {
			return fmt.Errorf("cannot get chain info: %v", err)
		}

		var getWitnessName = func(pks []string) []string {
			result := make([]string, len(pks))
			var wg sync.WaitGroup
			wg.Add(len(pks))
			for i, producerKey := range pks {
				i, producerKey := i, producerKey // bind current value to closure
				go func() {
					response, err := getContractStorage("vote_producer.empow", "producerKeyToId", producerKey)
					if err != nil {
						fmt.Printf("cannot get producer id of %v: %v", producerKey, err)
						return
					}
					result[i] = response.Data
					wg.Done()
				}()
			}
			wg.Wait()
			return result
		}

		currentPlist, pendingPlist := getWitnessName(chainInfo.WitnessList), getWitnessName(chainInfo.PendingWitnessList)
		fmt.Println("Current producer list:", currentPlist)
		fmt.Println("Pending producer list:", pendingPlist)
		return nil
	},
}

var pupdateCmd = &cobra.Command{
	Use:     "producer-update",
	Aliases: []string{"pupdate"},
	Short:   "Update producer info",
	Long:    `Update producer info`,
	Example: `  iwallet sys pupdate --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet sys pupdate --address EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB --pubkey XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
  iwallet sys pupdate --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt --location PEK --url iost.io --net_id 123`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		info, err := getProducerVoteInfo(target)
		if err != nil {
			return err
		}
		if publicKey == "" {
			publicKey = info.Pubkey
		}
		if location == "" {
			location = info.Loc
		}
		if url == "" {
			url = info.Url
		}
		if networkID == "" {
			networkID = info.NetId
		}
		return saveOrSendAction("vote_producer.empow", "updateProducer", target, publicKey, location, url, networkID)
	},
}

var predeemCmd = &cobra.Command{
	Use:     "producer-redeem [amount]",
	Aliases: []string{"predeem"},
	Short:   "Redeem the contribution value obtained by the block producing to EMPOW tokens",
	Long: `Redeem the contribution value obtained by the block producing to EMPOW tokens
	Omitting amount argument or zero amount will redeem all contribution value.`,
	Example: `  iwallet sys producer-redeem --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt
  iwallet sys producer-redeem 10 --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			if err := checkFloat(cmd, args[0], "amount"); err != nil {
				return err
			}
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		amount := "0"
		if len(args) > 0 {
			amount = args[0]
		}
		return saveOrSendAction("bonus.empow", "exchangeEMPOW", accountName, amount)
	},
}

var pwithdrawCmd = &cobra.Command{
	Use:     "producer-withdraw",
	Aliases: []string{"pwithdraw"},
	Short:   "Withdraw all voting reward for producer",
	Long:    `Withdraw all voting reward for producer`,
	Example: `  iwallet sys producer-withdraw --address EM2ZsDPRrJHHKgc7w719Ds9X9Z7QCcuMB4bFxMynDR2TYfQqt`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			target = accountName
		}
		return saveOrSendAction("vote_producer.empow", "candidateWithdraw", target)
	},
}

func init() {
	systemCmd.AddCommand(voteCmd)
	systemCmd.AddCommand(unvoteCmd)

	systemCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")
	registerCmd.Flags().StringVarP(&location, "location", "", "", "location info")
	registerCmd.Flags().StringVarP(&url, "url", "", "", "url address")
	registerCmd.Flags().StringVarP(&networkID, "net_id", "", "", "network ID")
	registerCmd.Flags().BoolVarP(&isPartner, "partner", "", false, "if is partner instead of producer")
	systemCmd.AddCommand(unregisterCmd)
	unregisterCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")
	systemCmd.AddCommand(pcleanCmd)
	pcleanCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")

	systemCmd.AddCommand(ploginCmd)
	ploginCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")
	systemCmd.AddCommand(plogoutCmd)
	plogoutCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")

	systemCmd.AddCommand(pinfoCmd)
	systemCmd.AddCommand(plistCmd)

	systemCmd.AddCommand(pupdateCmd)
	pupdateCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")
	pupdateCmd.Flags().StringVarP(&publicKey, "pubkey", "", "", "publick key")
	pupdateCmd.Flags().StringVarP(&location, "location", "", "", "location info")
	pupdateCmd.Flags().StringVarP(&url, "url", "", "", "url address")
	pupdateCmd.Flags().StringVarP(&networkID, "net_id", "", "", "network ID")

	systemCmd.AddCommand(predeemCmd)
	systemCmd.AddCommand(pwithdrawCmd)
	pwithdrawCmd.Flags().StringVarP(&target, "target", "", "", "target address (default is the account by flag --address himself/herself)")
}
