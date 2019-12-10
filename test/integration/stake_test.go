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

func prepareStake(t *testing.T, s *Simulator, acc *TestAccount) {
	s.Head.Number = 0

	// deploy issue.empow
	setNonNativeContract(s, "stake.empow", "stake.js", ContractPath)

	r, err := s.Call("stake.empow", "init", `[]`, acc0.ID, acc0.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}

	r, err = s.Call("stake.empow", "initAdmin", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}
}

func Test_Stake(t *testing.T) {
	ilog.Stop()
	Convey("test stake", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareIssue(s, acc0)
		prepareStake(t, s, acc0)
		prepareSocial(t, s, acc0)
		prepareFakeBase(t, s)

		s.Visitor.SetTokenBalance("em", acc2.ID, 1000*1e8)

		r, err := s.Call("stake.empow", "stake", fmt.Sprintf(`["%v","%v"]`, acc2.ID, 100), acc2.ID, acc2.KeyPair)

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		s.Visitor.Commit()

		So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 900*1e8)
		So(s.Visitor.TokenBalance("em", "stake.empow"), ShouldEqual, 100*1e8)
		So(database.MustUnmarshal(s.Visitor.Get("stake.empow-totalStakeAmount")), ShouldEqual, "100")
		So(database.MustUnmarshal(s.Visitor.Get("stake.empow-c_user_2")), ShouldEqual, "1")
		So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":0,"unstake":false,"amount":"100"}`)

		Convey("test multi package", func() {
			r, err = s.Call("stake.empow", "stake", fmt.Sprintf(`["%v", "%v"]`, acc2.ID, 200), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-c_user_2")), ShouldEqual, "2")

			r, err = s.Call("stake.empow", "stake", fmt.Sprintf(`["%v", "%v"]`, acc2.ID, 300), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-c_user_2")), ShouldEqual, "3")

			Convey("test withdraw all 2", func() {
				prepareNewProducerVote(t, s, acc0)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				currentBalance := s.Visitor.TokenBalance("em", acc2.ID)

				r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 1), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"200"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+16666666)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":345600,"unstake":false,"amount":"100"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":345600,"unstake":false,"amount":"200"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_2")), ShouldEqual, `{"lastBlockWithdraw":345600,"unstake":false,"amount":"300"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 40099999999) // currentBalance+(8333333+16666666+25000000)*2
			})

			Convey("test withdraw all 3", func() {
				prepareNewProducerVote(t, s, acc0)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				currentBalance := s.Visitor.TokenBalance("em", acc2.ID)

				r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 1), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"200"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+16666666)

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"100"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"200"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_2")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"300"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+8333333+16666666+25000000)
			})

			Convey("test withdraw all 4", func() {
				prepareNewProducerVote(t, s, acc0)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				currentBalance := s.Visitor.TokenBalance("em", acc2.ID)

				r, err = s.Call("stake.empow", "unstake", fmt.Sprintf(`["%v", %v]`, acc2.ID, 1), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":true,"amount":"200"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+16666666)

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"100"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":true,"amount":"200"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_2")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"300"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+8333333+16666666+25000000)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				currentBalance = s.Visitor.TokenBalance("em", acc2.ID)

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":345600,"unstake":false,"amount":"100"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":true,"amount":"200"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_2")), ShouldEqual, `{"lastBlockWithdraw":345600,"unstake":false,"amount":"300"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+8333333+25000000)
			})

			Convey("test withdraw all", func() {
				prepareNewProducerVote(t, s, acc0)

				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				currentBalance := s.Visitor.TokenBalance("em", acc2.ID)

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"100"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_1")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"200"}`)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_2")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"300"}`)
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, currentBalance+8333333+16666666+25000000)

				r, err = s.Call("stake.empow", "withdrawAll", fmt.Sprintf(`["%v"]`, acc2.ID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldContainSubstring, "All package can't withdraw")
			})
		})

		Convey("test topup", func() {

			prepareNewProducerVote(t, s, acc0)

			s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
			s.Head.Number += 24 * 60 * 60 * 2 // +1 day

			r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("em", "stake.empow"), ShouldEqual, int64(6859316438356))
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-i_1")), ShouldEqual, "684.93164383")

			currentBalance := s.Visitor.TokenBalance("em", acc2.ID)

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, int64(currentBalance+8333333)) // maximum 0.083% * 100
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":172800,"unstake":false,"amount":"100"}`)
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-restAmount")), ShouldEqual, "68493.0810496667")

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "package withdraw less than 1 day")
		})

		Convey("test withdraw within 1 year", func() {
			prepareNewProducerVote(t, s, acc0)

			for i := 0; i < 365; i++ {
				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
			}

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, int64(93041666665)) // maximum 0.083% * 100
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":63072000,"unstake":false,"amount":"100"}`)

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "package withdraw less than 1 day")
		})

		Convey("test withdraw within 7 days", func() {
			prepareNewProducerVote(t, s, acc0)

			for i := 0; i < 7; i++ {
				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
			}

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, int64(90058333333)) // maximum 0.083% * 100

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "package withdraw less than 1 day")
		})

		Convey("test withdraw each day within 7 days", func() {
			prepareNewProducerVote(t, s, acc0)

			for i := 0; i < 7; i++ {
				s.Head.Time += 24 * 60 * 60 * 1e9 // +1 day
				s.Head.Number += 24 * 60 * 60 * 2 // +1 day

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				lastBlockWithdraw := 172800 * (i + 1)
				So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, fmt.Sprintf(`{"lastBlockWithdraw":%v,"unstake":false,"amount":"100"}`, lastBlockWithdraw))
			}

			So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, int64(90058333331)) // maximum 0.083% * 100

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "package withdraw less than 1 day")
		})

		Convey("test withdraw less than 1 day", func() {
			s.Head.Time += 5 * 60 * 60 * 1e9 // +5 hours
			s.Head.Number += 5 * 60 * 60 * 2 // +5 hours

			r, err = s.Call("stake.empow", "withdraw", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "package withdraw less than 1 day")
		})

		Convey("test unstake", func() {
			r, err = s.Call("stake.empow", "unstake", fmt.Sprintf(`["%v", %v]`, acc2.ID, 0), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-p_user_2_0")), ShouldEqual, `{"lastBlockWithdraw":0,"unstake":true,"amount":"100"}`)
			So(database.MustUnmarshal(s.Visitor.Get("stake.empow-totalStakeAmount")), ShouldEqual, "0")

			r, err = s.Call("token.empow", "transfer", fmt.Sprintf(`["em", "%v", "%v", "%v", ""]`, acc2.ID, acc0.ID, 1000), acc2.ID, acc2.KeyPair)

			So(s.Visitor.FreezedTokenBalance("em", acc2.ID), ShouldEqual, 100*1e8)
			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "balance not enough")

			s.Head.Time += 3 * 24 * 60 * 60 * 1e9

			r, err = s.Call("token.empow", "transfer", fmt.Sprintf(`["em", "%v", "%v", "%v", ""]`, acc2.ID, acc0.ID, 1000), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
		})
	})
}
