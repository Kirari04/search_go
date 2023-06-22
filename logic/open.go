package logic

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

func Open(commandArr []string, matches *[]Data) bool {
	if !hasArgsCount(2, &commandArr) {
		return false
	}
	fileNr, err := strconv.Atoi(commandArr[1])
	if err != nil {
		fmt.Println("Invalid Parameter")
		return false
	}
	if fileNr <= 0 || len(*matches) < fileNr {
		fmt.Println("Selected file out of scope")
		return false
	}
	fmt.Printf("Opening: %v", *(*matches)[fileNr-1].Path)
	if err != nil {
		fmt.Printf("Failed to get abs path: %v\r\n", err)
	}
	cmd := exec.Command("explorer")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       fmt.Sprintf(` /select,"%s%s"`, *(*matches)[fileNr-1].Path, (*matches)[fileNr-1].Name),
		CreationFlags: 0,
	}
	cmd.Run()

	return true
}
