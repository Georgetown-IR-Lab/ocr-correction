package suggester

type Suggester interface {
  Suggest(string) []*Suggestion
}