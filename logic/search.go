package logic

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sync"
	"time"
)

func Search(text string, isRegex bool, entries *SafeEntries, matchsOutputLimit int, routineCount int, noOutput bool) bool {
	var reg *regexp.Regexp
	if isRegex {
		r, err := regexp.Compile(text)
		if err != nil {
			fmt.Println("regex can't be parsed")
			return false
		}
		reg = r
	}

	start := time.Now()
	entrieCountPerRoutine := int(math.Ceil(float64(entries.Count()) / float64(routineCount)))
	if !noOutput {
		log.Printf("Searching on %v routines\r\n", routineCount)
	}
	matchCount := 0
	wgSearch := sync.WaitGroup{}
	for ri := 0; ri < routineCount; ri++ {
		wgSearch.Add(1)
		go func(ri int) {
			defer wgSearch.Done()
			from := ri * entrieCountPerRoutine
			to := (ri + 1) * entrieCountPerRoutine
			for i := from; i < to; i++ {
				if matchsOutputLimit == 0 || matchsOutputLimit > matchCount {
					e, _ := entries.Get(i)
					if e.Path == nil {
						continue
					}
					e.Matched = false
					if !isRegex && e.Name == text {
						matchCount++
						if !noOutput {
							fmt.Printf("[%d] %v\r\n", i, e.Name)
						}
						e.Matched = true
					}
					if isRegex && reg.MatchString(e.Name) {
						matchCount++
						if !noOutput {
							fmt.Printf("[%d] %v\r\n", i, e.Name)
						}
						e.Matched = true
					}
				}
			}
		}(ri)
	}
	wgSearch.Wait()

	if !noOutput {
		fmt.Printf(
			"Found %v matches in %vms\r\n",
			matchCount,
			time.Now().UnixMilli()-start.UnixMilli(),
		)
	}
	if matchsOutputLimit != 0 && matchsOutputLimit <= matchCount {
		if !noOutput {
			fmt.Printf("You'r current limit on matches is set to %v, you can increase this amount by running --limit [number]\r\n", matchsOutputLimit)
		}
	}

	return true
}
