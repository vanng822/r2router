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
	seefor.Before(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			handler.ServeHTTP(w, r)
			log.Printf("took: %s", time.Now().Sub(start))
		})
	})
	// set label "say"
	seefor.After(r2router.AfterFunc(func(w http.ResponseWriter, r *http.Request, p r2router.Params, next func()) {
		p.AppSet("say", "Hello")
		next()
	}))
	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.Write([]byte(fmt.Sprintf("%s %s!", p.AppGet("say").(string), p.Get("name"))))
	})

	timer := seefor.UseTimer(nil)

	go http.ListenAndServe("127.0.0.1:8080", seefor)
	http.ListenAndServe("127.0.0.1:8081", timer)
}
