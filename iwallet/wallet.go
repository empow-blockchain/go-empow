package iwallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/empow-blockchain/go-empow/account"
	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/sdk"
	"github.com/spf13/cobra"
)

var (
	ownerKey         string
	activeKey        string
	initialRAM       int64
	initialBalance   int64
	initialGasPledge int64
)

type acc struct {
	Address string
	KeyPair *key
}

type accounts struct {
	Dir     string
	Account []*acc
}

// walletCmd represents the account command.
var accountCmd = &cobra.Command{
	Use:     "wallet",
	Aliases: []string{"wallet"},
	Short:   "Wallet manager",
	Long:    `Manage wallet in local storage`,
}

var encrypt bool
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create new wallet and save to local storage",
	Long:    `Create new wallet and save to local storage`,
	Example: `  iwallet wallet create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := account.NewAddress(nil, sdk.GetSignAlgoByName(signAlgo))
		if err != nil {
			return fmt.Errorf("failed to new key pair: %v", err)
		}

		var k key
		k.Algorithm = n.Algorithm.String()
		k.Address = n.Address
		k.Pubkey = common.Base58Encode(n.Pubkey)
		k.Seckey = common.Base58Encode(n.Seckey)

		if err != nil {
			return fmt.Errorf("failed create wallet: %v", err)
		}

		kp, err := NewKeyPairInfo(k.Seckey, k.Algorithm)
		if err != nil {
			return err
		}

		acc := AccountInfo{Address: k.Address, Keypairs: make(map[string]*KeyPairInfo, 0)}

		acc.Keypairs["active"] = kp
		acc.Keypairs["owner"] = kp

		acc.save(encrypt)

		fmt.Printf("Created %v\n", k.Address)
		return nil
	},
}

var viewCmd = &cobra.Command{
	Use:   "view [<address>]",
	Short: "View address by name or omit to show all addresses",
	Long:  `View address by name or omit to show all addresses`,
	Example: `  iwallet wallet view <address>
  iwallet wallet view`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := getAccountDir()
		if err != nil {
			return fmt.Errorf("failed to get address dir: %v", err)
		}
		a := accounts{}
		a.Dir = dir
		addAcc := func(ac *AccountInfo) {
			var k key
			k.Algorithm = ac.Keypairs["active"].KeyType
			k.Address = ac.Address
			k.Pubkey = ac.Keypairs["active"].PubKey
			if ac.isEncrypted() {
				k.Seckey = "---encrypted secret key---"
			} else {
				k.Seckey = ac.Keypairs["active"].RawKey
			}
			a.Account = append(a.Account, &acc{ac.Address, &k})
		}
		if len(args) < 1 {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, f := range files {
				ac, err := loadAccountFromFile(dir+"/"+f.Name(), false)
				if err != nil {
					continue
				}
				addAcc(ac)
			}
		} else {
			name := args[0]
			ac, err := loadAccountByName(name, false)
			if err != nil {
				return err
			}
			addAcc(ac)
		}
		info, err := json.MarshalIndent(a, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(info))
		return nil
	},
}

var importCmd = &cobra.Command{
	Use:   "import privateKey",
	Short: "Import an address by private key",
	Long:  `Import an address by private key`,
	Example: `  iwallet wallet import <private_key>
  iwallet wallet import <private_key>`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "privateKey"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		acc := AccountInfo{Address: "test", Keypairs: make(map[string]*KeyPairInfo, 0)}
		keys := strings.Split(args[0], ",")
		if len(keys) == 1 {
			key := keys[0]
			if len(strings.Split(key, ":")) != 1 {
				return fmt.Errorf("importing one key need not specifying permission")
			}
			kp, err := NewKeyPairInfo(key, signAlgo)
			if err != nil {
				return err
			}
			acc.Keypairs["active"] = kp
			acc.Keypairs["owner"] = kp
		} else {
			for _, permAndKey := range keys {
				splits := strings.Split(permAndKey, ":")
				if len(splits) != 2 {
					return fmt.Errorf("importing more than one keys need specifying permissions")
				}
				kp, err := NewKeyPairInfo(splits[1], signAlgo)
				if err != nil {
					return err
				}
				acc.Keypairs[splits[0]] = kp
			}
		}

		address := account.PubkeyToAddress(common.Base58Decode(acc.Keypairs["active"].PubKey))
		acc.Address = address
		err := acc.save(encrypt)
		if err != nil {
			return fmt.Errorf("failed to save address: %v", err)
		}
		fmt.Printf("import address %v done\n", address)
		return nil
	},
}

var dumpKeyCmd = &cobra.Command{
	Use:     "dumpkey address",
	Short:   "Print private key of the address to stdout",
	Long:    "Print private key of the address to stdout",
	Example: `  iwallet wallet dumpkey <address>`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "address"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		acc, err := loadAccountByName(args[0], true)
		if err != nil {
			return err
		}
		for k, v := range acc.Keypairs {
			fmt.Printf("%v:%v\n", k, v.RawKey)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete <address>",
	Aliases: []string{"del"},
	Short:   "Delete an address",
	Long:    `Delete an address`,
	Example: `  iwallet wallet delete EMvFUnDToqD4rFhckJCfkTHuufdSFPQpabrJs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "address"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		dir, err := getAccountDir()
		if err != nil {
			return fmt.Errorf("failed to get address dir: %v", err)
		}
		found := false
		sufs := []string{".json"}
		for _, algo := range ValidSignAlgos {
			sufs = append(sufs, "_"+algo)
		}
		for _, suf := range sufs {
			f := fmt.Sprintf("%s/%s%s", dir, address, suf)
			err = os.Remove(f)
			if err == nil {
				found = true
				fmt.Println("File", f, "has been removed.")
			}
			err = os.Remove(f + ".id")
			if err == nil {
				fmt.Println("File", f+".id", "has been removed.")
			}
			err = os.Remove(f + ".pub")
			if err == nil {
				fmt.Println("File", f+".pub", "has been removed.")
			}
		}
		if found {
			fmt.Println("Successfully deleted <", address, ">.")
		} else {
			fmt.Println("Address <", address, "> does not exist.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.PersistentFlags().BoolVarP(&encrypt, "encrypt", "", false, "whether to encrypt local key file")
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(importCmd)
	accountCmd.AddCommand(viewCmd)
	accountCmd.AddCommand(deleteCmd)
	accountCmd.AddCommand(dumpKeyCmd)
}
