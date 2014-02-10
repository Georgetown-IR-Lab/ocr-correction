package lexicon

import (
  "sort"
)

type Lexicon struct {  
  words []string
}

// Initializes and returns a new Lexicon
func NewLexicon() (*Lexicon) {
  l := new(Lexicon)
  l.words = make([]string, 0)
  return l
}

func (l *Lexicon) Add(word string) {
  l.words = append(l.words, word)
}

func (l *Lexicon) Sort() {
  sort.Strings(l.words[:])
}

func (l *Lexicon) Words() []string {
  return l.words
}