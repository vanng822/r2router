package r2router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRouteDuplicate(t *testing.T) {
	r := newRouteTree()
	assert.Panics(t, func() {
		r.addRoute("/users/:user/events", httpTestHandler)
		r.addRoute("/users/:user/events", httpTestHandler)
	})
}


func TestAddRouteDiffParamName(t *testing.T) {
	r := newRouteTree()
	assert.Panics(t, func() {
		r.addRoute("/users/:username/following/:target_user", httpTestHandler)
		r.addRoute("/users/:user/events", httpTestHandler)
	})
}

func TestSwapNodes(t *testing.T) {
	r := newRouteTree()
	r.addRoute("/users/:user/events", httpTestHandler)
	r.addRoute("/users/list/events", httpTestHandler)
	assert.Equal(t, len(r.root.children[0].children), 2)
	assert.Equal(t, r.root.children[0].children[0].path, "list")
	assert.True(t, r.root.children[0].children[1].paramNode)
	assert.Equal(t, r.root.children[0].children[1].paramName, "user")
}

func TestAddRouteBrokenParamName(t *testing.T) {
	r := newRouteTree()
	assert.Panics(t, func() {
		r.addRoute("/users/:/following", httpTestHandler)
	})
}

func TestAddRouteIndex(t *testing.T) {
	r := newRouteTree()
	r.addRoute("/", httpTestHandler)
	assert.NotNil(t, r.handler)
	assert.Panics(t, func() {
		r.addRoute("/", httpTestHandler)
	})
}

func TestMatchTrue(t *testing.T) {
	r := newRouteTree()
	r.addRoute("/users/:user/events", httpTestHandler)
	h, p, route := r.match("/users/vanng822/events")
	assert.NotNil(t, h)
	exectedP := &params_{}
	exectedP.appData = make(map[string]interface{})
	exectedP.requestParams = make(map[string]string)
	exectedP.requestParams["user"] = "vanng822"
	assert.Equal(t, p, exectedP)
	
	assert.Equal(t, route, "/users/:user/events")
}

func TestMatchIndex(t *testing.T) {
	r := newRouteTree()
	r.addRoute("/", httpTestHandler)
	h, p, route := r.match("/")
	assert.NotNil(t, h)
	assert.Nil(t, p)
	assert.Equal(t, route, "/")
}

func TestMatchFalse(t *testing.T) {
	r := newRouteTree()
	r.addRoute("/users/:user/events", httpTestHandler)
	h, p, route := r.match("/users/:user/orgs")
	assert.Nil(t, h)
	assert.Nil(t, p)
	assert.Equal(t, route, "")
}

