package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type fib struct {
	index    int
	sequence map[int]*big.Int
}

func (f *fib) next(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	current := f.index
	current += 1

	f.index = current

	if current == 1 {
		fmt.Fprintf(w, "1")
		return
	}

	a := f.sequence[current-1]
	b := f.sequence[current-2]

	sum := new(big.Int)
	sum.Add(a, b)

	f.sequence[current] = sum

	fmt.Fprintf(w, sum.String())
}

func (f *fib) current(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, f.sequence[f.index])
}

func (f *fib) previous(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	current := f.index

	if current == 0 {
		fmt.Fprint(w, f.sequence[current])
		return
	}

	current -= 1
	f.index = current

	fmt.Fprint(w, f.sequence[current])
}

func intializeSequence() *fib {
	fib := new(fib)
	fib.index = 0

	fib.sequence = make(map[int]*big.Int)
	fib.sequence[0] = big.NewInt(0)
	fib.sequence[1] = big.NewInt(1)

	return fib
}

func main() {

	fib := intializeSequence()

	router := httprouter.New()

	router.GET("/next", fib.next)
	router.GET("/current", fib.current)
	router.GET("/previous", fib.previous)

	log.Fatal(http.ListenAndServe(":8080", router))
}
