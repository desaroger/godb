package logs

import (
	"fmt"
	"log"
	"os"
)

var _info *log.Logger
var _error *log.Logger

func Initialize() {
	_info = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
}

func Info(format string, a ...any) {
	usrMsg := fmt.Sprintf(format, a...)
	_info.Println(usrMsg)
}

func Error(err error, format string, a ...any) {
	usrMsg := fmt.Sprintf(format, a...)

	_error.Printf("%s: %v\n", usrMsg, err)
}
