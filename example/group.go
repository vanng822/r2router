package main

import (
	"fmt"
	"github.com/unrolled/render"
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {
	renderer := render.New()
	router := r2router.NewRouter()
	router.Group("/repos/:owner/:repo", func(r *r2router.GroupRouter) {
		r.Get("/stats/contributors", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			data := r2router.M{
				"owner": p.Get("owner"),
				"repo":  p.Get("repo"),
			}
			renderer.JSON(w, http.StatusOK, data)
		})

		r.Get("/releases/:id/assets", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			w.Write([]byte(fmt.Sprintf("#v", p)))
		})

		r.Get("/issue", func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
			w.Write([]byte(fmt.Sprintf("#v", p)))
		})
	})

	http.ListenAndServe(":8080", router)
}
