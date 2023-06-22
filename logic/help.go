package logic

import "log"

func Help() bool {
	log.Println("#########################")
	log.Println("# search_go | help page #")
	log.Println("#########################")
	log.Println(`
	--open [number]			# requires a number
	--open 69				# opens the 69th file from the matched files & folders in the explorer
	
	--limit [number]		# requires a number
	--limit 100				# list only the first 100 matches (set to 0 if you want no limit)
	`)
	return true
}
