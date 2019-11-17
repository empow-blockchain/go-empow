package integration

import (
	"fmt"
	"testing"

	"github.com/empow-blockchain/go-empow/core/tx"
	. "github.com/empow-blockchain/go-empow/verifier"
)

func prepareStake(t *testing.T, s *Simulator, acc *TestAccount) {
	s.Head.Number = 0

	// deploy issue.empow
	setNonNativeContract(s, "stake.empow", "stake.js", ContractPath)

	r, err := s.Call("stake.empow", "init", `[]`, acc1.ID, acc1.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}

	r, err = s.Call("stake.empow", "initAdmin", fmt.Sprintf(`["%s"]`, acc1.ID), acc1.ID, acc1.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}
}
