package game

import (
	"fmt"
	"sync"

	"github.com/powerman/structlog"
)

type bank struct {
	mu         sync.Mutex
	balance    int
	nextCoin   int
	coinIssued []bool
}

func newBank(ctx Ctx, totalCash int) *bank {
	log := structlog.FromContext(ctx, nil)
	log.Debug("newBank", "totalCash", totalCash)
	return &bank{
		coinIssued: make([]bool, totalCash),
	}
}

func (b *bank) getBalance() (balance int, wallet []int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	walletSize := maxWalletSize
	if b.balance < walletSize {
		walletSize = b.balance
	}
	wallet = make([]int, walletSize)

	// Search for issued coins by looking back from last issued coin
	// because there is higher chance last issued coins wasn't spent yet.
	next := b.nextCoin - 1
	if next < 0 {
		next = len(b.coinIssued) - 1
	}
	iter := newRR(next, false, len(b.coinIssued))
	for i := range wallet {
		for !b.coinIssued[next] {
			next = iter.next()
		}
		wallet[i] = next
		next = iter.next()
	}

	return b.balance, wallet
}

func (b *bank) earn(amount int) (wallet []int, _ error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !(amount >= 1 && amount+b.balance > b.balance && amount+b.balance <= len(b.coinIssued)) {
		return nil, fmt.Errorf("%w: %d (balance=%d, overall coins=%d)", errWrongAmount, amount, b.balance, len(b.coinIssued))
	}

	wallet = make([]int, amount)
	iter := newRR(b.nextCoin, true, len(b.coinIssued))
	for i := range wallet {
		for b.coinIssued[b.nextCoin] {
			b.nextCoin = iter.next()
		}
		b.coinIssued[b.nextCoin] = true
		b.balance++
		wallet[i] = b.nextCoin
		b.nextCoin = iter.next()
	}
	return wallet, nil
}

func (b *bank) spend(wallet []int) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(wallet) > b.balance {
		return fmt.Errorf("%w: too many coins in the wallet", ErrBogusCoin)
	}
	for i := range wallet {
		coin := wallet[i]
		switch {
		case coin < 0 || coin >= len(b.coinIssued):
			err = ErrBogusCoin
		case !b.coinIssued[coin]:
			err = fmt.Errorf("%w: %d", ErrBogusCoin, coin)
		}
		if err != nil {
			for j := 0; j < i; j++ {
				b.coinIssued[wallet[j]] = true
			}
			return err
		}
		b.coinIssued[coin] = false
	}
	b.balance -= len(wallet)
	return nil
}
