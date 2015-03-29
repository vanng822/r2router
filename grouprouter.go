package r2router

import (
	"strings"
)


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

func (gr *GroupRouter) Get(path string, handler HandlerFunc) {
	gr.router.Get(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Head(path string, handler HandlerFunc) {
	gr.router.Head(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Post(path string, handler HandlerFunc) {
	gr.router.Post(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Put(path string, handler HandlerFunc) {
	gr.router.Put(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Delete(path string, handler HandlerFunc) {
	gr.router.Delete(gr.buildPath(path), handler)
}

func (gr *GroupRouter) Patch(path string, handler HandlerFunc) {
	gr.router.Patch(gr.buildPath(path), handler)
}