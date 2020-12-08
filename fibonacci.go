package main

import (
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/golang/groupcache/lru"
)

// mutex to lock fibStore
var mutex = &sync.Mutex{}

// struct storing current index in series and lru cache allowing access to numbers in the series
// by their corresponding indices
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

func (f *fibStore) getFromCache(index int) (*big.Int, error) {
	result, ok := f.cache.Get(index)
	if ok == false {
		return big.NewInt(-1), CacheMiss
	}
	return result.(*big.Int), nil
}

func (f *fibStore) addToCache(index int, value *big.Int) {
	f.cache.Add(index, value)
}

func (f *fibStore) buildSequenceToIndex(index int) {

	// build fibonacci sequence to index provided
	f.addToCache(0, big.NewInt(0))
	f.addToCache(1, big.NewInt(1))

	for i := 2; i <= index; i++ {
		sum := new(big.Int)
		a, err := f.getFromCache(i - 1)
		if err != nil {
			log.Fatal("Failed while recovering sequence.")
		}

		b, err := f.getFromCache(i - 2)
		if err != nil {
			log.Fatal("Failed while recovering sequence.")
		}

		sum.Add(a, b)

		f.addToCache(i, sum)
	}

	f.index = index
}

func (f *fibStore) attemptHardRecover(index int) *big.Int {

	// attempt to rebuild cache and retrieve desired index
	f.buildSequenceToIndex(index)
	value, err := f.getFromCache(index)

	// if rebuilding the cache and attempting to get that index doesn't work, there is a larger issue
	if err != nil {
		log.Fatal("Cache corrupt: hard restart required.")
	}

	return value
}

func (f *fibStore) getNext() *big.Int {
	mutex.Lock()
	defer mutex.Unlock()

	// increment current index
	current := f.index
	current += 1
	f.index = current

	// value at index 1 should be cached when index was 0 upon entering function
	if current == 1 {
		result, err := f.getFromCache(1)
		if err != nil {

			result = f.attemptHardRecover(1)
		}

		return result
	}

	result, err := f.getFromCache(current)
	if err != nil {
		if err == CacheMiss {
			// cache miss - calculate fibonacci
			a, err := f.getFromCache(current - 1)
			if err != nil {
				a = f.attemptHardRecover(current - 1)
			}

			b, err := f.getFromCache(current - 2)
			if err != nil {
				b = f.attemptHardRecover(current - 2)
			}

			sum := new(big.Int)
			sum.Add(a, b)

			f.addToCache(current, sum)

			return sum
		}
	}

	// cache hit
	return result

}

func (f *fibStore) getCurrent() *big.Int {
	result, err := f.getFromCache(f.index)
	if err != nil {
		result = f.attemptHardRecover(f.index)
	}
	return result
}

func (f *fibStore) getPrevious() *big.Int {
	mutex.Lock()
	defer mutex.Unlock()

	current := f.index

	// don't look for previous at index 0
	if current == 0 {
		// cache hit
		result, err := f.getFromCache(current)
		if err != nil {
			result = f.attemptHardRecover(current)
		}
		return result
	}

	current -= 1
	f.index = current
	result, err := f.getFromCache(current)
	if err != nil {
		if err == CacheMiss {
			// cache miss - calculate fibonacci
			a, err := f.getFromCache(current + 1)
			if err != nil {
				a = f.attemptHardRecover(current + 1)
			}

			b, err := f.getFromCache(current + 2)
			if err != nil {
				b = f.attemptHardRecover(current + 2)
			}
			result := new(big.Int)
			result.Sub(b, a)

			f.addToCache(current, result)
			return result
		}
	}
	// cache hit
	return result
}
