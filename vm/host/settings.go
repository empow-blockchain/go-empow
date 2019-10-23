package host

import "github.com/empow-blockchain/go-empow/core/contract"

// Setting in state db
type Setting struct {
	Costs map[string]contract.Cost `json:"costs"`
}
