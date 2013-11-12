package tagger

import (
    log "github.com/cihub/seelog"
    "github.com/wwwjscom/ocr_engine/db"
    "github.com/wwwjscom/go-sutils"
)

type Taggers struct {
    workers *int
    Queue chan *string // Tokens channel to process
    mysql chan *db.Mysql // database conn
    complete chan int // tracks complete workers
    Done chan bool // Signals that all workers are complete

    DictTokens []string
    NamesTokens []string
    GeoTokens []string
    MissingTokens []string
}

func (t *Taggers) Init(conns []*db.Mysql, workers *int) {
    t.workers = workers
    t.complete = make(chan int)
    t.Done = make(chan bool)
    t.Queue = make(chan *string)
    t.mysql = make(chan *db.Mysql)

    t.DictTokens = make([]string, 0, 100)
    t.NamesTokens = make([]string, 0, 100)
    t.GeoTokens = make([]string, 0, 100)
    t.MissingTokens = make([]string, 0, 100)

    // Don't block waiting for channel to be read
    go func() {
        for _, conn := range conns {
            t.mysql<- conn
        }
    }()
}

// Spawn the tagger workers with shared mysql conns
func (t *Taggers) Spawn() {
    for i := 0; i < *t.workers; i++ {
        go t.find(t.Queue, t.mysql)
    }

    t.wait_on_workers()
    t.Done<- true
}

// find all tokens in the channel until it's closed
func (t *Taggers) find(queue chan *string, mysql chan *db.Mysql) {
    i := 0
    for token := range queue {
        if len(*token) == 0 {
            log.Tracef("Caught empty str")
            continue
        }
        conn := <-mysql
        log.Tracef("Searching for token %s", *token)
        kind := t.search_all_tables(token, conn)
        go func() { mysql<- conn }()

        switch kind {
        case -1: t.MissingTokens = append(t.MissingTokens, *token)
        case 1: t.NamesTokens = append(t.NamesTokens, *token)
        case 2: t.DictTokens = append(t.DictTokens, *token)
        case 3: t.GeoTokens = append(t.GeoTokens, *token)
        }

        i++
    }

    log.Debugf("Worker, out.  Processed %d tokens", i)
    t.complete<- i
}

// Don't return until all workers have exited
func (t *Taggers) wait_on_workers() {
    done := 0
    processed := 0
    for count := range t.complete {
        done++
        processed += count
        if done == *t.workers {
            log.Debugf("%d tokens processed", processed)
            log.Debugf("%d missing tokens", len(t.MissingTokens))
            log.Debugf("%d names tokens", len(t.NamesTokens))
            log.Debugf("%d dict tokens", len(t.DictTokens))
            log.Debugf("%d geo tokens", len(t.GeoTokens))
            return
        }
    }
}

// Searches all tables unit it finds a match or no tables are left
func (t *Taggers) search_all_tables(token *string, conn *db.Mysql) int {
    escaped_token := sutils.EscapeAllQuotes(*token)
    names_q := "select * from names WHERE name = \"" + escaped_token + "\""
    dict_q  := "select * from dict WHERE word = \"" + escaped_token + "\""
    geo_q   := "select * from geo WHERE name = \"" + escaped_token + "\""

    if conn.Query(names_q) != nil {
        log.Tracef("%s found in names table", *token)
        return 1
    } else if conn.Query(dict_q) != nil {
        log.Tracef("%s found in dict table", *token)
        return 2
    } else if conn.Query(geo_q) != nil {
        log.Tracef("%s found in geo table", *token)
        return 3
    }
    log.Tracef("%s not found", *token)
    return -1
}
