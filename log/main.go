package log

import (
	"os"
    "github.com/withmandala/go-log"
)

var Log *log.Logger

func Init(debug bool) {
    Log = log.New(os.Stderr)
    if debug {
        Log.WithDebug()
    }
}
