package main

import (
	"bufio"
	"com/github/kirari04/search_go/logic"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rootdir = `C:\`
var pathSeperator = `\`
var maxdebth int = 10
var indexCount int = 0
var wg sync.WaitGroup
var wgSearch sync.WaitGroup
var silent = true
var isIndexed = false
var isRegex = true
var matchsOutputLimit = 100

var entries []logic.Data = make([]logic.Data, 100000)
var matchCount int
var matches []logic.Data

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("######################")
	fmt.Println("#      search_go     #")
	fmt.Println("######################")
	initENV()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()

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
			fmt.Printf(
				"It took %vms to index %v files & folders\r\n",
				time.Now().UnixMilli()-start.UnixMilli(),
				indexCount,
			)
			isIndexed = true
			continue
		}

		isCommand := false
		if strings.HasPrefix(text, "--") {
			isCommand = true
			commandArr := strings.Split(text, " ")
			if len(commandArr) == 0 {
				fmt.Println("The command is empty")
				continue
			}
			switch commandArr[0] {
			case "--help":
				if !logic.Help() {
					continue
				}
			case "--open":
				if !logic.Open(commandArr, &matches) {
					continue
				}
			case "--copy":
				if !logic.Copy(commandArr, &matches) {
					continue
				}
			case "--limit":
				if !logic.Limit(commandArr, &matchsOutputLimit) {
					continue
				}
			default:
				log.Printf("Command not found: %v\r\nRun --help for a list of all commands", commandArr[0])
				continue
			}
		}
		if isCommand {
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
		matchCount = 0
		matches = []logic.Data{}
		for _, entry := range entries {
			// to not open unecessary new go routines if it has already got to its limit we check here
			if matchsOutputLimit == 0 || matchsOutputLimit > matchCount {
				wgSearch.Add(1)
				go func(e logic.Data) {
					defer wgSearch.Done()
					if matchsOutputLimit == 0 || matchsOutputLimit > matchCount {
						if !isRegex && e.Name == text {
							matchCount++
							fmt.Printf("[%v] %v\r\n", e.Path, e.Name)
							matches = append(matches, e)
						}
						if isRegex && reg.MatchString(e.Name) {
							matchCount++
							fmt.Printf("[%v] %v\r\n", matchCount, e.Name)
							matches = append(matches, e)
						}
					}
				}(entry)
			}
		}
		wgSearch.Wait()
		fmt.Printf(
			"Found %v matches in %vms\r\n",
			matchCount,
			time.Now().UnixMilli()-start.UnixMilli(),
		)
		if matchsOutputLimit != 0 && matchsOutputLimit <= matchCount {
			fmt.Printf("You'r current limit on matches is set to %v, you can increase this amount by running --limit [number]\r\n", matchsOutputLimit)
		}
	}
}

func initENV() {
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
}

func listAllDirs(entries *[]logic.Data, dir string, debth int) {
	if entries == nil {
		panic("entries is nil")
	}
	mewEntries, err := os.ReadDir(dir)
	if err != nil && !silent {
		log.Printf("Warning: %v [%v]", err, dir)
		return
	}

	for _, nE := range mewEntries {
		*entries = append(*entries, logic.Data{
			Name: nE.Name(),
			Path: &dir,
		})
		indexCount++
	}

	for _, e := range mewEntries {
		wg.Add(1)
		go func(isDir bool, name string, entries_ptr *[]logic.Data) {
			defer wg.Done()
			if isDir && debth <= maxdebth {
				listAllDirs(entries_ptr, fmt.Sprintf("%s%s%s", dir, name, pathSeperator), debth+1)
			}
		}(e.IsDir(), e.Name(), entries)
	}
}
