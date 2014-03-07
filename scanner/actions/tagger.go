package actions

import (
    "bufio"
    "flag"
    "os"
    log "github.com/cihub/seelog"
    "github.com/wwwjscom/ocr_engine/tagger"
    "github.com/wwwjscom/ocr_engine/db"
    "github.com/wwwjscom/ocr_engine/scanner/filewriter"
)

type run_tagger_action struct {
    Args
    lexiconPath *string
    tokens []string
    dbUser *string
    dbPass *string
    dbName *string
    workers *int
    connPool []*db.Mysql
    useStemmer *bool
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
    a.workers = fs.Int("workers", 10, "Number of workers and db connections to make")
    a.useStemmer = fs.Bool("useStemmer", false, "Use a porter stemmer on the terms before assigning them to a group.")
}

func (a *run_tagger_action) Run() {
    SetupLogging(*a.verbosity)

    a.loadTokens()
    log.Info("Tokens loaded")

    log.Info("Filling connection pool")
    a.setupConnPool()

    taggers := new(tagger.Taggers)
    taggers.Init(a.connPool, a.workers, a.useStemmer)
    go taggers.Spawn()

    log.Info("Tagging")
    // For each token, find it in the db
    for i := range a.tokens {
        taggers.Queue <- &a.tokens[i]
    }

    close(taggers.Queue)
    <-taggers.Done

    // Write the missing tokens to disk
    fw := new(filewriter.TrecFileWriter)
    fw.Init("/tmp/missing_tokens")
    go fw.WriteAllTokens()
    for i := range taggers.MissingTokens {
        fw.StringChan<- &taggers.MissingTokens[i]
    }
    close(fw.StringChan)
    fw.Wait()

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

func (a *run_tagger_action) setupConnPool() {
    a.connPool = make([]*db.Mysql, *a.workers)
    for i := 0; i < *a.workers; i++ {
        a.connPool[i] = db.NewMySQLConn(*a.dbUser, *a.dbPass, *a.dbName)
    }
}
