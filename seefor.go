package r2router

import (
	//"fmt"
	"net/http"
)

type Middleware func(w http.ResponseWriter, req *http.Request, params Params, next func())

type Seefor struct {
	Router
	middlewares []Middleware
}

func NewSeeforRouter() *Seefor {
	c4 := &Seefor{}
	c4.middlewares = make([]Middleware, 0)
	c4.roots = make(map[string]*rootNode)
	c4.HandleMethodNotAllowed = true
	return c4
}

func (c4 *Seefor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if root, exist := c4.roots[req.Method]; exist {
		handler, params := root.match(req.URL.Path)
		if handler != nil {
			c4.handleMiddlewares(handler, w, req, params)
			return
		}
	}
	c4.Router.handleMissing(w, req)
}

func (c4 *Seefor) handleMiddlewares(handler Handler, w http.ResponseWriter, req *http.Request, params Params) {
	var next func()
	max := len(c4.middlewares) - 1
	counter := -1
	next = func() {
		if counter >= max {
			handler(w, req, params)
			return
		}
		counter += 1
		c4.middlewares[counter](w, req, params, next)
	}
	next()
}

func (c4 *Seefor) Use(middleware ...Middleware) {
	c4.middlewares = append(c4.middlewares, middleware...)
}

func (c4 *Seefor) UseHandler(handler Handler) {
	c4.Use(c4.Wrap(handler))
}

func (c4 *Seefor) Wrap(handler Handler) Middleware {
	return func(w http.ResponseWriter, req *http.Request, params Params, next func()) {
		handler(w, req, params)
		next()
	}
}
