package db

type Lexicon interface {
    Init(DB)
    Find(*string) bool
}
