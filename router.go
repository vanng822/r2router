package r2router

import (
	//"log"
	"net/http"
	"strings"
	//"time"
)

const (
	HTTP_METHOD_GET     = "GET"
	HTTP_METHOD_POST    = "POST"
	HTTP_METHOD_DELETE  = "DELETE"
	HTTP_METHOD_OPTIONS = "OPTIONS"
	HTTP_METHOD_HEAD    = "HEAD"
	HTTP_METHOD_PUT     = "PUT"
	HTTP_METHOD_PATCH   = "PATCH"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request, params Params)
}

// HandlerFunc define interface handler
type HandlerFunc func(w http.ResponseWriter, req *http.Request, params Params)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request, params Params) {
	h(w, req, params)
}

type Router struct {
	roots                  map[string]*rootNode
	HandleMethodNotAllowed bool
	MethodNotAllowed       http.HandlerFunc
	NotFound               http.HandlerFunc
}

// NewRouter return a new Router
func NewRouter() *Router {
	r := &Router{}
	r.roots = make(map[string]*rootNode)
	r.HandleMethodNotAllowed = true
	return r
}

// http Handler Interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//now := time.Now()
	if root, exist := r.roots[req.Method]; exist {
		handler, params, _ := root.match(req.URL.Path)
		if handler != nil {
			handler.ServeHTTP(w, req, params)
			//log.Println(time.Now().Sub(now))
			return
		}
	}
	r.handleMissing(w, req)
}

func (r *Router) handleMissing(w http.ResponseWriter, req *http.Request) {
	// if options find handler for different method
	if req.Method == HTTP_METHOD_OPTIONS {
		// build and serve options
		availableMethods := make([]string, 0, len(r.roots))
		for method := range r.roots {
			if root, exist := r.roots[method]; exist {
				handler, _, _ := root.match(req.URL.Path)
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

	if r.HandleMethodNotAllowed {
		for method := range r.roots {
			if root, exist := r.roots[method]; exist {
				handler, _, _ := root.match(req.URL.Path)
				if handler != nil {
					if r.MethodNotAllowed != nil {
						r.MethodNotAllowed(w, req)
					} else {
						http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
					}
					return
				}
			}
		}
	}

	if r.NotFound != nil {
		r.NotFound(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_GET, path, handler)
}

func (r *Router) Head(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_HEAD, path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_POST, path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_PUT, path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_DELETE, path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.AddHandler(HTTP_METHOD_PATCH, path, handler)
}

func (r *Router) AddHandler(method, path string, handler HandlerFunc) {
	if _, exists := r.roots[method]; !exists {
		r.roots[method] = newRouteTree()
	}
	r.roots[method].addRoute(path, HandlerFunc(handler))
}

// Group takes a path which typically a prefix for an endpoint
// It will call callback function with a group router which
// you can add handler for different request methods
func (r *Router) Group(path string, fn func(r *GroupRouter)) {
	gr := NewGroupRouter(r, path)
	fn(gr)
}


func (r *Router) Dump() string {
	s := ""
	for method, root := range r.roots {
		s += method + "\n" +root.dump()
	}
	
	return s
}