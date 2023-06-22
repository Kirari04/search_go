package main

import (
	"bufio"
	"com/github/kirari04/search_go/initializer"
	"com/github/kirari04/search_go/logic"
	"fmt"
	"log"
	"os"
	"strings"
)

func Console() {
	reader := bufio.NewReader(os.Stdin)
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
			initializer.Index(rootdir, &entries, silent, maxdebth, pathSeperator)
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
				if !logic.Open(commandArr, &entries) {
					continue
				}
			case "--copy":
				if !logic.Copy(commandArr, &entries, pathSeperator) {
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

		if !logic.Search(
			text,
			isRegex,
			&entries,
			matchsOutputLimit,
			8,
			false,
		) {
			continue
		}
	}
}
