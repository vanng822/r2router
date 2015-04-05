package main

import (
	"github.com/vanng822/r2router"
	"github.com/vanng822/recovery"
	"net/http"
)

func main() {
	seefor := r2router.NewSeeforRouter()
	rec := recovery.NewRecovery()
	rec.PrintStack = true
	seefor.Before(rec.Handler)
	seefor.Before(r2router.WrapBeforeHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Middleware panic")
	})))

	http.ListenAndServe(":8080", seefor)
}
