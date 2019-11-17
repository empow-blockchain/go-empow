package integration

import (
	"encoding/json"
	"testing"

	"github.com/empow-blockchain/go-empow/core/tx"

	"github.com/empow-blockchain/go-empow/ilog"

	"github.com/empow-blockchain/go-empow/common"
	. "github.com/empow-blockchain/go-empow/verifier"
	. "github.com/smartystreets/goconvey/convey"
)

func prepareIssue(s *Simulator, acc *TestAccount) (*tx.TxReceipt, error) {
	s.Head.Number = 0

	// deploy issue.empow
	setNonNativeContract(s, "issue.empow", "issue.js", ContractPath)
	s.Call("issue.empow", "init", `[]`, acc.ID, acc.KeyPair)

	witness := common.Witness{
		ID:      acc0.ID,
		Owner:   acc0.KeyPair.ReadablePubkey(),
		Active:  acc0.KeyPair.ReadablePubkey(),
		Balance: 55000000000,
	}
	params := []interface{}{
		acc0.ID,
		common.TokenInfo{
			FoundationAccount: acc1.ID,
			EMTotalSupply:     90000000000,
			EMDecimal:         8,
		},
		[]interface{}{witness},
	}
	b, _ := json.Marshal(params)
	r, err := s.Call("issue.empow", "initGenesis", string(b), acc.ID, acc.KeyPair)
	s.Visitor.Commit()
	return r, err
}

func Test_EMPOWIssue(t *testing.T) {
	ilog.Stop()
	Convey("test issue.empow", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		r, err := prepareIssue(s, acc0)

		Convey("test init", func() {
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(s.Visitor.TokenBalance("em", acc0.ID), ShouldEqual, int64(210*1e16))
		})

		prepareNewProducerVote(t, s, acc0)
		initProducer(t, s)

		Convey("test issueEM", func() {
			s.Head.Time += 4 * 3 * 1e9
			r, err := s.Call("issue.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
			s.Visitor.Commit()

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("em", "bonus.empow"), ShouldEqual, int64(7884322975))
			So(s.Visitor.TokenBalance("em", acc1.ID), ShouldEqual, int64(7884323211))
		})
	})
}
