package main

import (
	"fmt"
	"github.com/vanng822/r2router"
	"net/http"
)

func main() {

	routes := [][]string{
		{"GET", "/authorizations"},
		{"GET", "/authorizations/:id"},
		{"POST", "/authorizations"},
		{"DELETE", "/authorizations/:id"},

		// Activity
		{"GET", "/events"},
		{"GET", "/repos/:owner/:repo/events"},
		{"GET", "/networks/:owner/:repo/events"},
		{"GET", "/orgs/:org/events"},
		{"GET", "/users/:user/received_events"},
		{"GET", "/users/:user/received_events/public"},
		{"GET", "/users/:user/events"},
		{"GET", "/users/:user/events/public"},
		{"GET", "/users/:user/events/orgs/:org"},
		{"GET", "/feeds"},
		{"GET", "/notifications"},
		{"GET", "/repos/:owner/:repo/notifications"},
		{"PUT", "/notifications"},
		{"PUT", "/repos/:owner/:repo/notifications"},
		{"GET", "/notifications/threads/:id"},
		{"GET", "/notifications/threads/:id/subscription"},
		{"PUT", "/notifications/threads/:id/subscription"},
		{"DELETE", "/notifications/threads/:id/subscription"},
		{"GET", "/repos/:owner/:repo/stargazers"},
		{"GET", "/users/:user/starred"},
		{"GET", "/user/starred"},
		{"GET", "/user/starred/:owner/:repo"},
		{"PUT", "/user/starred/:owner/:repo"},
		{"DELETE", "/user/starred/:owner/:repo"},
		{"GET", "/repos/:owner/:repo/subscribers"},
		{"GET", "/users/:user/subscriptions"},
		{"GET", "/user/subscriptions"},
		{"GET", "/repos/:owner/:repo/subscription"},
		{"PUT", "/repos/:owner/:repo/subscription"},
		{"DELETE", "/repos/:owner/:repo/subscription"},
		{"GET", "/user/subscriptions/:owner/:repo"},
		{"PUT", "/user/subscriptions/:owner/:repo"},
		{"DELETE", "/user/subscriptions/:owner/:repo"},

		// Users
		{"GET", "/users/:user"},
		{"GET", "/user"},
		{"GET", "/users"},
		{"GET", "/user/emails"},
		{"POST", "/user/emails"},
		{"DELETE", "/user/emails"},
		{"GET", "/users/:user/followers"},
		{"GET", "/user/followers"},
		{"GET", "/users/:user/following"},
		{"GET", "/user/following"},
		{"GET", "/user/following/:user"},
		{"GET", "/users/:user/following/:target_user"},
		{"PUT", "/user/following/:user"},
		{"DELETE", "/user/following/:user"},
		{"GET", "/users/:user/keys"},
		{"GET", "/user/keys"},
		{"GET", "/user/keys/:id"},
		{"POST", "/user/keys"},
		{"DELETE", "/user/keys/:id"},
	}

	router := r2router.NewRouter()
	
	for _, route := range routes {
		switch route[0] {
		case "GET":
			router.Get(route[1], func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
				w.Write([]byte(fmt.Sprintf("#v", p)))
			})

		case "POST":
			router.Post(route[1], func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
				w.Write([]byte(fmt.Sprintf("#v", p)))
			})

		case "DELETE":
			router.Delete(route[1], func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
				w.Write([]byte(fmt.Sprintf("#v", p)))
			})

		case "PUT":
			router.Put(route[1], func(w http.ResponseWriter, r *http.Request, p r2router.Params) {
				w.Write([]byte(fmt.Sprintf("#v", p)))
			})
		}

	}

	http.ListenAndServe(":8080", router)
}
