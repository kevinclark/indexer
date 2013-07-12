package indexing

import (
  "bytes"
  "testing"
)

var testWrites = []struct {
  docs []Document
  out string
}{
  {[]Document{},
    // Header | End of Docs
    "searchme\x00\x00"}, // terminator + doc terminator
  {[]Document{{"path", "content"}},
    // Header | Doc paths + Terminator | term + terminator + num-docs + doc id (0)
    "searchme\x00path\x00\x00content\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00"},
}

func TestWrite(t *testing.T) {
  for testNum, testCase := range testWrites {
    i := NewIndex()
    buf := new(bytes.Buffer)
    for _, doc := range testCase.docs {
      i.Add(&doc)
    }
    i.Write(buf)

    result := string(buf.Bytes())
    if testCase.out != result {
      t.Fatalf("%d. Expected: %q Actual %q", testNum, testCase.out, result)
    }
  }
}

func TestLoadIndexWithBadHeader(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("Not the header you're looking for")
  i, err := LoadIndex(buf)
  if err == nil {
    t.Fatal("Got no error but missing header")
  }
  if i != nil {
    t.Fatal("Should have returned no index")
  }
}

func TestLoadIndexWithGoodHeader(t *testing.T) {
  buf := new(bytes.Buffer)
  buf.WriteString("searchme\x00path\x00\x00content\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00")
  i, err := LoadIndex(buf)
  if err != nil {
    t.Fatalf("Error: %q", err)
  }
  if i.pathToDocId["path"] != 0 {
    t.Fatalf("'path' should be doc id 0: %q", i.pathToDocId)
  }
  if len(i.postings["content"]) != 1 || i.postings["content"][0] != 0 {
    t.Fatalf("Should only have one doc for 'content' term: %v", i.postings["content"])
  }
}
