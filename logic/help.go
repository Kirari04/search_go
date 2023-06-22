package logic

import "fmt"

func Help() bool {
	fmt.Println("#########################")
	fmt.Println("# search_go | help page #")
	fmt.Println("#########################")
	fmt.Println(`
	--open [number]         # requires a number
	--open 69               # opens the 69th file from the matched files & folders in the explorer
	
	--limit [number]        # requires a number
	--limit 100             # list only the first 100 matches (set to 0 if you want no limit)

	--copy [number|*]       # requires a number or *
	--copy 100 ./destfolder # copies the 100th file into the destfolder
	--copy * ./destfolder   # copies the all matched files into the destfolder
	`)
	return true
}
