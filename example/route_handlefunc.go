package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

// Wrapper for http.HandlerFunc
func RouteHandleFunc(next http.HandlerFunc) r2router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ r2router.Params) {
		next(w, r)
	}
}

func main() {
	seefor := r2router.NewSeeforRouter()
	seefor.Get("/hello/world", RouteHandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	}))
	http.ListenAndServe("127.0.0.1:8080", seefor)
}
