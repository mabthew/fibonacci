package main

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/golang/groupcache/lru"
)

var mutex = &sync.Mutex{}

type fibStore struct {
	index int
	cache *lru.Cache
}

var CacheMiss = errors.New("Cache miss")

func intializeCache(size int) (*fibStore, error) {
	fib := new(fibStore)
	lruCache := lru.New(size)

	fib.cache = lruCache
	fib.addToCache(0, big.NewInt(0))
	fib.addToCache(1, big.NewInt(1))

	return &fibStore{0, lruCache}, nil
}

func (f *fibStore) getFromCache(idx int) (*big.Int, error) {
	mutex.Lock()
	result, ok := f.cache.Get(idx)
	mutex.Unlock()
	if ok == false {
		return big.NewInt(0), CacheMiss
	}
	return result.(*big.Int), nil
}

func (f *fibStore) addToCache(idx int, value *big.Int) {
	mutex.Lock()
	f.cache.Add(idx, value)
	mutex.Unlock()
}

func (f *fibStore) buildSequenceToIndex(recoveredIndex int) {
	for i := 2; i <= recoveredIndex; i++ {
		sum := new(big.Int)
		a, err := f.getFromCache(i - 1)
		if err != nil {
			panic(err)
		}

		b, err := f.getFromCache(i - 2)
		if err != nil {
			panic(err)
		}

		sum.Add(a, b)

		f.addToCache(i, sum)
	}

	f.index = recoveredIndex
}

func (f *fibStore) getNext() *big.Int {
	current := f.index
	current += 1
	fmt.Println("current:", current)

	f.index = current

	if current == 1 {
		result, err := f.getFromCache(1)
		if err != nil {
			panic(err)
		}

		return result
	}

	result, err := f.getFromCache(current)
	if err != nil {
		if err == CacheMiss {
			// cache miss
			a, err := f.getFromCache(current - 1)
			if err != nil {
				panic(err)
			}

			b, err := f.getFromCache(current - 2)
			if err != nil {

				panic(err)
			}

			sum := new(big.Int)
			sum.Add(a, b)

			f.addToCache(current, sum)

			return sum
		} else {
			panic(err)
		}
	}

	// cache hit
	return result

}

func (f *fibStore) getCurrent() *big.Int {
	result, err := f.getFromCache(f.index)
	if err != nil {
		panic(err)
	}

	return result
}

func (f *fibStore) getPrevious() *big.Int {
	current := f.index

	if current == 0 {
		// cache hit
		result, err := f.getFromCache(current)
		if err != nil {
			panic(err)
		}
		return result
	}

	current -= 1
	f.index = current
	result, err := f.getFromCache(current)
	if err != nil {
		if err == CacheMiss {
			// cache miss
			a, err := f.getFromCache(current + 1)
			if err != nil {
				panic(err)
			}

			b, err := f.getFromCache(current + 2)
			if err != nil {
				panic(err)
			}
			result := new(big.Int)
			result.Sub(b, a)

			f.addToCache(current, result)
			return result
		} else {
			panic(err)
		}
	}
	// cache hit
	return result
}
