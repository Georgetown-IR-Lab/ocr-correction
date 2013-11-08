package actions

import (
    "bufio"
    "flag"
    "os"
    log "github.com/cihub/seelog"
    "github.com/wwwjscom/ocr_engine/tagger"
    "github.com/wwwjscom/ocr_engine/db"
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

    a.loadTokens()
    log.Debug("Tokens loaded")

    log.Debug("Connecting to DB")
    mysql := new(db.Mysql)
    mysql =  mysql.New(*a.dbUser, *a.dbPass, *a.dbName)

    taggers := new(tagger.Taggers)
    taggers.Init(mysql)
    go taggers.Spawn()

    log.Debug("Tagging")
    // For each token, find it in the db
    for i := range a.tokens {
        taggers.Queue <- &a.tokens[i]
    }

    close(taggers.Queue)
    <-taggers.Done


    // If not found...
}

func (a *run_tagger_action) loadTokens() {
    a.tokens = make([]string, 10)

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
