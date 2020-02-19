package genesis

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/empow-blockchain/go-empow/account"
	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/core/block"
	"github.com/empow-blockchain/go-empow/core/contract"
	"github.com/empow-blockchain/go-empow/core/tx"
	"github.com/empow-blockchain/go-empow/crypto"
	"github.com/empow-blockchain/go-empow/db"
	"github.com/empow-blockchain/go-empow/ilog"
	"github.com/empow-blockchain/go-empow/verifier"
	"github.com/empow-blockchain/go-empow/vm/native"
)

// GenesisTxExecTime is the maximum execution time of a transaction in genesis block
var GenesisTxExecTime = 10 * time.Second

// GenGenesisByFile is create a genesis block by config file
func GenGenesisByFile(db db.MVCCDB, path string) (*block.Block, error) {
	v := common.LoadYamlAsViper(filepath.Join(path, "genesis.yml"))
	genesisConfig := &common.GenesisConfig{}
	if err := v.Unmarshal(genesisConfig); err != nil {
		ilog.Fatalf("Unable to decode into struct, %v", err)
	}
	genesisConfig.ContractPath = filepath.Join(path, "contract")
	return GenGenesis(db, genesisConfig)
}

func compile(id string, path string, name string) (*contract.Contract, error) {
	if id == "" || path == "" || name == "" {
		return nil, fmt.Errorf("arguments is error, id:%v, path:%v, name:%v", id, path, name)
	}
	cFilePath := filepath.Join(path, name)
	cAbiPath := filepath.Join(path, name+".abi")
	return contract.Compile(id, cFilePath, cAbiPath)
}

func genGenesisTx(gConf *common.GenesisConfig) (*tx.Tx, *account.Account, error) {
	witnessInfo := gConf.WitnessInfo
	// prepare actions
	var acts []*tx.Action
	adminInfo := gConf.AdminInfo

	// deploy token.empow
	acts = append(acts, tx.NewAction("system.empow", "initSetCode",
		fmt.Sprintf(`["%v", "%v"]`, "token.empow", native.SystemContractABI("token.empow", "1.0.0").B64Encode())))
	acts = append(acts, tx.NewAction("system.empow", "initSetCode",
		fmt.Sprintf(`["%v", "%v"]`, "token721.empow", native.SystemContractABI("token721.empow", "1.0.0").B64Encode())))
	// deploy gas.empow
	acts = append(acts, tx.NewAction("system.empow", "initSetCode",
		fmt.Sprintf(`["%v", "%v"]`, "gas.empow", native.SystemContractABI("gas.empow", "1.0.0").B64Encode())))
	// deploy issue.empow and create em token
	code, err := compile("issue.empow", gConf.ContractPath, "issue.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "issue.empow", code.B64Encode())))
	tokenInfo := gConf.TokenInfo
	tokenHolder := append(witnessInfo, adminInfo)
	params := []interface{}{
		adminInfo.Address,
		tokenInfo,
		tokenHolder,
	}
	b, _ := json.Marshal(params)
	acts = append(acts, tx.NewAction("issue.empow", "initGenesis", string(b)))
	// deploy auth.empow
	code, err = compile("auth.empow", gConf.ContractPath, "account.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "auth.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("auth.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deploy domain.empow
	acts = append(acts, tx.NewAction("system.empow", "initSetCode",
		fmt.Sprintf(`["%v", "%v"]`, "domain.empow", native.SystemContractABI("domain.empow", "0.0.0").B64Encode())))

	// deloy stake.empow
	code, err = compile("stake.empow", gConf.ContractPath, "stake.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "stake.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("stake.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deloy social.empow
	code, err = compile("social.empow", gConf.ContractPath, "social.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "social.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("social.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// new account
	acts = append(acts, tx.NewAction("auth.empow", "signUp", fmt.Sprintf(`["%v", "%v", "%v"]`, adminInfo.Address, adminInfo.Owner, adminInfo.Active)))
	// new account
	foundationInfo := gConf.FoundationInfo
	acts = append(acts, tx.NewAction("auth.empow", "signUp", fmt.Sprintf(`["%v", "%v", "%v"]`, foundationInfo.Address, foundationInfo.Owner, foundationInfo.Active)))

	for _, v := range witnessInfo {
		acts = append(acts, tx.NewAction("auth.empow", "signUp", fmt.Sprintf(`["%v", "%v", "%v"]`, v.Address, v.Owner, v.Active)))
	}
	invalidPubKey := "0"
	deadAccount := account.NewAccount("deadaddr")
	acts = append(acts, tx.NewAction("auth.empow", "signUp", fmt.Sprintf(`["%v", "%v", "%v"]`, deadAccount.Address, invalidPubKey, invalidPubKey)))

	// deploy bonus.empow
	code, err = compile("bonus.empow", gConf.ContractPath, "bonus.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "bonus.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("bonus.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deloy vote_point.empow
	code, err = compile("vote_point.empow", gConf.ContractPath, "vote_point.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "vote_point.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("vote_point.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deploy vote.empow
	code, err = compile("vote.empow", gConf.ContractPath, "vote_common.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "vote.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("vote.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deploy vote_producer.empow
	code, err = compile("vote_producer.empow", gConf.ContractPath, "vote_producer.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "vote_producer.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("vote_producer.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deploy base.empow
	code, err = compile("base.empow", gConf.ContractPath, "base.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "base.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("base.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))

	// deploy exchange.empow
	code, err = compile("exchange.empow", gConf.ContractPath, "exchange.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "exchange.empow", code.B64Encode())))

	for _, v := range witnessInfo {
		acts = append(acts, tx.NewAction("vote_producer.empow", "initProducer", fmt.Sprintf(`["%v", "%v"]`, v.Address, v.SignatureBlock)))
	}

	// pledge gas for admin
	gasPledgeAmount := 100
	acts = append(acts, tx.NewAction("gas.empow", "pledge", fmt.Sprintf(`["%v", "%v", "%v"]`, adminInfo.Address, adminInfo.Address, gasPledgeAmount)))

	// deploy ram.empow
	code, err = compile("ram.empow", gConf.ContractPath, "ram.js")
	if err != nil {
		return nil, nil, err
	}
	acts = append(acts, tx.NewAction("system.empow", "initSetCode", fmt.Sprintf(`["%v", "%v"]`, "ram.empow", code.B64Encode())))
	acts = append(acts, tx.NewAction("ram.empow", "initAdmin", fmt.Sprintf(`["%v"]`, adminInfo.Address)))
	var initialTotal int64 = 128 * 1024 * 1024 * 1024                           // 128GB at first
	var increaseInterval int64 = 10 * 60                                        // increase every 10 mins
	var increaseAmount int64 = 10 * (64 * 1024 * 1024 * 1024) / (365 * 24 * 60) // 64GB per year
	var reserveRAM = initialTotal * 3 / 10                                      // reserve for foundation
	acts = append(acts, tx.NewAction("ram.empow", "issue", fmt.Sprintf(`[%v, %v, %v, %v]`, initialTotal, increaseInterval, increaseAmount, reserveRAM)))

	adminInitialRAM := 100000
	acts = append(acts, tx.NewAction("ram.empow", "buy", fmt.Sprintf(`["%v", "%v", %v]`, adminInfo.Address, adminInfo.Address, adminInitialRAM)))
	acts = append(acts, tx.NewAction("token.empow", "transfer", fmt.Sprintf(`["ram","ram.empow", "%v", "%v", ""]`, foundationInfo.Address, reserveRAM)))

	for _, v := range witnessInfo {
		acts = append(acts, tx.NewAction("ram.empow", "buy", fmt.Sprintf(`["%v", "%v", %v]`, adminInfo.Address, v.Address, adminInitialRAM)))
	}

	acts = append(acts, tx.NewAction("gas.empow", "pledge", fmt.Sprintf(`["%v", "%v", "%v"]`, adminInfo.Address, foundationInfo.Address, gasPledgeAmount)))
	for _, v := range witnessInfo {
		acts = append(acts, tx.NewAction("gas.empow", "pledge", fmt.Sprintf(`["%v", "%v", "%v"]`, adminInfo.Address, v.Address, gasPledgeAmount)))
	}

	trx := tx.NewTx(acts, nil, 1000000000, 100, 0, 0, tx.ChainID)
	trx.Time = 0
	trx, err = tx.SignTx(trx, deadAccount.Address, []*account.KeyPair{})
	if err != nil {
		return nil, nil, err
	}
	trx.AmountLimit = append(trx.AmountLimit, &contract.Amount{Token: "*", Val: "unlimited"})
	return trx, deadAccount, nil
}

// GenGenesis is create a genesis block
func GenGenesis(db db.MVCCDB, gConf *common.GenesisConfig) (*block.Block, error) {
	t, err := time.Parse(time.RFC3339, gConf.InitialTimestamp)
	if err != nil {
		ilog.Fatalf("invalid genesis initial time string %v (%v).", gConf.InitialTimestamp, err)
		return nil, err
	}
	trx, _, err := genGenesisTx(gConf)
	if err != nil {
		return nil, err
	}

	blockHead := block.BlockHead{
		Version:    block.V0,
		ParentHash: nil,
		Number:     0,
		Witness:    "0",
		Time:       t.UnixNano(),
	}
	v := verifier.Verifier{}
	txr, err := v.Exec(&blockHead, db, trx, GenesisTxExecTime)
	if err != nil || txr.Status.Code != tx.Success {
		return nil, fmt.Errorf("exec tx failed, stop the pogram. err: %v, receipt: %v", err, txr)
	}
	blk := &block.Block{
		Head:     &blockHead,
		Sign:     &crypto.Signature{},
		Txs:      []*tx.Tx{trx},
		Receipts: []*tx.TxReceipt{txr},
	}
	blk.Head.TxMerkleHash = blk.CalculateTxMerkleHash()
	blk.Head.TxReceiptMerkleHash = blk.CalculateTxReceiptMerkleHash()
	blk.CalculateHeadHash()
	db.Commit(string(blk.HeadHash()))
	return blk, nil
}
