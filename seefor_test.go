package r2router

import (
	"github.com/stretchr/testify/assert"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSeeforInherits(t *testing.T) {
	router := NewSeefor()
	router.Get("/user/keys/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys"))
	})
	router.Post("/user/keys/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("POST:/user/keys"))
	})
	router.Put("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("PUT:/user/keys/:id," + p["id"]))
	})
	router.Delete("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("DELETE:/user/keys/:id," + p["id"]))
	})
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p["id"]))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, content, []byte("GET:/user/keys/:id,testing"))

	// Get all
	res, err = http.Get(ts.URL + "/user/keys")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, content, []byte("GET:/user/keys"))

	// post
	res, err = http.Post(ts.URL+"/user/keys", "", nil)
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, content, []byte("POST:/user/keys"))

	// put
	client := &http.Client{}
	req, err := http.NewRequest("PUT", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, content, []byte("PUT:/user/keys/:id,testing"))

	// delete
	req, err = http.NewRequest("DELETE", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, content, []byte("DELETE:/user/keys/:id,testing"))

	// options
	req, err = http.NewRequest("OPTIONS", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	res.Body.Close()
	assert.Contains(t, res.Header.Get("Allow"), "GET")
	assert.Contains(t, res.Header.Get("Allow"), "PUT")
	assert.Contains(t, res.Header.Get("Allow"), "DELETE")
}

func TestSeeforMiddleware(t *testing.T) {
	router := NewSeefor()
	
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p["id"] +  p["middleware"]))
	})
	router.Use(func(w http.ResponseWriter, r *http.Request, p Params, next func()) {
		p["middleware"] = "Test Middleware"
		next()
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, content, []byte("GET:/user/keys/:id,testingTest Middleware"))
}


func TestSeeforMultiMiddleware(t *testing.T) {
	router := NewSeefor()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p["id"] + p["middleware"] + p["hello"]))
	})
	router.Use(func(w http.ResponseWriter, r *http.Request, p Params, next func()) {
		p["middleware"] = "Test Middleware"
		p["hello"] = "World"
		next()
	})

	router.UseHandler(func(w http.ResponseWriter, r *http.Request, p Params) {
		p["middleware"] = "Middleware"
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, content, []byte("GET:/user/keys/:id,testingMiddlewareWorld"))
}

func TestSeeforMiddlewareWritten(t *testing.T) {
	router := NewSeefor()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p["id"] + p["middleware"] + p["hello"]))
	})
	router.Use(func(w http.ResponseWriter, r *http.Request, p Params, next func()) {
		w.WriteHeader(http.StatusNotFound)
	})


	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}
