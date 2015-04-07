package r2router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamsAppData(t *testing.T) {
	type User struct {
		Name string
		Age int
	}
	p := params_{}
	p.appData = make(map[interface{}]interface{})
	p.appData["hello"] = "World"
	p.AppSet("user", &User{"CPO", 56})
	
	assert.True(t, p.AppHas("hello"))
	assert.False(t, p.AppHas("world"))
	
	assert.Equal(t, p.AppGet("hello").(string), "World")
	
	user := p.AppGet("user").(*User)
	assert.Equal(t, user.Name, "CPO")
	assert.Equal(t, user.Age, 56)
}

func TestParamsRequestData(t *testing.T) {
	p := params_{}
	p.requestParams = make(map[string]string)
	p.requestParams["hello"] = "World"
	
	assert.True(t, p.Has("hello"))
	assert.False(t, p.Has("world"))
	
	assert.Equal(t, p.Get("hello"), "World")
}

