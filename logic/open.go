package logic

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

func Open(commandArr []string, entries *SafeEntries) bool {
	if !hasArgsCount(2, &commandArr) {
		return false
	}
	fileNr, err := strconv.Atoi(commandArr[1])
	if err != nil {
		fmt.Println("Invalid Parameter")
		return false
	}
	if err != nil {
		fmt.Printf("Failed to get abs path: %v\r\n", err)
	}

	file, err := entries.Get(fileNr)
	if err != nil {
		fmt.Printf("Failed to get file by fileNr %v: %v\r\n", fileNr, err)
		return false
	}
	if !file.Matched {
		fmt.Println("You can't select unmatched files")
		return false
	}
	fmt.Printf("Opening: %v\r\n", *file.Path)

	if runtime.GOOS == "windows" {
		cmd := exec.Command("explorer", fmt.Sprintf(` /select,%s%s`, *file.Path, file.Name))
		cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("xdg-open", fmt.Sprintf(`%s%s`, *file.Path, file.Name))
		cmd.Run()
	} else {
		fmt.Println("Plattform not supported")
	}

	return true
}
