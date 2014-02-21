package lexicon

type Term struct {
  Text string
  Count int
}

// Implementing Sort interface based on the Text field
type ByText []*Term
func (b ByText) Len() int           { return len(b) }
func (b ByText) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByText) Less(i, j int) bool { return b[i].Text < b[j].Text }


func NewTerm(text string) *Term {
  t := new(Term)
  t.Text = text
  t.Count = 1
  return t
}

func (t *Term) IncrementCount() int {
  t.Count++
  return t.Count
}

func (t1 *Term) Equal(t2 *Term) bool {
  if t1.Text != t2.Text {
    return false
  } else if t1.Count != t2.Count {
    return false
  }
  
  return true
}