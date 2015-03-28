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
