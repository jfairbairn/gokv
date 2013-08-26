package gokv

import (
	"fmt"
	"regexp"
	"math/big"
)

type IdGen struct {
	ids map[string]uint64
}

func NewIdGen() *IdGen {
	return &IdGen{ids: make(map[string]uint64)}
}

func (idgen *IdGen) NewId(prefix string) string {
	id := idgen.ids[prefix]+1
	idgen.ids[prefix] = id
	return fmt.Sprintf("%s.%d", prefix, id)
}

var keyRegex *regexp.Regexp

func init() {
	keyRegex = regexp.MustCompile("^(.*)\\.(\\d+)$") 
}

func (idgen *IdGen) OnKey(key string) {
	results := keyRegex.FindAllStringSubmatch(key, 1)
	if len(results) != 1 {
		return
	}
	if len(results[0]) != 3 {
		return
	}
	prefix := results[0][1]
	
	big := big.NewInt(0)
	big.SetString(results[0][2], 10)

	id := uint64(big.Int64())
	oldid := idgen.ids[prefix]
	if id > oldid {
		idgen.ids[prefix] = id
	}
}