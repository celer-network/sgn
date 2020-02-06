package testing

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	mnemonicEntropySize = 256

	passphraseFlag  = "passphrase"
	namePrefixFlag  = "prefix"
	countFlag       = "count"
	genesisCoinFlag = "coin"
)

func AccountsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Add accounts in batch",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			kb, err := keys.NewKeyBaseFromHomeFlag()
			if err != nil {
				return err
			}

			var addresses []string
			passphrase := viper.GetString(passphraseFlag)
			np := viper.GetString(namePrefixFlag)
			count := viper.GetInt(countFlag)

			for i := 0; i < count; i++ {
				name := fmt.Sprintf("%s_%d", np, i)
				info, err := kb.Get(name)
				if err == nil {
					log.Infof("Account %s has existed", name)
					printAccount(info)
					addresses = append(addresses, info.GetAddress().String())
					continue
				}

				entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
				if err != nil {
					return err
				}

				mnemonic, err := bip39.NewMnemonic(entropySeed[:])
				if err != nil {
					return err
				}

				info, err = kb.CreateAccount(name, mnemonic, passphrase, "", 0, 0)
				if err != nil {
					return err
				}

				log.Infof("Account %s created", name)
				printAccount(info)
				addresses = append(addresses, info.GetAddress().String())
			}

			jsonString, err := keys.MarshalJSON(addresses)
			if err != nil {
				return
			}
			log.Infoln("All addresses", string(jsonString))

			genesisCoin := viper.GetString(genesisCoinFlag)
			if genesisCoin != "" {
				for _, address := range addresses {
					cmd := exec.Command("sgn", "add-genesis-account", address, genesisCoin)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Dir, _ = filepath.Abs(".")
					if err := cmd.Run(); err == nil {
						log.Infof("Add address %s to genesis", address)
					}
				}
			}

			return
		},
	}

	cmd.Flags().String(passphraseFlag, "12345678", "account passphrase")
	cmd.Flags().String(namePrefixFlag, "transactor", "account prefix")
	cmd.Flags().Int(countFlag, 1, "account count")
	cmd.Flags().String(genesisCoinFlag, "", "amount of coin adding to genesis for the account")
	return cmd
}

func printAccount(info cKeys.Info) {
	out, err := cKeys.Bech32KeyOutput(info)
	if err != nil {
		return
	}

	jsonString, err := keys.MarshalJSON(out)
	if err != nil {
		return
	}

	log.Infof(string(jsonString))
	return
}
