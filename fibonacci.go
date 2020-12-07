package main

import (
	"math/big"
)

type fibStore struct {
	index    int
	sequence map[int]*big.Int
}

func intializeFibStore() *fibStore {
	f := new(fibStore)
	f.index = 0

	f.sequence = make(map[int]*big.Int)
	f.sequence[0] = big.NewInt(0)
	f.sequence[1] = big.NewInt(1)

	return f
}

func (f *fibStore) buildSequenceToIndex(recoveredIndex int) {
	for i := 2; i <= recoveredIndex; i++ {
		sum := new(big.Int)
		sum.Add(f.sequence[i-1], f.sequence[i-2])

		f.sequence[i] = sum
	}

	f.index = recoveredIndex
}

func (f *fibStore) getNext() *big.Int {
	current := f.index
	current += 1

	f.index = current

	if current == 1 {
		return f.sequence[current]
	}

	a := f.sequence[current-1]
	b := f.sequence[current-2]

	sum := new(big.Int)
	sum.Add(a, b)

	f.sequence[current] = sum
	return sum
}

func (f *fibStore) getCurrent() *big.Int {
	return f.sequence[f.index]
}

func (f *fibStore) getPrevious() *big.Int {
	current := f.index

	if current == 0 {
		return f.sequence[current]
	}

	current -= 1
	f.index = current
	return f.sequence[current]
}
