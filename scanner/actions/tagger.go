package actions

import (
    "bufio"
    "flag"
    "os"
    log "github.com/cihub/seelog"
    "github.com/wwwjscom/ocr_engine/tagger"
)

type run_tagger_action struct {
    Args
    lexiconPath *string
    tokens []string
    dbUser *string
    dbPass *string
    dbName *string
}

func Tagger() *run_tagger_action {
    return new(run_tagger_action)
}

func (a *run_tagger_action) Name() string {
    return "tagger"
}

func (a *run_tagger_action) DefineFlags(fs *flag.FlagSet) {
    a.AddDefaultArgs(fs)

    a.lexiconPath = fs.String("lexicon_path", "/tmp/tokens",
        "Path and filename to the lexicon file.  One word per line.")

    a.dbUser = fs.String("db.user", "", "")
    a.dbPass = fs.String("db.pass", "", "")
    a.dbName = fs.String("db.name", "", "")
}

func (a *run_tagger_action) Run() {
    SetupLogging(*a.verbosity)

    log.Debug("Connecting to DB")
    tagger.New(*a.dbUser, *a.dbPass, *a.dbName)
    log.Debug("Connected")
    a.loadTokens()
    log.Debug("Tokens loaded")

    // Open a connection to the db
    // For each token, find it in the db
    // If not found...
}

func (a *run_tagger_action) loadTokens() {
    a.tokens = make([]string, 100000)

    file, err := os.Open(*a.lexiconPath)
    defer file.Close()
    if err != nil {
        log.Criticalf("Error opening lexicon file %s: %s", *a.lexiconPath, err)
        panic(err)
    }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        a.tokens = append(a.tokens, scanner.Text())
    }
}
