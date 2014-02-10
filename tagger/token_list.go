package tagger

func (tm *TokenMap) Init() *TokenMap {
    tm = new(TokenMap)
    tm.tokenMap = make(map[*string]int)
    return tm
}

func (tm *TokenMap) Add(token *string) {
    tm.tokenMap[token] += 1
}

func (tm *TokenMap) Len() int {
    return len(tm.tokenMap)
}

func (tm *TokenMap) Map() map[*string]int {
    return tm.tokenMap
}
