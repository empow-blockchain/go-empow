package gobang

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/empow-blockchain/go-empow/sdk"

	"github.com/empow-blockchain/go-empow/test/performance/call"

	"encoding/json"

	"math/rand"

	"github.com/empow-blockchain/go-empow/account"
	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/core/tx"
	"github.com/empow-blockchain/go-empow/crypto"
	"github.com/empow-blockchain/go-empow/ilog"
	rpcpb "github.com/empow-blockchain/go-empow/rpc/pb"
)

var rootKey = "2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1"
var rootID = "EM2ZsSw7RWYC229Z1ib7ujKhken9GFR7dBkTTEbBWMKeLpVas"
var rootAcc *account.KeyPair
var contractID string
var iostSDK = sdk.NewEMPOWDevSDK()
var testID = "i" + strconv.FormatInt(time.Now().Unix(), 10)
var testAcc *account.KeyPair

type gobangHandle struct{}

const (
	chainID uint32 = 1024
)

func init() {
	gobang := new(gobangHandle)
	call.Register("gobang", gobang)
	iostSDK.SetChainID(chainID)
}

// Publish ...
func (t *gobangHandle) Prepare() error {
	var err error
	rootAcc, _ = account.NewKeyPair(common.Base58Decode(rootKey), crypto.Ed25519)
	codePath := os.Getenv("GOPATH") + "/src/github.com/empow-blockchain/go-empow/test/performance/handles/gobang/gobang.js"
	abiPath := codePath + ".abi"
	iostSDK.SetServer(call.GetClient(0).Addr())
	iostSDK.SetAccount("EM2ZsSw7RWYC229Z1ib7ujKhken9GFR7dBkTTEbBWMKeLpVas", rootAcc)
	iostSDK.SetTxInfo(100000, 1, 90, 0, nil)
	iostSDK.SetCheckResult(true, 3, 10)
	testAcc, err = account.NewKeyPair(nil, crypto.Ed25519)
	if err != nil {
		return err
	}
	k := testAcc.ReadablePubkey()
	_, err = iostSDK.CreateNewAccount(testID, k, k, 1000000, 10000, 100000)
	if err != nil {
		return err
	}
	err = iostSDK.PledgeForGasAndRAM(1500000, 10000000)
	if err != nil {
		return err
	}
	iostSDK.SetAccount(testID, testAcc)
	_, txHash, err := iostSDK.PublishContract(codePath, abiPath, "", false, "")
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(50) * time.Second)
	client := call.GetClient(0)
	resp, err := client.GetTxReceiptByTxHash(context.Background(), &rpcpb.TxHashRequest{Hash: txHash})
	if err != nil {
		return err
	}
	if tx.StatusCode(resp.StatusCode) != tx.Success {
		return fmt.Errorf("publish contract fail " + (resp.String()))
	}

	contractID = "Contract" + txHash
	return nil
}

func (t *gobangHandle) wait(i int, h string, id string) bool {
	g := rpcpb.GetContractStorageRequest{Id: contractID, Key: "games" + id, Field: "", ByLongestChain: true}
	for {
		v, err := iostSDK.GetContractStorage(&g)
		if err != nil {
			time.Sleep(3000 * time.Millisecond)
			continue
		}
		var f interface{}
		err = json.Unmarshal([]byte(v.Data), &f)
		if err != nil {
			ilog.Error(err)
		}
		if f == nil {
			continue
		}
		m := f.(map[string]interface{})
		if m["hash"] == h {
			if m["winner"] != nil {
				return true
			}
			return false
		}
		time.Sleep(3000 * time.Millisecond)
	}
}

func (t *gobangHandle) getGameID(h string) string {
	for {
		v, err := iostSDK.GetTxReceiptByTxHash(h)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if len(v.Returns) == 0 {
			continue
		}
		var f interface{}
		err = json.Unmarshal([]byte(v.Returns[0]), &f)
		if err != nil {
			ilog.Error(err)
		}
		if f == nil {
			continue
		}
		m := f.([]interface{})
		return m[0].(string)
	}
}

func (t *gobangHandle) Run(i int) (interface{}, error) {
	ilog.Info("run ", i)
	var board = make(map[string]bool)
	act := tx.NewAction(contractID, "newGameWith", fmt.Sprintf(`["%v"]`, testID))
	h := t.transfer(i, act, rootAcc, rootID)
	gameID := t.getGameID(h)
	t.wait(i, h, gameID)
	round := 0
	for {
		acc := rootAcc
		id := rootID
		var x, y int
		for {
			x = rand.Intn(15)
			y = rand.Intn(15)
			if board[strconv.Itoa(x)+","+strconv.Itoa(y)] == false {
				board[strconv.Itoa(x)+","+strconv.Itoa(y)] = true
				break
			}
		}
		if round%2 == 1 {
			acc = testAcc
			id = testID
		}
		act = tx.NewAction(contractID, "move", fmt.Sprintf(`[%v, %v, %v, "%v"]`, gameID, x, y, h))
		h = t.transfer(i, act, acc, id)
		r := t.wait(i, h, gameID)
		if r {
			break
		}
		round++
	}
	return "", nil
}

// Transfer ...
func (t *gobangHandle) transfer(i int, act *tx.Action, acc *account.KeyPair, id string) string {
	trx := tx.NewTx([]*tx.Action{act}, []string{}, 1000000000, 100, time.Now().Add(time.Second*time.Duration(10000)).UnixNano(), 0, chainID)
	stx, err := tx.SignTx(trx, id, []*account.KeyPair{acc})
	if err != nil {
		return fmt.Sprintf("signtx:%v err:%v", stx, err)
	}
	var txHash string
	txHash, err = call.SendTx(stx, i)
	if err != nil {
		return fmt.Sprintf("sendtx:%v  err:%v", txHash, err)
	}
	return txHash
}
