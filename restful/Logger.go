package restful

import (
	"fmt"
	"os"

	l4g "github.com/alecthomas/log4go"
)

type Logging struct {
	file    l4g.Logger
	console l4g.Logger
	isadd   bool //default is false,additional logging

}

func (lg *Logging) FileLogger() l4g.Logger {
	lg.file = l4g.NewLogger()

	var config Config
	config.LoadConfig()
	filepath, ok := config.Log["filepath"]
	if ok == false {
		fmt.Printf("The filepath '%s' for save log is not exist", filepath)
		panic(ok)
	}

	if lg.isadd == true {
		_, ok := os.Stat(filepath)
		if ok == nil {
			os.Truncate(filepath, 0)
		}
	}

	/* Can also specify manually via the following: (these are the defaults) */
	flw := l4g.NewFileLogWriter(filepath, false)
	flw.SetFormat("[%d %T] [%L] (%S) %M")
	flw.SetRotate(false)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(false)
	lg.file.AddFilter("file", l4g.DEBUG, flw)

	return lg.file
}

func (lg *Logging) ConsoleLogger() l4g.Logger {
	lg.console = l4g.NewLogger()
	lg.console.AddFilter("stdout", l4g.FINEST, l4g.NewConsoleLogWriter())
	return lg.console
}

func (lg *Logging) GetLogger() {
	lg.ConsoleLogger()
	lg.FileLogger()
}

func (lg *Logging) Trace(content string) {
	lg.console.Info(content)
	lg.file.Debug(content)
}

func (lg *Logging) Close() {
	lg.file.Close()
	lg.console.Close()
}
