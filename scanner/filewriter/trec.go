package filewriter

import (
    "os"
    "fmt"
    log "github.com/cihub/seelog"
)

type TrecFileWriter struct {
    filename string
    StringChan chan *string
    file *os.File
    Done chan bool
}

func (fw *TrecFileWriter) Init(filename string) {
    fw.StringChan = make(chan *string, 100000)
    fw.Done = make(chan bool)
    fw.filename = filename

    if file, err := os.Create(fw.filename); err != nil {
        panic(fmt.Sprintf("Unable to open file %s", fw.filename))
    } else {
        fw.file = file
    }

}

func (fw *TrecFileWriter) WriteAllTokens() {
    log.Debugf("Monitoring the writer channel")
    for t := range fw.StringChan {
        log.Debugf("Received %s. Writing it out to disk.", *t)
        fw.file.WriteString(*t + "\n")
    }
    fw.file.Close()
    log.Info("Sending exit signal")
    fw.Done<- true
    log.Info("Writer, out!")
}

func (fw *TrecFileWriter) Wait() {
    log.Info("Waiting on signal")
    <-fw.Done
    log.Info("Received done signal")
}
