// Package r2router provide a simple router. Suitable for API where you could map parameter from the url.
//
//	package main
// 
//	import (
//		"github.com/vanng822/r2router"
//		"net/http"
//	)
//
//
//	func main() {
//		router := r2router.NewRouter()
// 		router.Get("/users/:user", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
//			w.Write([]byte(p.Get("user")))
//		})
//		http.ListenAndServe(":8080", router)
//	}
package r2router