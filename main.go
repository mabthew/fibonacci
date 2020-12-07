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

type fileError struct {
	cause   string
	message error
}

func (e *fileError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.message)
}

func createOrReadBackupFile(path string) (int, error) {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		log.Println(&fileError{"File does not exist", err})

		var file, err = os.Create(path)
		if err != nil {
			log.Fatal("Failed to create backup file:", err)
		}
		defer file.Close()
		log.Println("Backup file successfully created. Starting sequence at index 0.")
		return 0, nil
	} else {
		var file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return 0, &fileError{"Failed to open backup file", err}
		}
		defer file.Close()

		fileinfo, err := file.Stat()
		if err != nil {
			return 0, err
		}

		filesize := fileinfo.Size()
		text := make([]byte, filesize)

		_, err = file.Read(text)
		if err != nil {
			return 0, &fileError{"Failed to read backup file", err}
		}

		index, err := strconv.Atoi(string(text))
		if err != nil {
			return 0, &fileError{"Failed to parse backup file", err}
		}

		log.Printf("Backup file recovered successfully, starting at index %v.\n", strconv.Itoa(index))
		return index, nil
	}
}

func writeToBackupFile(path string, index int) {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Println(&fileError{"Failed to open backup file", err})
	}
	defer file.Close()

	// backup to file
	err = file.Truncate(0)
	_, err = fmt.Fprintf(file, "%d", index)
	if err != nil {
		log.Println(&fileError{"Failed to write to backup file", err})
		return
	}

	// save file changes
	err = file.Sync()
	if err != nil {
		log.Println(&fileError{"Failed to save backup file changes", err})
		return
	}
}

func (f *fibStore) backup(path string) {
	for {
		<-time.After(2 * time.Second)
		go writeToBackupFile(path, f.index)
	}
}

func (f *fibStore) next(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := f.getNext()
	fmt.Fprint(w, result)
}

func (f *fibStore) current(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := f.getCurrent()
	fmt.Fprint(w, result)
}

func (f *fibStore) previous(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := f.getPrevious()
	fmt.Fprint(w, result)
}

func main() {

	const path = "backup.txt"
	const cacheSize = 10000

	fib, err := intializeCache(cacheSize)
	if err != nil {
		log.Fatal(err)
	}

	if index, err := createOrReadBackupFile(path); err != nil {
		log.Println(err)
		log.Println("Starting sequence at index 0.")
	} else {
		fib.buildSequenceToIndex(index)
	}

	go fib.backup(path)

	router := httprouter.New()

	router.GET("/next", fib.next)
	router.GET("/current", fib.current)
	router.GET("/previous", fib.previous)

	log.Println("Server listening on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", router))
}
