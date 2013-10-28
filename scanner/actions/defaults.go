package actions

import "flag"
import "os"
import "path/filepath"
import "regexp"
import "fmt"
import log "github.com/cihub/seelog"
import filereader "github.com/wwwjscom/ocr_engine/scanner/filereader"

type Args struct {
  docroot *string
  docpattern *string
  verbosity *int
}

func (a *Args) AddDefaultArgs(fs *flag.FlagSet) {

    a.docroot = fs.String("doc.root", "",
    `The root directory under which to find document`)

    a.docpattern = fs.String("doc.pattern", `^[^\.].+`,
    `A regular expression to match document names`)

    a.verbosity = fs.Int("v", 0, "Be verbose [1, 2, 3]")
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

var appConfig = `
  <seelog minlevel='%s'>
  <outputs formatid="scanner">
    <filter levels="critical,error,warn,info">
      <console formatid="scanner" />
    </filter>
    <filter levels="debug,trace">
      <console formatid="debug" />
    </filter>
  </outputs>
  <formats>
  <format id="scanner" format="[%%Time]:%%LEVEL:: %%Msg%%n" />
  <format id="debug" format="[%%Time]:%%LEVEL:%%Func:: %%Msg%%n" />
  </formats>
`

var config string

func SetupLogging(verbosity int) {

  switch verbosity {
  case 0:
    fallthrough
  case 1:
    config = fmt.Sprintf(appConfig, "warn")
  case 2:
    config = fmt.Sprintf(appConfig, "info")
  case 3:
    config = fmt.Sprintf(appConfig, "debug")
  default:
    config = fmt.Sprintf(appConfig, "trace")
  }

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))

	if err != nil {
		fmt.Println(err)
		return
	}

	log.ReplaceLogger(logger)
}

