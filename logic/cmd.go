package logic

import "log"

type Data struct {
	Name string
	Path *string
}

func hasArgsCount[T []string](count int, commandArguments *T) bool {
	res := !bool(count < len(*commandArguments) || count > len(*commandArguments))
	if !res {
		log.Printf("%v arguments are required\r\n", count-1)
	}

	return res
}
