package main

import (
	"fmt"
	"strings"

	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/db/kv/leveldb"
	"github.com/empow-blockchain/go-empow/vm/database"
)

func padTo(s string, ptn string, l int) string {
	if len(s) < l {
		return s + strings.Repeat(ptn, l-len(s))
	}
	return s
}

func printTokenBalance(db *leveldb.DB, tokenType string) {
	fmt.Println("############# ", tokenType, " balance ##############")
	prefix := "state/m-token.empow-TB"
	keys, err := db.Keys([]byte(prefix))
	if err != nil {
		panic(err)
	}
	suffix := "-" + tokenType
	decimalKey := "state/m-token.empow-TI" + tokenType + "-decimal"
	decimalRaw, err := db.Get([]byte(decimalKey))
	if err != nil {
		panic(err)
	}
	decimal := database.MustUnmarshal(string(decimalRaw))
	for _, k := range keys {
		if !strings.HasSuffix(string(k), suffix) {
			continue
		}
		rawValue, err := db.Get(k)
		if err != nil {
			panic(err)
		}
		v := database.MustUnmarshal(string(rawValue))
		f := common.Fixed{Value: v.(int64), Decimal: int(decimal.(int64))}
		tmp := string(k)[len(prefix):]
		user := tmp[:len(tmp)-len(suffix)]
		user = padTo(user, " ", 20)
		fmt.Printf("%v\t%v\n", user, f.ToString())
	}
	fmt.Println()
}

func printAll(db *leveldb.DB) { // nolint
	fmt.Println("######## all kvs #############")
	iter := db.NewIteratorByPrefix([]byte("state/")).(*leveldb.Iter)
	for iter.Next() {
		k := string(iter.Key())
		v := string(iter.Value())
		if len(v) > 100 {
			v = v[:100] + "..."
		}
		fmt.Printf("%v\t%v\n", k, v)
	}
}

func printRAMUsage(db *leveldb.DB) {
	fmt.Println("######## system ram usage #############")
	m := make(map[string]int)
	iter := db.NewIteratorByPrefix([]byte("state/")).(*leveldb.Iter)
	for iter.Next() {
		k := string(iter.Key())
		v := string(iter.Value())
		var owner string
		var ramUse int
		if strings.HasPrefix(k, "state/m-") && strings.HasPrefix(v, "@") {
			// map
			continue
		}
		if strings.HasPrefix(k, "state/c-") {
			cid := k[len("state/c-"):]
			var err error
			ownerRaw, err := db.Get([]byte("state/m-system.empow-contract_owner-" + cid))
			if err != nil {
				panic(err)
			}
			owner = string(ownerRaw)
			if owner == "" {
				if !strings.HasSuffix(cid, "em") {
					panic("non iost contracts should have owner")
				}
				//owner = "OWNER_" + cid
				owner = "[iost codes]"
			} else {
				idx := strings.LastIndex(owner, "@")
				if idx == -1 {
					panic("invalid contract owner format " + k + ":" + v + " " + owner)
				}
				owner = owner[:idx]
			}
			ramUse = len(v)
		} else {
			idx := strings.LastIndex(v, "@")
			if idx == -1 {
				panic("empty owner " + k + ":" + v)
			}
			owner = v[(idx + 1):]
			ramUse = idx
			//if !strings.HasPrefix(k, "state/b-base.empow-chain_info_") {
			//	fmt.Printf("%v\t%v\t%v\n", owner, k, v)
			//}
		}
		if owner == "" {
			owner = "[unknown]"
			//fmt.Printf("WHY!! %v %v %v\n", owner, k, v)
		}
		owner = padTo(owner, " ", 20)
		old, ok := m[owner]
		if ok {
			m[owner] = old + ramUse
		} else {
			m[owner] = ramUse
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}
	for k := range m {
		fmt.Printf("%v\t%v\n", k, m[k])
	}
	fmt.Println()
}

func main() {
	storagePath := "storage/StateDB"
	db, err := leveldb.NewDB(storagePath)
	defer func() {
		db.Close()
	}()
	if err != nil {
		panic(err)
	}
	printRAMUsage(db)
	printTokenBalance(db, "em")
	printTokenBalance(db, "ram")
	//printAll(db)
}
