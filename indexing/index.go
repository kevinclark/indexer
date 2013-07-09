package indexing

import (
	"net/mail"
	"os"
	"strings"
)

type Index struct {
	DocCounter  int64
	docIdToPath map[int64]string
	pathToDocId map[string]int64
	postings    map[string][]int64
}

func NewIndex() Index {
	var index Index
	index.postings = make(map[string][]int64)
	index.docIdToPath = make(map[int64]string)
	index.pathToDocId = make(map[string]int64)
	return index
}

func (i *Index) Add(path string) error {
	// TODO(kev): No duplicate adds
	file, err := os.Open(path)

	if err != nil {
		return nil
	}

	defer file.Close()
	var msg *mail.Message
	msg, err = mail.ReadMessage(file)
	if err != nil {
		return nil
	}

	termSet := make(map[string]bool)
	for _, term := range Tokenize(msg.Body) {
		termSet[term] = true
	}

	docId := i.DocId(path)

	for term, _ := range termSet {
		i.postings[term] = append(i.postings[term], docId)
	}

	return nil
}

func (i *Index) DocId(path string) int64 {
	docId, contains := i.pathToDocId[path]
	if !contains {
		docId = i.DocCounter
		i.pathToDocId[path] = docId
		i.docIdToPath[docId] = path
		i.DocCounter++
	}
	return docId
}

func (index *Index) TermPaths(term string) []string {
	docs, found := index.postings[strings.ToLower(term)]
	if !found {
		return make([]string, 0)
	}
	results := make([]string, len(docs))
	for i, docId := range docs {
		results[i] = index.docIdToPath[docId]
	}
	return results
}
