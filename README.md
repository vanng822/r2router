## r2router
A simple router which supports named parameter

## Example

	package main

	import (
		"github.com/vanng822/r2router"
		"net/http"
	)
	
	func main() {
		router := r2router.NewRouter()
		router.Get("/users/:user", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			w.Write([]byte(p["user"]))
		})
		http.ListenAndServe(":8080", router)
	}
