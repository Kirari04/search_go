package initializer

import (
	"com/github/kirari04/search_go/logic"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var indexCount int = 0
var silent bool
var maxdebth int = 100
var pathSeperator string = `/`

func Index(rootdir string, entries *logic.SafeEntries, newSilent bool, newMaxdebth int, newPathSeperator string) {
	pathSeperator = newPathSeperator
	maxdebth = newMaxdebth
	silent = newSilent
	start := time.Now()
	fmt.Printf("Indexing folders...\r\n")
	newEntries, _ := listAllDirs(rootdir, 0)
	entries.Set(*newEntries)
	fmt.Printf(
		"It took %vms to index %v files & folders\r\n",
		time.Now().UnixMilli()-start.UnixMilli(),
		indexCount,
	)
}

func listAllDirs(dir string, debth int) (*[]logic.Data, error) {
	mewEntries, err := os.ReadDir(dir)
	if err != nil && !silent {
		log.Printf("Warning: %v [%v]", err, dir)
		return nil, err
	}
	prepareEntryList := make([]logic.Data, len(mewEntries))
	for _, nE := range mewEntries {
		prepareEntryList = append(prepareEntryList, logic.Data{
			Name: nE.Name(),
			Path: &dir,
		})

		indexCount++
	}

	subentries := logic.SafeEntries{}
	var wg sync.WaitGroup
	for i := range mewEntries {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if mewEntries[i].IsDir() && debth <= maxdebth {
				subentry, err := listAllDirs(fmt.Sprintf("%s%s%s", dir, mewEntries[i].Name(), pathSeperator), debth+1)
				if err != nil {
					return
				}
				subentries.AddAll(*subentry)
			}
		}(i)
	}
	wg.Wait()

	prepareEntryList = append(prepareEntryList, *subentries.Value()...)
	return &prepareEntryList, nil
}
