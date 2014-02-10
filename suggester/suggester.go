package suggester

type Suggester interface {
    Init(lexicons ...Lexicon)
    Suggest(*string) []string
}
