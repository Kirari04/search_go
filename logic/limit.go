package logic

import (
	"log"
	"strconv"
)

func Limit(commandArr []string, matchsOutputLimit *int) bool {
	if !hasArgsCount(2, &commandArr) {
		return false
	}
	nr, err := strconv.Atoi(commandArr[1])
	if err != nil {
		log.Println("Invalid Parameter")
		return false
	}
	if nr < 0 {
		log.Println("Selected limit out of scope")
		return false
	}
	*matchsOutputLimit = nr

	return true
}
