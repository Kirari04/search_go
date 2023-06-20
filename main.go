package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var rootdir = `C:\`
var maxdebth int = 10
var indexCount int = 0
var wg sync.WaitGroup
var wgSearch sync.WaitGroup
var silent = true
var isIndexed = false
var isRegex = true

type Data struct {
	Name string
	Path *string
}

var entries []Data = make([]Data, 100000)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("######################")
	fmt.Println("#      search_go     #")
	fmt.Println("######################")

	if os.Getenv("ROOTDIR") != "" {
		if _, err := os.ReadDir(rootdir); err != nil {
			log.Panicf("Failed to set ROOTDIR: %v", err)
		}
		rootdir = os.Getenv("ROOTDIR")
	}

	if os.Getenv("MAXDEBTH") != "" {
		maxdebth_tmp, err := strconv.Atoi(os.Getenv("MAXDEBTH"))
		if err != nil {
			log.Panicf("Failed to set MAXDEBTH: %v", err)
		}
		if maxdebth_tmp < 0 {
			log.Panicf("Failed to set MAXDEBTH: MAXDEBTH has to be >= 0")
		}
		maxdebth = maxdebth_tmp
	}

	if os.Getenv("SILENT") != "" {
		silent_tmp, err := strconv.ParseBool(os.Getenv("SILENT"))
		if err != nil {
			log.Panicf("Failed to set SILENT: %v", err)
		}
		silent = silent_tmp
	}

	if os.Getenv("REGEX") != "" {
		isregex_tmp, err := strconv.ParseBool(os.Getenv("REGEX"))
		if err != nil {
			log.Panicf("Failed to set REGEX: %v", err)
		}
		isRegex = isregex_tmp
	}

	fmt.Printf(
		"ROOTDIR='%v'\r\nMAXDEBTH='%v'\r\nSILENT='%v'\r\nREGEX='%v'\r\n",
		rootdir,
		maxdebth,
		silent,
		isRegex,
	)

	for {
		if !isIndexed {
			fmt.Printf("the folders aren't indexed yet - press enter to index them.\r\n")
		}
		fmt.Print("-> ")
		line, _, _ := reader.ReadLine()
		text := string(line)
		if text == "" && isIndexed {
			fmt.Println("search term can't be empty")
			continue
		}

		if !isIndexed {
			start := time.Now()
			fmt.Printf("Indexing folders...\r\n")
			listAllDirs(&entries, rootdir, 0)
			wg.Wait()
			log.Printf(
				"It took %vms to index %v files & folders\r\n",
				time.Now().UnixMilli()-start.UnixMilli(),
				indexCount,
			)
			isIndexed = true
			continue
		}

		var reg *regexp.Regexp
		if isRegex {
			r, err := regexp.Compile(text)
			if err != nil {
				fmt.Println("regex can't be parsed")
				continue
			}
			reg = r
		}

		start := time.Now()
		var matchCount int
		for i := range entries {
			wgSearch.Add(1)
			go func(e *Data) {
				defer wgSearch.Done()
				if !isRegex {
					if (*e).Name == text {
						matchCount++
						fmt.Printf("[%v] %v\r\n", e.Path, e.Name)
					}
				}
				if isRegex {
					if reg.MatchString((*e).Name) {
						matchCount++
						fmt.Printf("[%v] %v\r\n", *e.Path, e.Name)
					}
				}

			}(&entries[i])
		}
		wgSearch.Wait()
		fmt.Printf(
			"Found %v matches in %vms\r\n",
			matchCount,
			time.Now().UnixMilli()-start.UnixMilli(),
		)
	}
}

func listAllDirs(entries *[]Data, dir string, debth int) {
	mewEntries, err := os.ReadDir(dir)
	if err != nil && !silent {
		log.Printf("Warning: %v [%v]", err, dir)
		return
	}

	wg.Add(1)
	defer wg.Done()
	for _, nE := range mewEntries {
		*entries = append(*entries, Data{
			Name: nE.Name(),
			Path: &dir,
		})
		indexCount++
	}

	var wgListAllDirs sync.WaitGroup
	for _, e := range mewEntries {
		wgListAllDirs.Add(1)
		go func(e fs.DirEntry) {
			defer wgListAllDirs.Done()
			if e.IsDir() && debth <= maxdebth {
				wg.Add(1)
				defer wg.Done()
				listAllDirs(entries, fmt.Sprintf("%v/%v", dir, e.Name()), debth+1)
			}
		}(e)
	}
	wgListAllDirs.Wait()
}
