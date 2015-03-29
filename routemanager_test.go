package r2router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlFor(t *testing.T) {
	m := NewRouteManager()
	m.Add("some::for", "/some/:key/for")

	assert.Equal(t, m.UrlFor("some::for", map[string]interface{}{"key": 100}), "/some/100/for")

	assert.Equal(t, m.UrlFor("some::for", map[string]interface{}{"key": 10.5}), "/some/10.5/for")

	assert.Equal(t, m.UrlFor("some::for", map[string]interface{}{"key": "thing"}), "/some/thing/for")

	assert.Equal(t, m.UrlFor("some::for", map[string]interface{}{"key": "thing", "your": "kitty"}), "/some/thing/for?your=kitty")

	assert.Panics(t, func() {
		m.UrlFor("some::for", map[string]interface{}{})
	})
}

func TestAddNoneUnique(t *testing.T) {
	m := NewRouteManager()
	m.Add("same::for", "/some/:key/for")
	assert.Panics(t, func() {
		m.Add("same::for", "/some/:key/for")
	})
}

func TestPathForMissing(t *testing.T) {
	m := NewRouteManager()
	m.Add("same::for", "/some/:key/for")
	assert.Panics(t, func() {
		m.PathFor("same::too")
	})
}
