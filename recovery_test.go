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
	router.Before(NewRecovery())
	
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


func TestSeeforRecoveryPrintStack(t *testing.T) {
	router := NewSeeforRouter()
	rec := NewRecovery()
	rec.PrintStack = true
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