package lexicon

import (
  "testing"
)

func TestAddWord(t *testing.T) {
  lex := NewLexicon()
  
  if lex.Size() != 0 {
    t.Error("Words is not of expected length")
  }
  
  lex.Add("Foo")
  
  if lex.Size() != 1 {
    t.Error("Words is not of expected length")    
  }  
}

func TestHasWord(t *testing.T) {
  
  tests := []struct {
    word string
    expected *Term
  }{
    { "Foo", NewTerm("foo") },
    { "foo", NewTerm("foo") },
    { "Bar", nil },
  }
  
  lex := NewLexicon()
  lex.Add("Foo")
  
  for _, tc := range tests {
    if term, _ := lex.Has(tc.word); term == tc.expected {
      continue
    } else if term, has := lex.Has(tc.word); has && !term.Equal(tc.expected) {
      t.Errorf("Given %s, expected %s, but received %s", tc.word, tc.expected, term)
    }
  }
}

func TestSortByText(t *testing.T) {
  lex := NewLexicon()
  lex.Add("ZZZ")
  lex.Add("AAA")
  lex.Add("AAA")
  result := lex.SortByText()
  
  if result[0].Text != "aaa" || result[1].Text != "zzz" {
    t.Error("Not sorted correctly...")
  }
}

func TestUniquiness(t *testing.T) {
  lex := NewLexicon()
  lex.Add("ZZZ")
  lex.Add("ZZZ")
  lex.Add("AAA") 
  lex.Add("AAA") 
  lex.Add("aaa")
  
  if lex.Size() != 2 {
    t.Errorf("Expected size of 2, found %d", lex.Size())
  }
  
  if term, has := lex.Has("aaa"); has && term.Count != 3 {
    t.Errorf("Expected 3, got %d", term.Count)
  }

  if term, has := lex.Has("ZzZ"); has && term.Count != 2 {  
    t.Errorf("Expected 2, got %d", term.Count)
  }
}

func TestTerms(t *testing.T) {
  lex := NewLexicon()
  lex.Add("ZZZ")
  lex.Add("AAA") 
  lex.Add("AAA") 

  terms := lex.terms
  if len(terms) != 2 {
    t.Errorf("Wrong size")
  }
}