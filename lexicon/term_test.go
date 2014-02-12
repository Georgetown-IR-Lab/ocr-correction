package lexicon

import "testing"

func TestNewTerm(t *testing.T) {
  term := NewTerm("foo")
  if term.Text != "foo" {
    t.Error("Unexpected text returned")
  }
  
  if term.Count != 1 {
    t.Errorf("Count was expected to be 1, but was %d", term.Count)
  }
}

func TestIncrementCount(t *testing.T) {
  term := NewTerm("foo")
  if term.IncrementCount() != 2 {
    t.Errorf("Expected count to be 2, but was %d", term.Count)
  }
}

func TestEqual(t *testing.T) {
  tests := []struct {
    term1 *Term
    term2 *Term
    equal bool
  }{
    { NewTerm("Foo"), NewTerm("Foo"), true},
    { NewTerm("Foo"), NewTerm("Bar"), false},
  }
  
  for _, tc := range tests {
    if tc.equal && tc.term1.Equal(tc.term2) {
      // all good -- expected to be equal and it is
    } else if !tc.equal && !tc.term1.Equal(tc.term2) {
      // all good -- expect to be unequal and it is
    } else {
      t.Errorf("Expected terms to be equal (%s): %s %s", tc.equal, tc.term1, tc.term2)
    }
  }
}