## r2router

A simple router which supports named parameter. Idea for API or backend without any static content. This package contains 2 router types. Router is for pure routing and Seefor supports middleware stacks.

Middlewares are divided into 2 groups, one runs before routing and one runs after routing. Before middleware is thought for serving static, logging, recovery from panic and so on. After middleware is thought for pre-processing data before executing endpoint handler. One can do this by using AppSet method on Params. This mean that Before middlewares are always executed, except when a middleware cancels and does not call next(), meanwhile After middlewares are only call if a route is hit. Each After middleware also has a chance to response and stop calling next().

By default there is no middleware added. See list bellow for suitable middlewares, including recovery which has been moved from this repos.

One interesting feature this package has is measurement of endpoint performance. The timer measures how long time it takes average for each route. The timer itself is a http.Handler so one can use it to serve these statistics locally (see example/timer.go).

There is also a route manager for registering all your routes. This can use for setting up endpoints and also use for building urls.

[![GoDoc](https://godoc.org/github.com/vanng822/r2router?status.svg)](https://godoc.org/github.com/vanng822/r2router)


## Example

### Router
```go
package main

import (
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {
	router := r2router.NewRouter()
	router.Get("/users/:user", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.Write([]byte(p.Get("user")))
	})
	http.ListenAndServe(":8080", router)
}
```
	
Demo: http://premailer.isgoodness.com/
	
### Measuring endpoint performance using Timer

```go
package main

import (
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {
	router := r2router.NewSeeforRouter()
	router.Group("/hello", func(r *r2router.GroupRouter) {
		r.Get("/kitty", func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
			w.Write([]byte("Mau"))
		})

		r.Get("/duck", func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
			w.Write([]byte("Crispy"))
		})

		r.Get("/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			w.Write([]byte(p.Get("name")))
		})
	})
	timer := router.UseTimer(nil)
	
	go http.ListenAndServe("127.0.0.1:8080", router)
	http.ListenAndServe("127.0.0.1:8081", timer)
}
```

Demo: http://premailer.isgoodness.com/timers

### Middleware

```go	
package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
	"log"
	"time"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	// measure time middleware
	seefor.Before(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			handler.ServeHTTP(w, r)
			log.Printf("took: %s", time.Now().Sub(start))
		})
	})
	// set label "say"
	seefor.After(r2router.Wrap(func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		p.AppSet("say", "Hello")
	}))
	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.Write([]byte(fmt.Sprintf("%s %s!", p.AppGet("say").(string), p.Get("name"))))
	})
	
	http.ListenAndServe(":8080", seefor)
}
```	

## Middlewares & Recommended packages

### Middlewares

Panic recovery https://github.com/vanng822/recovery

Access log https://github.com/vanng822/accesslog

Basic auth https://github.com/goji/httpauth

CORS https://github.com/rs/cors

Quick security wins https://github.com/unrolled/secure

HMAC authentication https://github.com/auroratechnologies/vangoh

### Generic packages

Rendering https://github.com/unrolled/render
