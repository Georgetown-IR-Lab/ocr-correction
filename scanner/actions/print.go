package actions

import "fmt"
import "flag"
import log "github.com/cihub/seelog"
import filereader "github.com/wwwjscom/ocr_engine/scanner/filereader"
import filewriter "github.com/wwwjscom/ocr_engine/scanner/filewriter"


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

    a.tokenOutputPath = fs.String("token_output_path", "/tmp/tokens",
        "Path and filename to output the (non-unique, non-sorted) tokens")
}


func (a *print_tokens_action) Run() {
    SetupLogging(*a.verbosity)

    writer := new(filewriter.TrecFileWriter)
    writer.Init(*a.tokenOutputPath)
    go writer.WriteAllTokens()

    docStream := make(chan filereader.Document)

    walker := new(DocWalker)
    walker.WalkDocuments(*a.docroot, *a.docpattern, docStream)

    for doc := range docStream {

        for t := range doc.Tokens() {
            log.Debugf("Adding token: %s", t)
            writer.StringChan <- &t.Text
        }

        fmt.Printf("Document %s (%d tokens)\n", doc.Identifier(), doc.Len())
    }
    close(writer.StringChan)
}

