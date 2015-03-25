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
	seefor.Use(func(w http.ResponseWriter, r *http.Request, p r2router.Params, next func()) {
		start := time.Now()
		next()
		log.Printf("took: %s", time.Now().Sub(start)) 
	})
	// set label "say"
	seefor.Use(func(w http.ResponseWriter, r *http.Request, p r2router.Params, next func()) {
		p.AppSet("say", "Hello")
		next()
	})
	seefor.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.Write([]byte(fmt.Sprintf("%s %s!", p.AppGet("say").(string), p.Get("name"))))
	})
	
	http.ListenAndServe(":8080", seefor)
}