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

type RecoveryOptions struct {
	Logger     Logger
	StackAll   bool
	StackSize  int
	PrintStack bool
}

type Recovery struct {
	options *RecoveryOptions
	next http.Handler
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, rec.options.StackSize)
			stack = stack[:runtime.Stack(stack, rec.options.StackAll)]
			format := "PANIC: %s\n%s"
			rec.options.Logger.Printf(format, err, stack)

			if rec.options.PrintStack {
				fmt.Fprintf(w, format, err, stack)
			} else {
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}
	}()
	rec.next.ServeHTTP(w, req)
}

func NewRecoveryOptions() *RecoveryOptions {
	return &RecoveryOptions{
		Logger:     log.New(os.Stdout, "[seefor] ", 0),
		StackAll:   false,
		StackSize:  1024 * 8,
		PrintStack: false,
	}
}

func NewRecovery(options *RecoveryOptions) Before {
	if options == nil {
		options = NewRecoveryOptions()
	}
	return func(next http.Handler) http.Handler {
		rec := &Recovery{
			options: options,
			next: next,
		}
		return rec
	}
}
