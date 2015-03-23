package r2router

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	HTTP_METHOD_GET     = "GET"
	HTTP_METHOD_POST    = "POST"
	HTTP_METHOD_DELETE  = "DELETE"
	HTTP_METHOD_OPTIONS = "OPTIONS"
	HTTP_METHOD_HEAD    = "HEAD"
	HTTP_METHOD_PUT     = "PUT"
)

type Params map[string]string

type Handler func(http.ResponseWriter, *http.Request, Params)

type Router struct {
	roots               map[string]*rootNode
	optionsAllowMethods []string
}

func NewRouter() *Router {
	r := &Router{}
	r.roots = make(map[string]*rootNode)
	r.optionsAllowMethods = []string{
		HTTP_METHOD_GET,
		HTTP_METHOD_POST,
		HTTP_METHOD_PUT,
		HTTP_METHOD_DELETE,
		HTTP_METHOD_HEAD,
	}
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	if root, exist := r.roots[req.Method]; exist {
		handler, params := root.match(req.URL.Path)
		if handler != nil {
			handler(w, req, params)
			log.Println(time.Now().Sub(now))
			return
		}
	}

	// if options find handler for different method
	if req.Method == HTTP_METHOD_OPTIONS {
		availableMethods := make([]string, 0, len(r.optionsAllowMethods))
		for _, method := range r.optionsAllowMethods {
			if root, exist := r.roots[method]; exist {
				handler, _ := root.match(req.URL.Path)
				if handler != nil {
					availableMethods = append(availableMethods, method)
				}
			}
		}
		if len(availableMethods) > 0 {
			w.Header().Add("Allow", strings.Join(availableMethods, ", "))
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) Get(path string, handler Handler) {
	if _, exists := r.roots[HTTP_METHOD_GET]; !exists {
		r.roots[HTTP_METHOD_GET] = newRouteTree()
	}
	r.roots[HTTP_METHOD_GET].addRoute(path, handler)
}

func (r *Router) Head(path string, handler Handler) {
	if _, exists := r.roots[HTTP_METHOD_HEAD]; !exists {
		r.roots[HTTP_METHOD_HEAD] = newRouteTree()
	}
	r.roots[HTTP_METHOD_HEAD].addRoute(path, handler)
}

func (r *Router) Post(path string, handler Handler) {
	if _, exists := r.roots[HTTP_METHOD_POST]; !exists {
		r.roots[HTTP_METHOD_POST] = newRouteTree()
	}
	r.roots[HTTP_METHOD_POST].addRoute(path, handler)
}

func (r *Router) Put(path string, handler Handler) {
	if _, exists := r.roots[HTTP_METHOD_PUT]; !exists {
		r.roots[HTTP_METHOD_PUT] = newRouteTree()
	}
	r.roots[HTTP_METHOD_PUT].addRoute(path, handler)
}

func (r *Router) Delete(path string, handler Handler) {
	if _, exists := r.roots[HTTP_METHOD_DELETE]; !exists {
		r.roots[HTTP_METHOD_DELETE] = newRouteTree()
	}
	r.roots[HTTP_METHOD_DELETE].addRoute(path, handler)
}

func (r *Router) Group(path string, fn func(r *GroupRouter)) {
	gr := NewGroupRouter(r, path)
	fn(gr)
}

type GroupRouter struct {
	router *Router
	path   string
}

func NewGroupRouter(r *Router, path string) *GroupRouter {
	gr := &GroupRouter{}
	gr.router = r
	gr.path = strings.TrimRight(path, "/")
	return gr
}

func (gr *GroupRouter) buildPath(path string) string {
	return gr.path + "/" + strings.TrimLeft(path, "/")
}

func (gr *GroupRouter) Get(path string, handler Handler) {
	gr.router.Get(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Head(path string, handler Handler) {
	gr.router.Head(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Post(path string, handler Handler) {
	gr.router.Post(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Put(path string, handler Handler) {
	gr.router.Put(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Delete(path string, handler Handler) {
	gr.router.Delete(gr.buildPath(path), handler)
}
