package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"log"
	"net/http"
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

	timer := seefor.UseTimer(nil)

	go http.ListenAndServe("127.0.0.1:8080", seefor)
	http.ListenAndServe("127.0.0.1:8081", timer)
}
