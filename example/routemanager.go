package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	rmanager := r2router.NewRouteManager()
	// register and use it at the same time
	seefor.Get(rmanager.Add("hello::name", "/hello/:name"), func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		fmt.Fprintf(w, "Hello %s!", p.Get("name"))
	})

	// Or register first elsewhere and get it later
	rmanager.Add("redirect::name", "/redirect/:name")

	seefor.Get(rmanager.PathFor("redirect::name"), func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
		// Building url for routename "hello::name" and redirect
		http.Redirect(w, r, rmanager.UrlFor("hello::name", r2router.P{"name": []string{p.Get("name")}}), http.StatusFound)
	})

	http.ListenAndServe("127.0.0.1:8080", seefor)
}
