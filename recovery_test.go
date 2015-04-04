package r2router

import (
	"github.com/stretchr/testify/assert"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSeeforRecovery(t *testing.T) {
	router := NewSeeforRouter()
	router.Before(NewRecovery(NewRecoveryOptions()))

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		panic("This shouldn't crash Seefor")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(content), "Internal Server Error")
}

func TestSeeforRecoveryMiddlewarePanic(t *testing.T) {
	router := NewSeeforRouter()
	options := NewRecoveryOptions()
	options.PrintStack = true
	rec := NewRecovery(options)
	router.Before(rec)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		panic("This shouldn't crash Seefor")
	})

	router.Before(WrapBeforeHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Middleware panic")
	})))

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Contains(t, string(content), "Middleware panic")
	assert.NotContains(t, string(content), "This shouldn't crash Seefor")
}

func TestSeeforRecoveryPrintStack(t *testing.T) {
	router := NewSeeforRouter()
	options := NewRecoveryOptions()
	options.PrintStack = true
	rec := NewRecovery(options)
	router.Before(rec)

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		panic("This shouldn't crash Seefor")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Contains(t, string(content), "This shouldn't crash Seefor")
}
