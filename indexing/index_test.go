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
    // Header | Doc paths + Terminator | term + terminator + doc id (0)
    "searchme\x00path\x00\x00content\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"},
}

func TestWrite(t *testing.T) {
  for _, testCase := range testWrites {
    i := NewIndex()
    buf := new(bytes.Buffer)
    for _, doc := range testCase.docs {
      i.Add(doc)
    }
    i.Write(buf)

    result := string(buf.Bytes())
    if testCase.out != result {
      t.Fatalf("Expected: %q Actual %q", testCase.out, result)
    }
  }
}

