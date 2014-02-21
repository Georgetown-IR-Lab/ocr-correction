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
  terms map[string]*Term
}

// Initializes and returns a new Lexicon
func NewLexicon() (*Lexicon) {
  l := new(Lexicon)
  l.terms = make(map[string]*Term)
  return l
}

// Adds the term if it doesnt exist, otherwise it'll increment the
// count of that term
func (l *Lexicon) Add(word string) {
  term, has_term := l.Has(word)

  if has_term {
    term.IncrementCount()
  } else {
    term := NewTerm(downcase(word))
    l.terms[term.Text] = term
  }  
}

func (l *Lexicon) SortByText() []*Term {
  terms := l.Terms()
  sort.Sort(ByText(terms))
  return terms
}


func (l *Lexicon) Terms() []*Term {
  terms := make([]*Term, 0)
  for _, v := range l.terms {
    terms = append(terms, v)
  }
  return terms
}


func (l *Lexicon) Size() int {
  return len(l.terms)
}

func (l *Lexicon) Has(word string) (*Term, bool) {
  term, contains := l.terms[downcase(word)]
  return term, contains
}

// Just a lazy way to simplify downcasing terms
func downcase(word string) string {
  return strings.ToLower(word)
}