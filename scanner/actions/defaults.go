package actions

import "flag"
import "fmt"
import log "github.com/cihub/seelog"

type Args struct {
  docroot *string
  docpattern *string
  verbosity *int
}

func (a *Args) AddDefaultArgs(fs *flag.FlagSet) {

    a.verbosity = fs.Int("v", 0, "Be verbose (less) [1, 2, 3] (more)")
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

