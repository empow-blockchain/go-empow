package host

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/empow-blockchain/go-empow/ilog"
	"github.com/empow-blockchain/go-empow/vm/database"

	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/core/contract"
)

// Teller handler of iost
type Teller struct {
	h         *Host
	cost      map[string]contract.Cost
	cacheCost contract.Cost
}

// NewTeller new teller
func NewTeller(h *Host) Teller {
	return Teller{
		h:    h,
		cost: make(map[string]contract.Cost),
	}
}

// Costs ...
func (t *Teller) Costs() map[string]contract.Cost {
	return t.cost
}

// GasPaid ...
func (t *Teller) GasPaid(publishers ...string) int64 {
	var publisher string
	if len(publishers) > 0 {
		publisher = publishers[0]
	} else {
		publisher = t.h.Context().Value("publisher").(string)
	}
	v, ok := t.cost[publisher]
	if !ok {
		return 0
	}
	return v.ToGas()
}

// ClearCosts ...
func (t *Teller) ClearCosts() {
	t.cost = make(map[string]contract.Cost)
}

// ClearRAMCosts ...
func (t *Teller) ClearRAMCosts() {
	newCost := make(map[string]contract.Cost)
	for k, c := range t.cost {
		if c.Net != 0 || c.CPU != 0 {
			newCost[k] = contract.NewCost(0, c.Net, c.CPU)
		}
	}
	t.cost = newCost
}

// AddCacheCost ...
func (t *Teller) AddCacheCost(c contract.Cost) {
	t.cacheCost.AddAssign(c)
}

// CacheCost ...
func (t *Teller) CacheCost() contract.Cost {
	return t.cacheCost
}

// FlushCacheCost ...
func (t *Teller) FlushCacheCost() {
	t.PayCost(t.cacheCost, "")
	t.cacheCost = contract.Cost0()
}

// ClearCacheCost ...
func (t *Teller) ClearCacheCost() {
	t.cacheCost = contract.Cost0()
}

// PayCost ...
func (t *Teller) PayCost(c contract.Cost, who string) {
	//fmt.Printf("paycost [%v] %v(%v)\n", who, c, c.ToGas())
	costMap := make(map[string]contract.Cost)
	if c.CPU > 0 || c.Net > 0 {
		costMap[who] = contract.Cost{CPU: c.CPU, Net: c.Net}
	}
	for _, item := range c.DataList {
		if oc, ok := costMap[item.Payer]; ok {
			oc.AddAssign(contract.Cost{Data: item.Val, DataList: []contract.DataItem{item}})
			costMap[item.Payer] = oc
		} else {
			costMap[item.Payer] = contract.Cost{Data: item.Val, DataList: []contract.DataItem{item}}
		}
	}
	for who, c := range costMap {
		if oc, ok := t.cost[who]; ok {
			oc.AddAssign(c)
			t.cost[who] = oc
		} else {
			t.cost[who] = c
		}
	}
}

// IsProducer check account is producer
func (t *Teller) IsProducer(acc string) bool {
	pm := t.h.DB().Get("vote_producer.empow-producerMap")
	pmStr := database.Unmarshal(pm)
	if _, ok := pmStr.(error); ok {
		return false
	}
	producerMap, err := simplejson.NewJson([]byte(pmStr.(string)))
	if err != nil {
		return false
	}
	_, ok := producerMap.CheckGet(acc)
	return ok
}

// DoPay ...
func (t *Teller) DoPay(witness string, gasRatio int64) (paidGas *common.Fixed, err error) {
	for payer, costOfPayer := range t.cost {
		gas := &common.Fixed{
			Value:   gasRatio * costOfPayer.ToGas(),
			Decimal: database.GasDecimal,
		}
		if !gas.IsZero() {
			err := t.h.CostGas(payer, gas)
			if err != nil {
				return nil, fmt.Errorf("pay gas cost failed: %v %v %v", err, payer, gas)
			}

			var enableReferrerReward = false
			if enableReferrerReward {
				// reward 15% gas to account referrer
				if !t.h.IsContract(payer) {
					acc, _ := ReadAuth(t.h.DB(), payer)
					if acc == nil {
						ilog.Fatalf("invalid account %v", payer)
					}
					if acc.Referrer != "" && t.IsProducer(acc.Referrer) {
						reward := gas.TimesF(0.1)
						t.h.ChangeTGas(acc.Referrer, reward, true)
					}
				}
			}

		}

		if payer == t.h.Context().Value("publisher").(string) {
			paidGas = gas
		}
		// contracts in "em" domain will not pay for ram
		if !strings.HasSuffix(payer, ".empow") {
			var ramPayer string
			if t.h.IsContract(payer) {
				p, _ := t.h.GlobalMapGet("system.empow", "contract_owner", payer)
				var ok bool
				ramPayer, ok = p.(string)
				if !ok {
					ilog.Fatalf("DoPay failed: contract %v has no owner", payer)
				}
			} else {
				ramPayer = payer
			}

			ram := costOfPayer.Data
			currentRAM := t.h.db.TokenBalance("ram", ramPayer)
			if currentRAM < ram {
				err = fmt.Errorf("pay ram failed. id: %v need %v, actual %v", ramPayer, ram, currentRAM)
				return
			}
			t.h.db.SetTokenBalance("ram", ramPayer, currentRAM-ram)
			t.h.db.ChangeUsedRAMInfo(ramPayer, ram)
		}
	}
	return
}

// Privilege ...
func (t *Teller) Privilege(id string) int {
	am, ok := t.h.ctx.Value("auth_list").(map[string]int)
	if !ok {
		return 0
	}
	i, ok := am[id]
	if !ok {
		i = 0
	}
	return i
}
