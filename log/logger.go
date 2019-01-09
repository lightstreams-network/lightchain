package log

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/ethereum/go-ethereum/log"
	tmtLog "github.com/tendermint/tendermint/libs/log"
)

var glogger *log.GlogHandler

func init() {
	//usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	usecolor := true
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger = log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))
}

// Setup sets up the logging infrastructure
func SetupLogger(lvl int) error {
	glogger.Verbosity(log.Lvl(lvl))
	log.Root().SetHandler(glogger)
	return nil
}

// Interface assertions
var _ tmtLog.Logger = (*logger)(nil)

// ---------------------------
// Logger - wraps the logger in tmlibs

type logger struct {
	keyvals []interface{}
}

// Logger returns a new instance of an lightchain logger. With() should
// be called upon the returned instance to set default keys
// #unstable
func NewLogger() tmtLog.Logger {
	logger := logger{keyvals: make([]interface{}, 0)}
	return logger
}

// Debug proxies everything to the go-ethereum logging facilities
// #unstable
func (l logger) Debug(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Debug(msg, ctx...)
}

// Info proxies everything to the go-ethereum logging facilities
// #unstable
func (l logger) Info(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Info(msg, ctx...)
}

// Error proxies everything to the go-ethereum logging facilities
// #unstable
func (l logger) Error(msg string, ctx ...interface{}) {
	ctx = append(l.keyvals, ctx...)
	log.Error(msg, ctx...)
}

// With proxies everything to the go-ethereum logging facilities
// #unstable
func (l logger) With(ctx ...interface{}) tmtLog.Logger {
	l.keyvals = append(l.keyvals, ctx...)
	return l
}
