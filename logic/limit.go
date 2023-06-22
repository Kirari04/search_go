package logic

import (
	"fmt"
	"strconv"
)

func Limit(commandArr []string, matchsOutputLimit *int) bool {
	if !hasArgsCount(2, &commandArr) {
		return false
	}
	nr, err := strconv.Atoi(commandArr[1])
	if err != nil {
		fmt.Println("Invalid Parameter")
		return false
	}
	if nr < 0 {
		fmt.Println("Selected limit out of scope")
		return false
	}
	*matchsOutputLimit = nr

	return true
}
