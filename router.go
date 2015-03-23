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
)

// Holding value for named parameters
type Params map[string]string

// Handler define interface handler
type Handler func(http.ResponseWriter, *http.Request, Params)

type Router struct {
	roots               map[string]*rootNode
	
	HandleMethodNotAllowed bool
	
	MethodNotAllowedHandler http.HandlerFunc
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
		handler, params := root.match(req.URL.Path)
		if handler != nil {
			handler(w, req, params)
			//log.Println(time.Now().Sub(now))
			return
		}
	}
	
	// if options find handler for different method
	if req.Method == HTTP_METHOD_OPTIONS {
		// build and serve options
		availableMethods := make([]string, 0, len(r.roots))
		for method := range r.roots {
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
	
	if r.HandleMethodNotAllowed {
		for method := range r.roots {
			if root, exist := r.roots[method]; exist {
				handler, _ := root.match(req.URL.Path)
				if handler != nil {
					if r.MethodNotAllowedHandler != nil {
						r.MethodNotAllowedHandler(w, req)
					} else {
						http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed) 
					}
					return
				}
			}
		}	
	}

	http.NotFound(w, req)
}

func (r *Router) Get(path string, handler Handler) {
	r.AddHandler(HTTP_METHOD_GET, path, handler)
}

func (r *Router) Head(path string, handler Handler) {
	r.AddHandler(HTTP_METHOD_HEAD, path, handler)
}

func (r *Router) Post(path string, handler Handler) {
	r.AddHandler(HTTP_METHOD_POST, path, handler)
}

func (r *Router) Put(path string, handler Handler) {
	r.AddHandler(HTTP_METHOD_PUT, path, handler)
}

func (r *Router) Delete(path string, handler Handler) {
	r.AddHandler(HTTP_METHOD_DELETE, path, handler)
}

func (r *Router) AddHandler(method, path string, handler Handler) {
	if _, exists := r.roots[method]; !exists {
		r.roots[method] = newRouteTree()
	}
	r.roots[method].addRoute(path, handler)
}

// Group takes a path which typically a prefix for an endpoint
// It will call callback function with a group router which
// you can add handler for different request methods 
func (r *Router) Group(path string, fn func(r *GroupRouter)) {
	gr := NewGroupRouter(r, path)
	fn(gr)
}

// GroupRouter is a helper for grouping endpoints
// All methods are proxy of Router
// Suitable for grouping different methods of an endpoint
type GroupRouter struct {
	router *Router
	path   string
}
// NewGroupRouter return GroupRouter which is a helper
// to construct a group of endpoints, such example could
// be API-version or different methods for an endpoint
// You should always use router.Group instead of using this directly
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
