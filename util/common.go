package util

import (
	"encoding/json"
	"math/rand"
)

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func JsonStr(i interface{}) string {
	bs, _ := json.Marshal(i)
	return string(bs)
}

func RandomInt64() int64 {
	return rand.Int63()
}
