// By default, the lexicon is not case sensitive.
// This lexicon will not allow duplicate terms.  If
// you attempt to Add(xxx) a duplicate, it will simply
// increment the counter on that Term object.

package lexicon

import (
  "strings"
  "sort"
)

type Lexicon struct {  
  terms []*Term
}

// Initializes and returns a new Lexicon
func NewLexicon() (*Lexicon) {
  l := new(Lexicon)
  l.terms = make([]*Term, 0)
  return l
}

// Adds the term if it doesnt exist, otherwise it'll increment the
// count of that term
func (l *Lexicon) Add(word string) {
  term := l.Has(word)
  
  if term == nil {
    term := NewTerm(downcase(&word))
    l.terms = append(l.terms, term) 
  } else {
    term.IncrementCount()
  }  
}

func (l *Lexicon) SortByText() {
  sort.Sort(ByText(l.terms))
}

func (l *Lexicon) Terms() []*Term {
  return l.terms
}

func (l *Lexicon) Size() int {
  return len(l.terms)
}

func (l *Lexicon) Has(word string) *Term {
  for _, t := range l.terms {
    if t.Text == downcase(&word) {
      return t
    }
  }
  return nil
}

// Just a lazy way to simplify downcasing terms
func downcase(word *string) string {
  return strings.ToLower(*word)
}