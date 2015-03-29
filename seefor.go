package r2router

import (
	//"fmt"
	"net/http"
	"time"
)

// Middleware defines how a middle should look like.
// Middlewares are executed in the order they were inserted.
// A middleware can choose to response to a request and not call next
// for continue with the next middleware/handler
type Middleware func(w http.ResponseWriter, req *http.Request, params Params, next func())

// Seefor is a subtype of Router.
// It supports a simple middleware layers.
// Middlewares are always executed before handler,
// no matter where or when they are added.
// And middlewares are executed in the order they were inserted.
type Seefor struct {
	Router
	middlewares []Middleware
	timer       *Timer
}

// NewSeeforRouter for creating a new instance of Seefor router
func NewSeeforRouter() *Seefor {
	c4 := &Seefor{}
	c4.middlewares = make([]Middleware, 0)
	c4.roots = make(map[string]*rootNode)
	c4.HandleMethodNotAllowed = true
	return c4
}

// Implementing http handler interface.
// This is a override of Router.ServeHTTP for handling middlewares
func (c4 *Seefor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if root, exist := c4.roots[req.Method]; exist {
		handler, params, route := root.match(req.URL.Path)
		if handler != nil {
			if c4.timer != nil {
				c4.timeit(route, handler, w, req, params)
			} else {
				c4.handleMiddlewares(handler, w, req, params)
			}
			return
		}
	}
	c4.Router.handleMissing(w, req)
}

func (c4 *Seefor) timeit(route string, handler HandlerFunc, w http.ResponseWriter, req *http.Request, params Params) {
	start := time.Now()
	c4.handleMiddlewares(handler, w, req, params)
	c4.timer.Get(route).Accumulate(start, time.Now())
}

func (c4 *Seefor) handleMiddlewares(handler HandlerFunc, w http.ResponseWriter, req *http.Request, params Params) {
	var next func()

	max := len(c4.middlewares)
	if max == 0 {
		handler(w, req, params)
		return
	}

	counter := 0
	next = func() {
		if counter >= max {
			handler(w, req, params)
			return
		}
		middleware := c4.middlewares[counter]
		counter += 1
		middleware(w, req, params, next)
	}
	next()
}

// Use is for adding middleware
func (c4 *Seefor) Use(middleware ...Middleware) {
	c4.middlewares = append(c4.middlewares, middleware...)
}

// UseHandler is for adding a HandlerFunc as middleware
// it will first wrap handler to middleware and add it to the stack
func (c4 *Seefor) UseHandler(handler HandlerFunc) {
	c4.Use(c4.Wrap(handler))
}

// Wrap for wrapping a handler to middleware
func (c4 *Seefor) Wrap(handler HandlerFunc) Middleware {
	return func(w http.ResponseWriter, req *http.Request, params Params, next func()) {
		handler(w, req, params)
		next()
	}
}

// UseTimer set timer for meaturing endpoint performance.
// If timer is nil then a new timer will be created.
// You can serve statistics internal using Timer as handler
func (c4 *Seefor) UseTimer(timer *Timer) *Timer {
	if timer == nil {
		timer = NewTimer()
	}
	c4.timer = timer
	
	return c4.timer
}