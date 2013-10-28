package logging

import "fmt"
import log "github.com/cihub/seelog"
import "testing"

func SetupTestLogging() {
	var appConfig = `
  <seelog type="sync" minlevel='%s'>
  <outputs formatid="scanner">
    <filter levels="critical,error,warn,info">
      <console formatid="scanner" />
    </filter>
    <filter levels="debug,trace">
      <console formatid="debug" />
    </filter>
  </outputs>
  <formats>
  <format id="scanner" format="test: [%%LEV] %%Msg%%n" />
  <format id="debug" format="test: [%%LEV] %%Func :: %%Msg%%n" />
  </formats>
  </seelog>
`

var config string
if testing.Verbose() {
  config =  fmt.Sprintf(appConfig, "trace")
} else {
  config =  fmt.Sprintf(appConfig, "info")
}

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))

	if err != nil {
		fmt.Println(err)
		return
	}

	log.ReplaceLogger(logger)
}
