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