package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var path = "backup.txt"

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

func (f *fibStore) backup() {
	for {
		<-time.After(2 * time.Second)
		go writeToBackupFile(f.index)
	}
}

func (f *fibStore) next(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, f.getNext().String())
}

func (f *fibStore) current(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, f.getCurrent().String())
}

func (f *fibStore) previous(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Fprint(w, f.getPrevious().String())
}

func main() {

	fib := intializeFibStore()

	router := httprouter.New()

	index := createOrReadBackupFile()

	if index != 0 {
		fib.buildSequenceToIndex(index)
	}

	go fib.backup()

	router.GET("/next", fib.next)
	router.GET("/current", fib.current)
	router.GET("/previous", fib.previous)

	fmt.Printf("Listening on port 8080...\n")

	log.Fatal(http.ListenAndServe(":8080", router))
}
