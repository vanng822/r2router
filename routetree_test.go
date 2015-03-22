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
