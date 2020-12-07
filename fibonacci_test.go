package main

import (
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestGetFromCache(t *testing.T) {
	// arrange
	fib, err := intializeCache(2)
	if err != nil {
		log.Fatal(err)
	}

	var tests = []struct {
		a    int
		want *big.Int
	}{
		{0, big.NewInt(0)},
		{1, big.NewInt(1)},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.a)
		t.Run(testname, func(t *testing.T) {

			// act
			fib.getFromCache(tt.a)
			ans, err := fib.getFromCache(tt.a)
			if err != nil {
				t.Fatal(err)
			}

			// assert
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestAddToCache(t *testing.T) {
	// arrange
	fib, err := intializeCache(5)
	if err != nil {
		log.Fatal(err)
	}

	var tests = []struct {
		index int
		value *big.Int
	}{
		{2, big.NewInt(1)},
		{3, big.NewInt(2)},
		{4, big.NewInt(3)},
		{5, big.NewInt(5)},
		{10, big.NewInt(55)},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.index)
		t.Run(testname, func(t *testing.T) {

			// act
			fib.addToCache(tt.index, tt.value)
			ans, err := fib.getFromCache(tt.index)
			if err != nil {
				t.Fatal(err)
			}

			// assert
			if ans.Cmp(tt.value) != 0 {
				t.Errorf("got %d, want %d", ans, tt.value)
			}
		})
	}
}

func TestBuildSequenceToIndex(t *testing.T) {
	// arrange
	fib, err := intializeCache(500)
	if err != nil {
		t.Fatal(err)
	}

	const testCases = 12

	var numStrings [testCases]string
	numStrings[0] = "0"
	numStrings[1] = "1"
	numStrings[2] = "1"
	numStrings[3] = "2"
	numStrings[4] = "3"
	numStrings[5] = "5"
	numStrings[6] = "55"
	numStrings[7] = "75025"
	numStrings[8] = "12586269025"
	numStrings[9] = "354224848179261915075"
	numStrings[10] = "222232244629420445529739893461909967206666939096499764990979600"
	numStrings[11] = "139423224561697880139724382870407283950070256587697307264108962948325571622863290691557658876222521294125"

	var bigNums [testCases]*big.Int
	for i := 0; i < testCases; i++ {
		n := new(big.Int)
		n, ok := n.SetString(numStrings[i], 10)
		if !ok {
			t.Fatal("Errored creating test cases during SetString call.")
		}
		bigNums[i] = n
	}

	var tests = []struct {
		a    int
		want *big.Int
	}{
		{0, bigNums[0]},
		{1, bigNums[1]},
		{2, bigNums[2]},
		{3, bigNums[3]},
		{4, bigNums[4]},
		{5, bigNums[5]},
		{10, bigNums[6]},
		{25, bigNums[7]},
		{50, bigNums[8]},
		{100, bigNums[9]},
		{300, bigNums[10]},
		{500, bigNums[11]},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.a)
		t.Run(testname, func(t *testing.T) {

			// act
			fib.buildSequenceToIndex(tt.a)
			ans, err := fib.getFromCache(tt.a)
			if err != nil {
				t.Fatal(err)
			}

			// assert
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestGetCurrent(t *testing.T) {
	// arrange
	fib, err := intializeCache(20)
	if err != nil {
		log.Fatal(err)
	}

	const testCases = 12

	var tests = []struct {
		index int
		want  *big.Int
	}{
		{0, big.NewInt(0)},
		{1, big.NewInt(1)},
		{2, big.NewInt(1)},
		{3, big.NewInt(2)},
		{4, big.NewInt(3)},
		{5, big.NewInt(5)},
		{6, big.NewInt(8)},
		{7, big.NewInt(13)},
		{8, big.NewInt(21)},
		{9, big.NewInt(34)},
		{10, big.NewInt(55)},
		{11, big.NewInt(89)},
	}

	fib.buildSequenceToIndex(12)

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.want)
		t.Run(testname, func(t *testing.T) {

			// act
			fib.index = tt.index
			ans := fib.getCurrent()

			// assert
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestGetNext(t *testing.T) {

	// arrange
	fib, err := intializeCache(20)
	if err != nil {
		log.Fatal(err)
	}

	const testCases = 12

	var tests = []struct {
		want *big.Int
	}{
		{big.NewInt(1)},
		{big.NewInt(1)},
		{big.NewInt(2)},
		{big.NewInt(3)},
		{big.NewInt(5)},
		{big.NewInt(8)},
		{big.NewInt(13)},
		{big.NewInt(21)},
		{big.NewInt(34)},
		{big.NewInt(55)},
		{big.NewInt(89)},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.want)
		t.Run(testname, func(t *testing.T) {

			// act
			ans := fib.getNext()

			// assert
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestGetPrevious(t *testing.T) {
	// arrange
	fib, err := intializeCache(20)
	if err != nil {
		log.Fatal(err)
	}

	var tests = []struct {
		want *big.Int
	}{
		{big.NewInt(89)},
		{big.NewInt(55)},
		{big.NewInt(34)},
		{big.NewInt(21)},
		{big.NewInt(13)},
		{big.NewInt(8)},
		{big.NewInt(5)},
		{big.NewInt(3)},
		{big.NewInt(2)},
		{big.NewInt(1)},
		{big.NewInt(1)},
		{big.NewInt(0)},
	}

	fib.buildSequenceToIndex(12)

	fib.index = 12

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.want)
		t.Run(testname, func(t *testing.T) {
			// act
			ans := fib.getPrevious()

			// assert
			if ans.Cmp(tt.want) != 0 {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
