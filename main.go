package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var path = "backup.txt"

type fib struct {
	index    int
	sequence map[int]*big.Int
}

func createOrReadBackupFile() int {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println("No backup file found, starting at index 0.")
		var file, err = os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		return 0
	} else {
		var file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fileinfo, err := file.Stat()
		if err != nil {
			panic(err)
		}

		filesize := fileinfo.Size()
		text := make([]byte, filesize)

		_, err = file.Read(text)
		if err != nil {
			panic(err)
		}

		index, err := strconv.Atoi(string(text))
		if err != nil || string(text) == "" {
			fmt.Println("Backup file corrupted, starting at index 0.")
			return 0
		}

		fmt.Printf("Backup file opened successfully, starting at index %v.\n", strconv.Itoa(index))
		return index
	}
}

func writeToBackupFile(index int) {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// backup to file
	err = file.Truncate(0)
	_, err = fmt.Fprintf(file, "%d", index)
	if err != nil {
		panic(err)
	}
	// Save file changes.
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func (f *fib) recover(recoveredIndex int) {

	for i := 2; i <= recoveredIndex; i++ {
		sum := new(big.Int)
		sum.Add(f.sequence[i-1], f.sequence[i-2])

		f.sequence[i] = sum
	}

	f.index = recoveredIndex
}

func (f *fib) backup() {
	for {
		<-time.After(2 * time.Second)
		go writeToBackupFile(f.index)
	}
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

	index := createOrReadBackupFile()

	if index != 0 {
		fib.recover(index)
	}

	go fib.backup()

	router.GET("/next", fib.next)
	router.GET("/current", fib.current)
	router.GET("/previous", fib.previous)

	fmt.Printf("Listening on port 8080...\n")

	log.Fatal(http.ListenAndServe(":8080", router))
}
