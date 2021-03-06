package integration

import (
	"fmt"
	"testing"

	"github.com/empow-blockchain/go-empow/ilog"
	"github.com/empow-blockchain/go-empow/vm/database"

	"github.com/empow-blockchain/go-empow/core/tx"
	. "github.com/empow-blockchain/go-empow/verifier"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_IssueBonus(t *testing.T) {
	ilog.Stop()
	Convey("test bonus.empow", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareFakeBase(t, s)

		// deploy issue.empow
		setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
		s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)

		Convey("test issueContribute", func() {
			s.Head.Witness = acc1.KeyPair.ReadablePubkey()
			s.Head.Number = 1

			r, err := s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","12345678"]}]`, acc1.ID), acc1.ID, acc1.KeyPair)
			s.Visitor.Commit()

			So(err, ShouldBeNil)
			So(r.Status.Code, ShouldEqual, tx.Success)
			So(s.Visitor.TokenBalance("contribute", acc1.ID), ShouldEqual, int64(78217486))
		})
	})
}

func Test_ExchangeEMPOW(t *testing.T) {
	ilog.Stop()
	Convey("test bonus.empow", t, func() {
		s := NewSimulator()
		defer s.Clear()

		s.Head.Number = 0
		createAccountsWithResource(s)
		prepareIssue(s, acc0)
		prepareNewProducerVote(t, s, acc0)
		initProducer(t, s)
		prepareFakeBase(t, s)

		// deploy bonus.empow
		setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
		s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)

		Convey("test exchangeEMPOW", func() {
			createToken(t, s, acc0)

			// set bonus pool
			s.Call("token.empow", "issue", fmt.Sprintf(`["%v", "%v", "%v"]`, "em", "bonus.empow", "1000"), acc0.ID, acc0.KeyPair)

			// gain contribute
			s.Head.Witness = acc1.KeyPair.ReadablePubkey()
			s.Head.Number = 1
			r, err := s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc1.ID, 1), acc1.ID, acc1.KeyPair)
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			s.Visitor.Commit()

			So(s.Visitor.TokenBalance("contribute", acc1.ID), ShouldEqual, int64(78217486))

			s.Head.Witness = acc2.KeyPair.ReadablePubkey()
			s.Head.Number = 2
			r, err = s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc2.ID, 123456789), acc2.ID, acc2.KeyPair)
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			s.Visitor.Commit()

			So(s.Visitor.TokenBalance("contribute", acc2.ID), ShouldEqual, int64(78217486))

			r, err = s.Call("bonus.empow", "exchangeEMPOW", fmt.Sprintf(`["%v", "%v"]`, acc1.ID, "0.5"), acc1.ID, acc1.KeyPair)
			s.Visitor.Commit()

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(s.Visitor.TokenBalance("contribute", acc1.ID), ShouldEqual, int64(28217486))
			So(s.Visitor.TokenBalance("em", acc1.ID), ShouldEqual, int64(50000000))
			So(s.Visitor.TokenBalance("em", "bonus.empow"), ShouldEqual, int64(99950000000))
		})
	})
}

func Test_UpdateBonus(t *testing.T) {
	ilog.Stop()
	Convey("test update bonus", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareFakeBase(t, s)

		r, err := prepareIssue(s, acc0)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")

		prepareStake(t, s, acc0)
		prepareSocial(t, s, acc0)

		// deploy issue.empow
		err = setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
		So(err, ShouldBeNil)
		r, err = s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")

		prepareNewProducerVote(t, s, acc0)
		initProducer(t, s)

		Convey("test update bonus 1", func() {
			s.Head.Witness = acc1.KeyPair.ReadablePubkey()
			s.Head.Number = 1

			So(database.MustUnmarshal(s.Visitor.Get("bonus.empow-blockContrib")), ShouldEqual, `"0.78217486"`)

			for i := 0; i < 7; i++ {
				s.Head.Time += 86400 * 1e9

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				r, err = s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","12345678"]}]`, acc1.ID), acc1.ID, acc1.KeyPair)
				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
			}

			So(s.Visitor.TokenBalance("em", "bonus.empow"), ShouldEqual, int64(95921961516586))
			So(s.Visitor.TokenBalance("contribute", acc1.ID), ShouldEqual, int64(547582439))
			So(database.MustUnmarshal(s.Visitor.Get("bonus.empow-blockContrib")), ShouldEqual, `"0.78277523"`)

			for i := 0; i < 7; i++ {
				s.Head.Time += 86400 * 1e9

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				r, err = s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","12345678"]}]`, acc1.ID), acc1.ID, acc1.KeyPair)
				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
			}

			So(s.Visitor.TokenBalance("em", "bonus.empow"), ShouldEqual, int64(191917531200059))
			So(s.Visitor.TokenBalance("contribute", acc1.ID), ShouldEqual, int64(1095585169))
			So(database.MustUnmarshal(s.Visitor.Get("bonus.empow-blockContrib")), ShouldEqual, `"0.78337592"`)
		})
	})
}
