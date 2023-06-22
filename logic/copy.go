package logic

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var copyWg sync.WaitGroup

func Copy(commandArr []string, matches *[]Data) bool {
	if !hasArgsCount(3, &commandArr) {
		return false
	}

	absPath, err := filepath.Abs(commandArr[2])
	if err != nil {
		fmt.Println("Invalid copy destination")
		return false
	}
	pathOpen, err := os.Open(absPath)
	if os.IsPermission(err) {
		fmt.Println("No permission to access copy destination")
		return false
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			fmt.Printf("Failed to create folder: %v\r\n", err)
			return false
		}
	} else if err != nil {
		fmt.Printf("Failed to read folder: %v", err)
		return false
	}
	pathOpen.Close()

	copyAll := false
	if commandArr[1] == "*" {
		copyAll = true
	}

	fileNr := 0
	if !copyAll {
		nr, err := strconv.Atoi(commandArr[1])
		if err != nil {
			fmt.Println("Invalid Parameter")
			return false
		}
		if nr <= 0 || len(*matches) < nr {
			fmt.Println("Selected file out of scope")
			return false
		}
		fileNr = nr
	}

	if !copyAll {
		from := fmt.Sprintf("%v%v", *(*matches)[fileNr].Path, (*matches)[fileNr].Name)
		to := fmt.Sprintf("%v\\%v", absPath, (*matches)[fileNr].Name)
		if err := copy(from, to); err != nil {
			log.Printf("Failed to copy file: %v [ %v ] => [ %v ]", err, from, to)
		}
	} else {
		for index, matchingFile := range *matches {
			copyWg.Add(1)
			go func(matchingFile Data, index int) {
				defer copyWg.Done()
				from := fmt.Sprintf("%v%v", *matchingFile.Path, matchingFile.Name)
				to := fmt.Sprintf("%v\\nr%d-%v", absPath, index, matchingFile.Name)
				if err := copy(from, to); err != nil {
					log.Printf("Failed to copy file: %v [ %v ] => [ %v ]", err, from, to)
				}
			}(matchingFile, index)
		}
		copyWg.Wait()
	}

	return true
}

func copy(from string, to string) error {
	bytesRead, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(to, bytesRead, 0755)
	if err != nil {
		return err
	}

	return nil
}