package server

import (
	"io"
	"log"
	"os"
)

const (
	infoPrefix  = "[IM] "
	debugPrefix = "[DEBUG] "
	flag        = log.Ldate | log.Ltime | log.Lshortfile
)

var (
	Log   = log.New(os.Stdout, infoPrefix, flag)
	Debug = log.New(io.Discard, debugPrefix, flag)
)

func init() {
	if DebugOn {
		Debug.SetOutput(os.Stdout)
	}
}

const DebugOn = false
