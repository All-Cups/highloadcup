package app_test

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

func BenchmarkBankToWallet(b *testing.B) {
	const (
		supply    = 10_000_000
		walletMax = 1_000
	)
	bank := make([]uint64, supply/8/8)
	for i := range bank {
		bank[i] = rand.Uint64()
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wallet := make([]int32, 0, walletMax)
	LOOP:
		for i := range bank {
			for bit := 0; bit < 64; bit++ {
				if bank[i]&(1<<bit) != 0 {
					wallet = wallet[:len(wallet)+1]
					wallet[len(wallet)-1] = int32(i*64 + bit)
					if len(wallet) == walletMax {
						break LOOP
					}
				}
			}
		}
		json.NewEncoder(ioutil.Discard).Encode(wallet)
	}
}
