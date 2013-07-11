package indexing

import (
  "strings"
  "testing"
)

var tokenizeTests = []struct{
  in string
  out []string
}{
  // Single word
  {"foo", []string{"foo"}},
  // Space separated multi word
  {"foo bar", []string{"foo", "bar"}},
  // Tab separated multi word
  {"foo\tbar", []string{"foo", "bar"}},
  // Trailing space
  {"foo bar ", []string{"foo", "bar"}},
  // Sanitize end of single word
  {"foo.", []string{"foo"}},
  // Sanitize middle of single word
  {"o.k.", []string{"ok"}},
  // Sanitize multi word
  {"o.k. k.c.", []string{"ok", "kc"}},
}

func TestTokenize(t *testing.T) {
  for _, testCase := range tokenizeTests {
    tokens := Tokenize(strings.NewReader(testCase.in))
    if len(testCase.out) != len(tokens) {
      t.Fatal("Incorrect number of tokens. Expected: ",
        testCase.out, " but got", tokens)

      for i, token := range testCase.out {
        if tokens[i] != token {
          t.Fatal("Diff at index ", i,
            ". Expected: ", testCase.out, " actual: ", tokens)
        }
      }
    }
  }
}
