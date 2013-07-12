package main

import (
	"flag"
	"fmt"
	"github.com/kevinclark/indexer/indexing"
	"os"
	"path/filepath"
)

func query(path string) {
  file, err := os.Open(path)
  if err != nil {
    fmt.Println("Unable to load index: ", err)
  }
  defer file.Close()

  index, err := indexing.LoadIndex(file)

  if err != nil {
    fmt.Println("Unable to load index: ", err)
  }
  

	for {
		fmt.Println("Query?")
		var query string
		_, err := fmt.Scanf("%s", &query)
		if err != nil {
			panic(err)
		}
		for _, path := range index.TermPaths(query) {
			fmt.Println("Found in: ", path)
		}
	}
}

func write(path string, outPath string) {
  index := indexing.NewIndex()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

    doc, err := indexing.LoadFromMail(path)
    if err != nil {
      fmt.Println("Unable to load mail: ", err)
      return nil
    }
		index.Add(doc)

		return nil
	})

	fmt.Println("Processed ", index.DocCounter, " files")
  file, err := os.Create(outPath)
  if err != nil {
    fmt.Println("Unable to write index: ", err)
  }
  defer file.Close()
  index.Write(file)
}

// Walk the given file path, adding emails to the index
func main() {
	// The path we're going to index comes from the command line
	flag.Parse()

  args := flag.Args()

  switch args[0] {
  case "index":
    write(args[1], args[2])
  case "query":
    query(args[1])
  }
}
