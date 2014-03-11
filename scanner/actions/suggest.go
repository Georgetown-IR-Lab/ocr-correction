package actions

import (
    "fmt"
    "bufio"
    "flag"
    "os"
    "sync"
    log "github.com/cihub/seelog"
    segments "github.com/wwwjscom/go-segments"
    "github.com/wwwjscom/go-segments/mysql_db"
//    "github.com/wwwjscom/ocr_engine/suggester"
//    "github.com/wwwjscom/ocr_engine/db"
    "github.com/wwwjscom/ocr_engine/scanner/filewriter"
)

type run_suggest_action struct {
    Args
    lexiconPath *string
    tokens []string
    dbUser *string
    dbPass *string
    dbName *string
    workers *int
    topK *int
    connPool chan *mysql_db.Mysql
}

func Suggest() *run_suggest_action {
    return new(run_suggest_action)
}

func (a *run_suggest_action) Name() string {
    return "suggest"
}

func (a *run_suggest_action) DefineFlags(fs *flag.FlagSet) {
    a.AddDefaultArgs(fs)

    a.lexiconPath = fs.String("missing_tokens_path", "/tmp/missing_tokens",
        "Path and filename to the lexicon file.  One word per line.")

    a.dbUser = fs.String("db.user", "", "")
    a.dbPass = fs.String("db.pass", "", "")
    a.dbName = fs.String("db.name", "", "")
    a.workers = fs.Int("workers", 1, "Number of workers and db connections to make")
    a.topK = fs.Int("topk", 10, "Top-k suggested terms to return for each word")
}

func (a *run_suggest_action) Run() {
    wg := new(sync.WaitGroup)
    
    worker_queue := make(chan int, *a.workers)
    for i:=0; i < *a.workers; i++ {
        worker_queue<-1
    }
    
    SetupLogging(*a.verbosity)

    // Write the suggested tokens to disk
    fw := new(filewriter.TrecFileWriter)
    fw.Init("/tmp/suggested_tokens")
    go fw.WriteAllTokens()

    a.loadTokens()
    log.Info("Tokens loaded")

    log.Info("Filling connection pool")
    a.setupConnPool()


    tables_to_search := []string{"names", "geo", "dict"}

    log.Debugf("Tokens size: %d", len(a.tokens))
    for i, token := range a.tokens {
        
        // Sync with the worker queue to prevent a million go threads from being
        // created on startup
        <-worker_queue
        
        wg.Add(1)
        go func(token string, wg *sync.WaitGroup, i int) {
//            ed_suggester := suggester.NewEditDistanceSuggester(a.connPool, tables_to_search)
//            
//            log.Debugf("Finding suggestions for %s", token)
//            suggestions := ed_suggester.Suggest(token)
//
//            suggestion_string := fmt.Sprintf("%s ::: ", token)
//            for _, suggestion := range suggestions {
//                suggestion_string += fmt.Sprintf("%s::%d, ", suggestion.Text, suggestion.Confidence)
//            }
//            
//            // Chop off the ending of ", "
//            suggestion_string = suggestion_string[:len(suggestion_string)-2]
//            
//            fw.StringChan <- &suggestion_string
            
            dbConn := <-a.connPool
            suggestions := segments.Suggest(token, dbConn, tables_to_search)
            a.connPool <- dbConn
            
            suggestion_string := fmt.Sprintf("%s ::: ", token)
            for i, sug := range suggestions {
                suggestion_string += fmt.Sprintf("%s::%f, ", sug.Term, sug.Confidence)
                if i >= *a.topK {
                    break
                }
            }
            // Chop off the ending of ", "
            suggestion_string = suggestion_string[:len(suggestion_string)-2]
            fw.StringChan <- &suggestion_string
            
            log.Infof("%d remaining", len(a.tokens)-i)
            worker_queue<-1
            
            wg.Done()
        }(token, wg, i)
    }

    wg.Wait()
    close(fw.StringChan)
    fw.Wait()
}

func (a *run_suggest_action) loadTokens() {
    a.tokens = make([]string, 0)

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

func (a *run_suggest_action) setupConnPool() {
    a.connPool = make(chan *mysql_db.Mysql, *a.workers)
    for i := 0; i < *a.workers; i++ {
        a.connPool <- segments.NewDBConn(*a.dbUser, *a.dbPass, *a.dbName)
    }
//    a.connPool = make(chan *db.Mysql, *a.workers)
//    for i := 0; i < *a.workers; i++ {
//        a.connPool<- db.NewMySQLConn(*a.dbUser, *a.dbPass, *a.dbName)
//    }
}
