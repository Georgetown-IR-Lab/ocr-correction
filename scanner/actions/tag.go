package actions

import (
    "flag"
    //log "github.com/cihub/seelog"
)

type tagger struct {
    Args
    lexiconPath *string
}

func Tagger() *tagger {
    return new(tagger)
}

func (a *tagger) Name() string {
    return "tagger"
}

func (a *tagger) DefineFlags(fs *flag.FlagSet) {
    a.AddDefaultArgs(fs)

    a.lexiconPath = fs.String("lexicon_path", "",
        "Path and filename to the lexicon file.  One word per line.")
}

func (a *tagger) Run() {
    SetupLogging(*a.verbosity)
    // go
}
