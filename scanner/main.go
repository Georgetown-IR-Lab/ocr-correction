package main

import "github.com/wwwjscom/ocr_engine/scanner/actions"
import "github.com/cwacek/go-subcommand"
import log "github.com/cihub/seelog"
import "runtime"

func main() {
	defer log.Flush()
    runtime.GOMAXPROCS(runtime.NumCPU())
	Run()
}

func Run() {
  actions.SetupLogging(0)
  subcommand.Parse( true, actions.PrintTokens(), actions.Tagger())
}

