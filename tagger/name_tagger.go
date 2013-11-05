package tagger

import (
    //"os"
    log "github.com/cihub/seelog"
    "github.com/wwwjscom/ocr_engine/db"
)

type Taggers struct {
    MAX_WORKERS int
    Queue chan *string // Tokens channel to process
    mysql chan *db.Mysql // database conn
    complete chan int // tracks complete workers
    Done chan bool // Signals that all workers are complete
}

func (t *Taggers) Init(conn *db.Mysql) {
    t.MAX_WORKERS = 10
    t.complete = make(chan int)
    t.Done = make(chan bool)
    t.Queue = make(chan *string)
    t.mysql = make(chan *db.Mysql)

    // Don't block waiting for channel to be read
    go func() { t.mysql<- conn }() 
}

// Spawn the tagger workers with a shared mysql conn
func (t *Taggers) Spawn() {
    for i := 0; i < t.MAX_WORKERS; i++ {
        go t.find(t.Queue, t.mysql)
    }

    t.wait_on_workers()
    t.Done<- true
}

// find all tokens in the channel until it's closed
func (t *Taggers) find(queue chan *string, mysql chan *db.Mysql) {
    i := 0
    for token := range queue {
        log.Tracef("Searching for token %s", *token)
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
        if done == t.MAX_WORKERS {
            log.Debugf("%d tokens processed", processed)
            return
        }
    }
}
