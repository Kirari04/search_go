package benchmark_test

import (
	"com/github/kirari04/search_go/initializer"
	"com/github/kirari04/search_go/logic"
	"fmt"
	"testing"
)

func BenchmarkMP4Search(b *testing.B) {
	rootdir := "/home/lev/"
	entries := logic.SafeEntries{}
	silent := true
	maxdebth := 10
	pathSeperator := `/`
	isRegex := true
	matchsOutputLimit := 0
	initializer.Index(
		rootdir,
		&entries,
		silent,
		maxdebth,
		pathSeperator,
	)
	b.ResetTimer()

	for threads := 1; threads < 40; threads++ {
		b.Run(fmt.Sprintf("input_size_%d", threads), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logic.Search(
					"^*.mp4$",
					isRegex,
					&entries,
					matchsOutputLimit,
					threads,
					true,
				)
			}
		})
	}

}
