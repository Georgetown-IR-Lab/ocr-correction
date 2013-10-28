package filereader

import "strconv"

type DocumentId uint32

func (id DocumentId) MarshalJSON() ([]byte, error) {
  return []byte(strconv.Itoa(int(id))), nil
}

type DocInfo interface {
  OrigIdent() string
  Identifier() DocumentId
  Len() int /* The number of tokens in this document */
}

type Document interface {
  DocInfo
  Tokens() <-chan *Token
  Add(*Token) /* Add a token, setting the DocId and position if necessary */
}

type FileReader interface {
  Init(string)
  ReadAll() <-chan Document
  Read() Document
  Path() string
}

