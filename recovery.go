package r2router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Recovery struct {
	Logger     Logger
	StackAll   bool
	StackSize  int
	PrintStack bool
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, req *http.Request, next func()) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			format := "PANIC: %s\n%s"
			rec.Logger.Printf(format, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(w, format, err, stack)
			} else {
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}
	}()

	next()
}

func NewRecovery() *Recovery {
	rec := &Recovery{
		Logger:     log.New(os.Stdout, "[seefor] ", 0),
		StackAll:   false,
		StackSize:  1024 * 8,
		PrintStack: false,
	}
	return rec
}
