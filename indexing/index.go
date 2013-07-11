package indexing

import (
	"bufio"
	"encoding/binary"
	"io"
	"strings"
)

const (
  magic = "searchme"
)

type Index struct {
	DocCounter  int64
	docIdToPath map[int64]string
	pathToDocId map[string]int64
	postings    map[string][]int64
}

func NewIndex() *Index {
	var index Index
	index.postings = make(map[string][]int64)
	index.docIdToPath = make(map[int64]string)
	index.pathToDocId = make(map[string]int64)
	return &index
}

func (i *Index) Add(doc Document) error {
	termSet := make(map[string]bool)
	for _, term := range TokenizeString(doc.Content) {
		termSet[term] = true
	}

	docId := i.DocId(doc.Path)

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

func (index *Index) Write(w io.Writer) {
  // List of files, null terminated. Doc ids correspond to index. Ends with an empty filename. 
  io.WriteString(w, magic + "\x00")
  for i := int64(0); i < index.DocCounter; i++ {
    io.WriteString(w, index.docIdToPath[i] + "\x00")
  }
  io.WriteString(w, "\x00")

  for term, docIds := range index.postings {
    // TERM \x00 64bit doc ids until \x00
    io.WriteString(w, term)
    io.WriteString(w, "\x00")
    for _, id := range docIds {
      binary.Write(w, binary.BigEndian, id)
      io.WriteString(w, "\x00")
    }
  }
}

func Load(reader io.Reader) (*Index, error) {
  index := NewIndex()
  r := bufio.NewReader(reader)
  magic, err := r.ReadBytes('\x00')
  if err != nil {
    return nil, err
  }
  if string(magic) != "searchme\x00" {
    panic("Bad format. Magic bytes not detected.")
  }
  var name []byte;
  for i := int64(0); len(name) != 1; name, err = r.ReadBytes('\x00') {
    pth := string(name[:len(name) - 2])
    index.docIdToPath[i] = pth // skip \x00
    index.pathToDocId[pth] = i // skip \x00
    i++
  }

  return index, nil
}
