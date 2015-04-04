package main

import (
	"github.com/goji/httpauth"
	"github.com/vanng822/r2router"
	"net/http"
)

// Example how to use basic auth together with r2router
func main() {
	router := r2router.NewSeeforRouter()
	
	// basic auth for entire router
	router.Before(httpauth.SimpleBasicAuth("testuser", "testpw"))
	
	router.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		w.Write([]byte( p.Get("name")))
	})

	http.ListenAndServe("127.0.0.1:8080", router)
}