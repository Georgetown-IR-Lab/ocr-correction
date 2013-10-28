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
}

func (fw *TrecFileWriter) Init(filename string) {
    fw.StringChan = make(chan *string)
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
    log.Flush()
    fw.file.Close()
    log.Debugf("Exiting")
}
