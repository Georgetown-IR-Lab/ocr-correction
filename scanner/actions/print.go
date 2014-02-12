package actions

import "fmt"
import "flag"
import "os"
import "path/filepath"
import "regexp"
import log "github.com/cihub/seelog"
import filereader "github.com/wwwjscom/ocr_engine/scanner/filereader"
import filewriter "github.com/wwwjscom/ocr_engine/scanner/filewriter"
import lexicon "github.com/wwwjscom/ocr_engine/lexicon"

func PrintTokens() *print_tokens_action {
    return new(print_tokens_action)
}

type print_tokens_action struct {
    Args

    workers chan string
    output chan string
    worker_count int
    tokenOutputPath *string
}

func (a *print_tokens_action) Name() string {
    return "print_tokens"
}

func (a *print_tokens_action) DefineFlags(fs *flag.FlagSet) {
    a.AddDefaultArgs(fs)

    a.docroot = fs.String("doc.root", "",
    `The root directory under which to find document`)

    a.docpattern = fs.String("doc.pattern", `^[^\.].+`,
    `A regular expression to match document names`)

    a.tokenOutputPath = fs.String("token_output_path", "/tmp/tokens",
        "Path and filename to output the (non-unique, non-sorted) tokens")
}


func (a *print_tokens_action) Run() {
    SetupLogging(*a.verbosity)
    
    lex := lexicon.NewLexicon()

    writer := new(filewriter.TrecFileWriter)
    writer.Init(*a.tokenOutputPath)
    go writer.WriteAllTokens()

    docStream := make(chan filereader.Document)

    walker := new(DocWalker)
    walker.WalkDocuments(*a.docroot, *a.docpattern, docStream)

    for doc := range docStream {

        for t := range doc.Tokens() {
            log.Tracef("Adding token: %s", t)            
            lex.Add(t.Text)            
        }

        log.Debugf("Document %s (%d tokens)\n", doc.Identifier(), doc.Len())
    }
    
    log.Debug("Sorting Words")
    lex.SortByText()
    
    log.Debug("Writing Words")
    // Write out all the words
    for _, t := range lex.Terms() {
      writer.StringChan <- &t.Text
    }
    
    log.Info("Done reading from the docStream")
    close(writer.StringChan)

    // Wait for the writer to finish
    writer.Wait()
}


type DocWalker struct {
  output  chan filereader.Document
  workers  chan string
  worker_count int
  filepattern string
}

func (d *DocWalker) WalkDocuments(docroot, pattern string, out chan filereader.Document) {

  d.output = out
  d.workers = make(chan string)
  d.worker_count = 0
  d.filepattern = pattern

  fmt.Println("Reading in documents")
  filepath.Walk(docroot, d.read_file)

  go d.signal_when_done()
}

func (d *DocWalker) signal_when_done() {
  for {
    select {
    case file := <- d.workers:
      d.worker_count -= 1
      log.Infof("Worker for %s done. Waiting for %d workers.", file, d.worker_count)
      if d.worker_count <= 0 {
          fmt.Println("Finished reading documents")
          close(d.output)
        return
      }
    }
  }
}

func (d *DocWalker) read_file( path string, info os.FileInfo, err error) error {

  if info.Mode().IsRegular() {
    file := filepath.Base(path)

    log.Debugf("Trying file %s", file)

    matched, err := regexp.MatchString(d.filepattern, file);
    log.Debugf("File match: %v, error: %v", matched, err)
    if matched && err == nil {

      fr := new(filereader.TrecFileReader)
      fr.Init(path)

      go func() {
          for doc := range fr.ReadAll() {
            d.output <- doc
          }
          d.workers <- fr.Path()
          return
      }()

      d.worker_count += 1
      /*log.Errorf("Now have %d workers", d.worker_count)*/
    }
  }
  return nil
}
