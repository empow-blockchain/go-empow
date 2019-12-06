package rpc

import (
	"fmt"
	"strings"
)

func checkHashValid(hash string) error {
	maxLen := 100
	if len(hash) >= maxLen {
		return fmt.Errorf("Hash len(%v) is greater then %v", len(hash), maxLen)
	}
	return nil
}

func checkIDValid(id string) error {
	if strings.HasPrefix(id, "Contract") || strings.HasSuffix(id, ".empow") {
		if len(id) >= 100 {
			return fmt.Errorf("address invalid. ContractID length should be less then 100 - %v", len(id))
		}
		return nil
	}

	if len(id) != 49 {
		return fmt.Errorf("wrong address")
	}

	for _, v := range id {
		if !((v >= 'A' && v <= 'z') || (v >= '0' && v <= '9' || v == '_')) {
			return fmt.Errorf("id invalid. id contains invalid character - %v", v)
		}
	}
	return nil
}
