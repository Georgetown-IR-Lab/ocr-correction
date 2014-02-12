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
    if tc.expected == lex.Has(tc.word) {
      continue
    } else if !lex.Has(tc.word).Equal(tc.expected) {
      t.Errorf("Given %s, expected %s, but received %s", tc.word, tc.expected, lex.Has(tc.word))
    }
  }
}

func TestSortByText(t *testing.T) {
  lex := NewLexicon()
  lex.Add("ZZZ")
  lex.Add("AAA")
  lex.SortByText()
  
  if lex.terms[0].Text != "aaa" || lex.terms[1].Text != "zzz" {
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
  
  if count := lex.Has("aaa").Count; count != 3 {
    t.Errorf("Expected 3, got %d", count)
  }
  
  if count := lex.Has("ZzZ").Count; count != 2 {
    t.Errorf("Expected 2, got %d", count)
  }
}