## r2router

A simple router which supports named parameter. Idea for API or backend without any static content. This package contains 2 router types. Router is for pure routing and Seefor supports middleware stacks.

Middlewares are divided into 2 groups, one runs before routing and one runs after routing. Before middleware is thought for serving static, logging, recovery from panic and so on. After middleware is thought for pre-processing data before executing endpoint handler. One can do this by using AppSet method on Params. This mean that Before middlewares are always executed, except when a middleware cancels and does not call next(), meanwhile After middlewares are only call if a route is hit. Each After middleware also has a chance to response and stop calling next().

By default there is no middleware added. See list bellow for suitable middlewares, including recovery which has been moved from this repos.

One interesting feature this package has is measurement of endpoint performance. The timer measures how long time it takes average for each route. The timer itself is a http.Handler so one can use it to serve these statistics locally (see example/timer.go).

There is also a route manager for registering all your routes. This can use for setting up endpoints and also use for building urls.

[![GoDoc](https://godoc.org/github.com/vanng822/r2router?status.svg)](https://godoc.org/github.com/vanng822/r2router)
[![Build Status](https://travis-ci.org/vanng822/r2router.svg?branch=master)](https://travis-ci.org/vanng822/r2router)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/vanng822/r2router) [![](http://gocover.io/_badge/github.com/vanng822/r2router)](http://gocover.io/github.com/vanng822/r2router)

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

If your routes don't contain named params and you have existing http.HandlerFunc then you can wrap as bellow

```go

package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

// Wrapper for http.HandlerFunc, similar can be done for http.Handler
func RouteHandlerFunc(next http.HandlerFunc) r2router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
		next(w, r)
	}
}

func main() {
	seefor := r2router.NewSeeforRouter()
	seefor.Get("/hello/world", RouteHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	}))
	http.ListenAndServe("127.0.0.1:8080", seefor)
}
```
	
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
	seefor.Before(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("took: %s", time.Now().Sub(start))
		})
	})
	// set label "say"
	seefor.After(func(next r2router.Handler) r2router.Handler {
		return r2router.HandlerFunc(func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			p.AppSet("say", "Hello")
			next.ServeHTTP(w, r, p)
		})
	})

	seefor.After(r2router.Wrap(func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		p.AppSet("goodbye", "Bye bye")
	}))

	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "%s %s!\n%s", p.AppGet("say").(string), p.Get("name"), p.AppGet("goodbye"))
	})
	
	http.ListenAndServe(":8080", seefor)
}
```	

If you want to add middlewares for a specific route then you should create your own wrapper. This way you will have full control over your middlewares.
You could do something like bellow

```go

package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

// Your own route middle wrapper
func RouteMiddleware(next r2router.HandlerFunc) r2router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		p.AppSet("say", "Hello")
		next(w, r, p)
	}
}

func main() {
	seefor := r2router.NewSeeforRouter()
	seefor.Get("/hello/:name", RouteMiddleware(func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "%s %s!", p.AppGet("say").(string), p.Get("name"))
	}))
	http.ListenAndServe("127.0.0.1:8080", seefor)
}
```



### Route manager

```go	
package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	rmanager := r2router.NewRouteManager()
	// register and use it at the same time
	seefor.Get(rmanager.Add("hello::name", "/hello/:name"), func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "Hello %s!", p.Get("name"))
	})

	// Or register first elsewhere and get it later
	rmanager.Add("redirect::name", "/redirect/:name")

	seefor.Get(rmanager.PathFor("redirect::name"), func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		// Building url for routename "hello::name" and redirect
		http.Redirect(w, r, rmanager.UrlFor("hello::name", r2router.P{"name": []string{p.Get("name")}}), http.StatusFound)
	})

	http.ListenAndServe("127.0.0.1:8080", seefor)
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

Static content https://github.com/hypebeast/gojistatic

### Generic packages

Rendering https://github.com/unrolled/render

Session https://github.com/gorilla/sessions

HTTP request data binding & validation https://github.com/mholt/binding


