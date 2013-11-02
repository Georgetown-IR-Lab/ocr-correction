package main

import "github.com/wwwjscom/ocr_engine/scanner/actions"
import "github.com/cwacek/go-subcommand"
import log "github.com/cihub/seelog"

func main() {
	defer log.Flush()
	Run()
}

func Run() {
  actions.SetupLogging(0)
  subcommand.Parse( true, actions.PrintTokens(), actions.Tagger())
}

