package integration

import (
	"fmt"
	"testing"

	"github.com/empow-blockchain/go-empow/core/tx"
	"github.com/empow-blockchain/go-empow/ilog"
	. "github.com/empow-blockchain/go-empow/verifier"
	"github.com/empow-blockchain/go-empow/vm/database"
	. "github.com/smartystreets/goconvey/convey"
)

func prepareBase(t *testing.T, s *Simulator, acc *TestAccount) {
	// deploy base.empow
	setNonNativeContract(s, "base.empow", "base.js", ContractPath)
	r, err := s.Call("base.empow", "init", `[]`, acc.ID, acc.KeyPair)
	So(err, ShouldBeNil)
	So(r.Status.Code, ShouldEqual, tx.Success)
	s.Visitor.Commit()
}

func Test_Base(t *testing.T) {
	ilog.Stop()
	Convey("test Base", t, func() {
		s := NewSimulator()
		defer s.Clear()

		s.Head.Number = 0

		createAccountsWithResource(s)
		prepareToken(t, s, acc0)
		prepareNewProducerVote(t, s, acc0)
		for _, acc := range testAccounts[:6] {
			r, err := s.Call("vote_producer.empow", "initProducer", fmt.Sprintf(`["%v", "%v"]`, acc.ID, acc.KeyPair.ReadablePubkey()), acc.ID, acc.KeyPair)
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
		}

		// deploy bonus.empow
		setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
		r, err := s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")

		prepareBase(t, s, acc0)

		s.Head.Number = 200
		s.Head.Witness = "test_witness_01"
		re, err := s.Call("base.empow", "exec", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc0.ID, 12345678), acc0.ID, acc0.KeyPair)
		So(err, ShouldBeNil)
		So(re.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get("base.empow-witness_produced")), ShouldEqual, `{"test_witness_01":1}`)

		s.Head.Number++
		s.Head.Witness = "test_witness_02"
		re, err = s.Call("base.empow", "exec", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc0.ID, 12345678), acc0.ID, acc0.KeyPair)
		So(err, ShouldBeNil)
		So(re.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get("base.empow-witness_produced")), ShouldEqual, `{"test_witness_01":1,"test_witness_02":1}`)

		s.Head.Number++
		s.Head.Witness = "test_witness_02"
		re, err = s.Call("base.empow", "exec", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc0.ID, 12345678), acc0.ID, acc0.KeyPair)
		So(err, ShouldBeNil)
		So(re.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get("base.empow-witness_produced")), ShouldEqual, `{"test_witness_01":1,"test_witness_02":2}`)
	})
}
