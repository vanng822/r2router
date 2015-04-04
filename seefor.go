package r2router

import (
	//"fmt"
	"net/http"
	"time"
)

// Shortcut for map[string]interface{}
// Helpful to build data for json response
type M map[string]interface{}

// Before defines middleware interface.
// Before middlewares are for handling request before routing.
// Before middlewares are executed in the order they were inserted.
// A middleware can choose to response to a request and not call the handler
type Before func(handler http.Handler) http.Handler

// After defines how a middleware should look like.
// After middlewares are for handling request after routing.
// After middlewares are executed in the order they were inserted.
// A middleware can choose to response to a request and not call next
// for continue with the next middleware/handler
type After interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request, params Params, next func())
}

// AfterFunc for wrapping a function that have this signature to use as
// After middleware
type AfterFunc func(w http.ResponseWriter, req *http.Request, params Params, next func())

func (a AfterFunc) ServeHTTP(w http.ResponseWriter, req *http.Request, params Params, next func()) {
	a(w, req, params, next)
}

// Seefor is a subtype of Router.
// It supports a simple middleware layers.
// Middlewares are always executed before handler,
// no matter where or when they are added.
// And middlewares are executed in the order they were inserted.
type Seefor struct {
	Router
	befores []Before
	afters  []After
	timer   *Timer
}

// NewSeeforRouter for creating a new instance of Seefor router
func NewSeeforRouter() *Seefor {
	c4 := &Seefor{}
	c4.afters = make([]After, 0)
	c4.befores = make([]Before, 0)
	c4.roots = make(map[string]*rootNode)
	c4.HandleMethodNotAllowed = true
	return c4
}

// Implementing http handler interface.
// This is a override of Router.ServeHTTP for handling middlewares
func (c4 *Seefor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	started := time.Now()
	c4.handleBeforeMiddlewares(w, req, func() {
		beforeEnd := time.Now()
		if root, exist := c4.roots[req.Method]; exist {
			handler, params, route := root.match(req.URL.Path)
			if handler != nil {
				if c4.timer != nil {
					c4.timeit(route, started, beforeEnd, handler, w, req, params)
				} else {
					c4.handleAfterMiddlewares(handler, w, req, params)
				}
				return
			}
		}
		c4.Router.handleMissing(w, req)
	})
}

func (c4 *Seefor) handleBeforeMiddlewares(w http.ResponseWriter, req *http.Request, nextHandler func()) {
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		nextHandler()
	})

	for i := len(c4.befores) - 1; i >= 0; i-- {
		handler = c4.befores[i](handler)
	}

	handler.ServeHTTP(w, req)
}

func (c4 *Seefor) timeit(route string, started, beforeEnd time.Time, handler Handler, w http.ResponseWriter, req *http.Request, params Params) {
	after := time.Now()
	c4.handleAfterMiddlewares(handler, w, req, params)
	c4.timer.Get(route).Accumulate(started, beforeEnd, after, time.Now())
}

func (c4 *Seefor) handleAfterMiddlewares(handler Handler, w http.ResponseWriter, req *http.Request, params Params) {
	var next func()

	max := len(c4.afters)
	if max == 0 {
		handler.ServeHTTP(w, req, params)
		return
	}

	counter := 0
	next = func() {
		if counter >= max {
			handler.ServeHTTP(w, req, params)
			return
		}
		middleware := c4.afters[counter]
		counter += 1
		middleware.ServeHTTP(w, req, params, next)
	}
	next()
}

// Before is for adding middleware for running before routing
func (c4 *Seefor) Before(middleware ...Before) {
	c4.befores = append(c4.befores, middleware...)
}

// After is for adding middleware for running after
func (c4 *Seefor) After(middleware ...After) {
	c4.afters = append(c4.afters, middleware...)
}

// Wrap for wrapping a handler to After middleware
// Be aware that it will not be able to stop execution propagation
// That is it will continue to execute the next middleware/handler
func Wrap(handler HandlerFunc) After {
	return AfterFunc(func(w http.ResponseWriter, req *http.Request, params Params, next func()) {
		handler(w, req, params)
		next()
	})
}

// WrapHandler for wrapping a http.Handler to After middleware.
// Be aware that it will not be able to stop execution propagation
// That is it will continue to execute the next middleware/handler
func WrapHandler(handler http.Handler) After {
	return AfterFunc(func(w http.ResponseWriter, req *http.Request, _ Params, next func()) {
		handler.ServeHTTP(w, req)
		next()
	})
}

// WrapBeforeHandler for wrapping a http.Handler to Before middleware.
// Be aware that it will not be able to stop execution propagation
// That is it will continue to execute the next middleware/handler
func WrapBeforeHandler(handler http.Handler) Before {
	return func(their http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			their.ServeHTTP(w, req)
		})
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
