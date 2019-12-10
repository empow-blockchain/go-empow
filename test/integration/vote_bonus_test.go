package integration

import (
	"fmt"
	"testing"

	"github.com/empow-blockchain/go-empow/core/tx"
	"github.com/empow-blockchain/go-empow/ilog"
	. "github.com/empow-blockchain/go-empow/verifier"
	"github.com/empow-blockchain/go-empow/vm/database"
	"github.com/stretchr/testify/assert"
)

func Test_VoteBonus(t *testing.T) {
	ilog.Stop()
	s := NewSimulator()
	defer s.Clear()

	s.Head.Number = 0

	createAccountsWithResource(s)
	prepareFakeBase(t, s)
	prepareIssue(s, acc0)
	prepareStake(t, s, acc0)
	prepareSocial(t, s, acc0)
	prepareNewProducerVote(t, s, acc0)
	initProducer(t, s)

	// deploy bonus.empow
	setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
	s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)

	s.Head.Number = 1
	for _, acc := range testAccounts[6:] {
		r, err := s.Call("vote_producer.empow", "applyRegister", fmt.Sprintf(`["%v", "%v", "loc", "url", "netId", true]`, acc.ID, acc.KeyPair.ReadablePubkey()), acc.ID, acc.KeyPair)
		assert.Nil(t, err)
		assert.Equal(t, tx.Success, r.Status.Code)
		r, err = s.Call("vote_producer.empow", "approveRegister", fmt.Sprintf(`["%v"]`, acc.ID), acc0.ID, acc0.KeyPair)
		assert.Nil(t, err)
		assert.Equal(t, tx.Success, r.Status.Code)
		r, err = s.Call("vote_producer.empow", "logInProducer", fmt.Sprintf(`["%v"]`, acc.ID), acc.ID, acc.KeyPair)
		assert.Nil(t, err)
		assert.Equal(t, tx.Success, r.Status.Code)
	}

	s.Visitor.SetTokenBalance("vote", acc0.ID, 1e17)
	s.Visitor.SetTokenBalance("vote", acc2.ID, 1e17)

	for idx, acc := range testAccounts {
		voter := acc0
		if idx > 0 {
			r, err := s.Call("vote_producer.empow", "vote", fmt.Sprintf(`["%v", "%v", "%v"]`, voter.ID, acc.ID, idx*2e7), voter.ID, voter.KeyPair)
			assert.Nil(t, err)
			assert.Empty(t, r.Status.Message)
			assert.Equal(t, fmt.Sprintf(`{"votes":"%d","deleted":0,"clearTime":-1}`, idx*2e7), database.MustUnmarshal(s.Visitor.MGet("vote.empow-v_1", acc.ID)))
		}
		voter = acc2
		r, err := s.Call("vote_producer.empow", "vote", fmt.Sprintf(`["%v", "%v", "%v"]`, voter.ID, acc.ID, 2e7), voter.ID, voter.KeyPair)
		assert.Nil(t, err)
		assert.Empty(t, r.Status.Message)
		assert.Equal(t, fmt.Sprintf(`{"votes":"%d","deleted":0,"clearTime":-1}`, (idx+1)*2e7), database.MustUnmarshal(s.Visitor.MGet("vote.empow-v_1", acc.ID)))
	}

	for idx, acc := range testAccounts {
		s.Head.Witness = acc.KeyPair.ReadablePubkey()
		for i := 0; i <= idx; i++ {
			s.Head.Number++
			r, err := s.Call("base.empow", "issueContribute", fmt.Sprintf(`[{"parent":["%v","%v"]}]`, acc.ID, 1), acc.ID, acc.KeyPair)
			assert.Nil(t, err)
			assert.Empty(t, r.Status.Message)
		}
		assert.Equal(t, int64(78217486*(idx+1)), s.Visitor.TokenBalance("contribute", acc.ID))
	}
	assert.Equal(t, `{"user_1":["20000000",1,"0"],"user_2":["40000000",1,"0"],"user_3":["60000000",1,"0"],"user_4":["80000000",1,"0"],"user_5":["100000000",1,"0"],"user_6":["120000000",1,"0"],"user_7":["140000000",1,"0"],"user_8":["160000000",1,"0"],"user_9":["180000000",1,"0"]}`, database.MustUnmarshal(s.Visitor.MGet("vote.empow-u_1", acc0.ID)))
	assert.Equal(t, `{"user_0":["20000000",1,"0"],"user_1":["20000000",1,"0"],"user_2":["20000000",1,"0"],"user_3":["20000000",1,"0"],"user_4":["20000000",1,"0"],"user_5":["20000000",1,"0"],"user_6":["20000000",1,"0"],"user_7":["20000000",1,"0"],"user_8":["20000000",1,"0"],"user_9":["20000000",1,"0"]}`, database.MustUnmarshal(s.Visitor.MGet("vote.empow-u_1", acc2.ID)))
	s.Head.Time += 5073358980

	r, err := s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(100402187976), s.Visitor.TokenBalance("em", "vote_producer.empow"))
	assert.Equal(t, int64(804375953), s.Visitor.TokenBalance("em", "bonus.empow"))

	for i := 0; i < 10; i++ {
		s.Visitor.SetTokenBalance("em", testAccounts[i].ID, 0)
	}

	// 0. normal withdraw
	r, err = s.Call("bonus.empow", "exchangeEMPOW", fmt.Sprintf(`["%s","0"]`, acc1.ID), acc1.ID, acc1.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(0), s.Visitor.TokenBalance("contribute", acc1.ID))
	assert.Equal(t, int64(156434972), s.Visitor.TokenBalance("em", acc1.ID))
	assert.Equal(t, int64(100402187976), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	s.Head.Time += 86400 * 1e9
	r, err = s.Call("bonus.empow", "exchangeEMPOW", fmt.Sprintf(`["%s","%s"]`, acc1.ID, "0.00000001"), acc1.ID, acc1.KeyPair)
	assert.Nil(t, err)
	assert.Contains(t, r.Status.Message, "invalid amount: negative or greater than contribute")

	// withdraw by admin
	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc1.ID), acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(156434972+14625017), s.Visitor.TokenBalance("em", acc1.ID)) // 14625017 = (402187976*(2/55))
	assert.Equal(t, int64(100402187976-14625017), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc1.ID), acc1.ID, acc1.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(156434972+14625017), s.Visitor.TokenBalance("em", acc1.ID)) // not change
	assert.Equal(t, int64(100402187976-14625017), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	// 1. unregistered withdraw
	r, err = s.Call("vote_producer.empow", "forceUnregister", fmt.Sprintf(`["%v"]`, acc3.ID), acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	r, err = s.Call("vote_producer.empow", "unregister", fmt.Sprintf(`["%v"]`, acc3.ID), acc3.ID, acc3.KeyPair)
	assert.Nil(t, err)
	assert.Contains(t, r.Status.Message, "producer in pending list or in current list, can't unregister")

	r, err = s.Call("bonus.empow", "exchangeEMPOW", fmt.Sprintf(`["%s","0"]`, acc3.ID), acc3.ID, acc3.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(0), s.Visitor.TokenBalance("contribute", acc3.ID))
	assert.Equal(t, int64(312869944), s.Visitor.TokenBalance("em", acc3.ID))
	assert.Equal(t, int64(100387562959), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	s.Head.Time += 86400 * 1e9
	r, err = s.Call("bonus.empow", "exchangeEMPOW", fmt.Sprintf(`["%s","%s"]`, acc3.ID, "0.00000001"), acc3.ID, acc3.KeyPair)
	assert.Nil(t, err)
	assert.Contains(t, r.Status.Message, "invalid amount: negative or greater than contribute")

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc3.ID), acc3.ID, acc3.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(312869944+29250034), s.Visitor.TokenBalance("em", acc3.ID)) // 29250034 = (402187976*(4/55))
	assert.Equal(t, int64(100387562959-29250034), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc3.ID), acc3.ID, acc3.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(312869944+29250034), s.Visitor.TokenBalance("em", acc3.ID)) // not change
	assert.Equal(t, int64(100387562959-29250034), s.Visitor.TokenBalance("em", "vote_producer.empow"))
}

func TestCriticalVoteCase(t *testing.T) {
	ilog.Stop()
	s := NewSimulator()
	defer s.Clear()

	s.Head.Number = 0

	createAccountsWithResource(s)
	prepareFakeBase(t, s)
	prepareIssue(s, acc0)
	prepareStake(t, s, acc0)
	prepareSocial(t, s, acc0)
	prepareNewProducerVote(t, s, acc0)
	initProducer(t, s)

	// deploy bonus.empow
	setNonNativeContract(s, "bonus.empow", "bonus.js", ContractPath)
	s.Call("bonus.empow", "init", `[]`, acc0.ID, acc0.KeyPair)
	s.Head.Number = 1

	s.Visitor.SetTokenBalance("vote", acc0.ID, 1e17)

	for _, acc := range testAccounts[6:] {
		r, err := s.Call("vote_producer.empow", "applyRegister", fmt.Sprintf(`["%v", "%v", "loc", "url", "netId", true]`, acc.ID, acc.KeyPair.ReadablePubkey()), acc.ID, acc.KeyPair)
		assert.Nil(t, err)
		assert.Equal(t, tx.Success, r.Status.Code)
	}
	for idx, acc := range testAccounts {
		voter := acc0
		r, err := s.Call("vote_producer.empow", "vote", fmt.Sprintf(`["%v", "%v", "%v"]`, voter.ID, acc.ID, (idx+1)*4e5), voter.ID, voter.KeyPair)
		assert.Nil(t, err)
		assert.Empty(t, r.Status.Message)
		assert.Equal(t, fmt.Sprintf(`{"votes":"%d","deleted":0,"clearTime":-1}`, (idx+1)*4e5), database.MustUnmarshal(s.Visitor.MGet("vote.empow-v_1", acc.ID)))
		if idx == 0 {
			assert.Nil(t, database.MustUnmarshal(s.Visitor.Get("vote_producer.empow-candAllKey")))
		} else {
			assert.Equal(t, fmt.Sprintf(`"%d"`, (idx+1)*(idx+2)*4e5/2-4e5), database.MustUnmarshal(s.Visitor.Get("vote_producer.empow-candAllKey")))
		}
	}
	assert.Equal(t, `{"user_0":["400000",1,"0"],"user_1":["800000",1,"0"],"user_2":["1200000",1,"0"],"user_3":["1600000",1,"0"],"user_4":["2000000",1,"0"],"user_5":["2400000",1,"0"],"user_6":["2800000",1,"0"],"user_7":["3200000",1,"0"],"user_8":["3600000",1,"0"],"user_9":["4000000",1,"0"]}`, database.MustUnmarshal(s.Visitor.MGet("vote.empow-u_1", acc0.ID)))
	s.Head.Time += 5073358980
	r, err := s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(100402187976), s.Visitor.TokenBalance("em", "vote_producer.empow"))
	assert.Equal(t, int64(804375953), s.Visitor.TokenBalance("em", "bonus.empow"))

	for i := 0; i < 10; i++ {
		s.Visitor.SetTokenBalance("em", testAccounts[i].ID, 0)
	}
	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc6.ID), acc6.ID, acc6.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(52135478), s.Visitor.TokenBalance("em", acc6.ID)) // 52135478 = (402187976*(7/54))
	assert.Equal(t, int64(100402187976-52135478), s.Visitor.TokenBalance("em", "vote_producer.empow"))

	s.Visitor.SetTokenBalance("vote", acc2.ID, 1e17)
	r, err = s.Call("vote_producer.empow", "vote", fmt.Sprintf(`["%v", "%v", "%v"]`, acc2.ID, acc0.ID, 1e5), acc2.ID, acc2.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, `"22100000"`, database.MustUnmarshal(s.Visitor.Get("vote_producer.empow-candAllKey")))

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(0), s.Visitor.TokenBalance("em", acc0.ID)) // not changed

	s.Head.Time += 24*3600*1e9 + 1
	r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(6949666534929), s.Visitor.TokenBalance("em", "vote_producer.empow"))
	assert.Equal(t, int64(13699437340816), s.Visitor.TokenBalance("em", "bonus.empow"))

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(154961911367), s.Visitor.TokenBalance("em", acc0.ID))

	r, err = s.Call("vote_producer.empow", "unvote", fmt.Sprintf(`["%v", "%v", "%v"]`, acc2.ID, acc0.ID, 1e5), acc2.ID, acc2.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, `"21600000"`, database.MustUnmarshal(s.Visitor.Get("vote_producer.empow-candAllKey")))

	s.Head.Time += 24*3600*1e9 + 1
	r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)

	r, err = s.Call("vote_producer.empow", "candidateWithdraw", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	assert.Nil(t, err)
	assert.Empty(t, r.Status.Message)
	assert.Equal(t, int64(154961911367), s.Visitor.TokenBalance("em", acc0.ID)) // not changed
}
