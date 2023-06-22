package main

import (
	"com/github/kirari04/search_go/logic"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

var rootdir = `C:\`
var pathSeperator = `\`
var maxdebth int = 10
var silent = true
var isIndexed = false
var isRegex = true
var matchsOutputLimit = 100

var entries logic.SafeEntries = logic.SafeEntries{}

func main() {
	fmt.Println("######################")
	fmt.Println("#      search_go     #")
	fmt.Println("######################")

	initENV()

	Console()
}

func initENV() {
	if os.Getenv("ROOTDIR") != "" {
		rootdir = os.Getenv("ROOTDIR")
		if _, err := os.ReadDir(rootdir); err != nil {
			log.Panicf("Failed to set ROOTDIR: %v", err)
		}
	} else {
		switch runtime.GOOS {
		case "windows":
			rootdir = "c:\\"
		default:
			rootdir = "/"
		}
	}

	if os.Getenv("PATHSEP") != "" {
		pathSeperator = os.Getenv("PATHSEP")
	} else {
		switch runtime.GOOS {
		case "windows":
			pathSeperator = "\\"
		default:
			pathSeperator = "/"
		}
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
