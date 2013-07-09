package main

import (
	"flag"
	"fmt"
	"github.com/kevinclark/indexer/indexing"
	"os"
	"path/filepath"
)

// Walk the given file path, adding emails to the index
func main() {
	// The path we're going to index comes from the command line
	flag.Parse()

	index := indexing.NewIndex()

	filepath.Walk(flag.Args()[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		index.Add(path)

		return nil
	})

	fmt.Println("Processed ", index.DocCounter, " files")

	for {
		fmt.Println("Query?")
		var query string
		_, err := fmt.Scanf("%s", &query)
		if err != nil {
			panic("crap")
		}
		for _, path := range index.TermPaths(query) {
			fmt.Println("Found in: ", path)
		}
	}
}
