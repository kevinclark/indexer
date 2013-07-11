package indexing

import (
  "bytes"
  "net/mail"
  "os"
)

type Document struct {
  Path string
  Content string
}

func LoadFromFile(path string, transform func (*os.File) (string, error)) (*Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

  content, err := transform(file)
  if err != nil {
    return nil, err
  }

  return &Document{path, content}, nil
}

func LoadFromMail(path string) (*Document, error) {
  return LoadFromFile(path, func (file *os.File) (string, error) {
    msg, err := mail.ReadMessage(file)
    if err != nil {
      return "", err
    }
    buf := new(bytes.Buffer)
    buf.ReadFrom(msg.Body)

    return string(buf.Bytes()), nil
  })
}
