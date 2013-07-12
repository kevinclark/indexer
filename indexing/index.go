package indexing

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
  magic = "searchme"
)

type Index struct {
	DocCounter  uint64
	docIdToPath map[uint64]string
	pathToDocId map[string]uint64
	postings    map[string][]uint64
}

func NewIndex() *Index {
	var index Index
	index.postings = make(map[string][]uint64)
	index.docIdToPath = make(map[uint64]string)
	index.pathToDocId = make(map[string]uint64)
	return &index
}

func (i *Index) Add(doc *Document) error {
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

func (i *Index) DocId(path string) uint64 {
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
  for i := uint64(0); i < index.DocCounter; i++ {
    io.WriteString(w, index.docIdToPath[i] + "\x00")
  }
  io.WriteString(w, "\x00")

  for term, docIds := range index.postings {
    // TERM \x00 32-bit-number-of-docs 64bit doc ids until
    io.WriteString(w, term)
    io.WriteString(w, "\x00")
    binary.Write(w, binary.BigEndian, uint32(len(docIds)))
    for _, id := range docIds {
      binary.Write(w, binary.BigEndian, id)
    }
  }
}

func stripNull(b []byte) string {
  return string(b[:len(b) - 1])
}

func LoadIndex(reader io.Reader) (*Index, error) {
  index := NewIndex()
  r := bufio.NewReader(reader)

  magic, _ := r.ReadBytes('\x00')
  if stripNull(magic) != "searchme" {
    return nil, errors.New(fmt.Sprintf("Bad format. Magic bytes not detected: %q", magic))
  }

  // Read docs
  var i uint64
  for p, _ := r.ReadBytes('\x00'); len(p) >= 2; p, _ = r.ReadBytes('\x00') {
    path := stripNull(p)
    index.docIdToPath[i] = path // skip \x00
    index.pathToDocId[path] = i // skip \x00
    i++
  }

  // Read terms
  for {
    t, err := r.ReadBytes('\x00')
    if err != nil {
      return index, nil
    }
    term := stripNull(t)
    var size uint32
    binary.Read(r, binary.BigEndian, &size)
    docs := make([]uint64, size)
    for j := uint32(0); j < size; j++ {
      var docId uint64
      binary.Read(r, binary.BigEndian, &docId)
      docs[j] = docId
    }
    index.postings[term] = docs
  }

  return index, nil
}
