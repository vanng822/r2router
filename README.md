## r2router
A simple router which supports named parameter. Idea for API or backend without any static content.

The idea of "middleware" here is for pre-processing data before executing handler. This means that they are always executed after routing and before handler. If no route is matched then no middleware is executed.

[![GoDoc](https://godoc.org/github.com/vanng822/r2router?status.svg)](https://godoc.org/github.com/vanng822/r2router)


## Example

### Router

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

### Middleware
	
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
		seefor.Before(r2router.BeforeFunc(func(w http.ResponseWriter, r *http.Request, next func()) {
			start := time.Now()
			next()
			log.Printf("took: %s", time.Now().Sub(start)) 
		}))
		// set label "say"
		seefor.After(r2router.AfterFunc(func(w http.ResponseWriter, r *http.Request, p r2router.Params, next func()) {
			p.AppSet("say", "Hello")
			next()
		}))
		seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			w.Write([]byte(fmt.Sprintf("%s %s!", p.AppGet("say").(string), p.Get("name"))))
		})
		
		http.ListenAndServe(":8080", seefor)
	}
	
### Measuring endpoint performance using Timer

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
		