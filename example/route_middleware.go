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
