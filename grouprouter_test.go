package r2router

import (
	"github.com/stretchr/testify/assert"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroup(t *testing.T) {
	router := NewRouter()
	router.Group("/user/keys", func(r *GroupRouter) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("GET:/user/keys"))
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("POST:/user/keys"))
		})
		r.Put("/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("PUT:/user/keys/:id," + p.Get("id")))
		})
		r.Delete("/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("DELETE:/user/keys/:id," + p.Get("id")))
		})
		r.Get("/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
		})
		r.Head("/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("HEAD:/user/keys/:id," + p.Get("id")))
		})
		
		r.Patch("/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("PATCH:/user/keys/:id," + p.Get("id")))
		})
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	// get
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(content), "GET:/user/keys/:id,testing")

	// Get all
	res, err = http.Get(ts.URL + "/user/keys")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, string(content), "GET:/user/keys")

	// post
	res, err = http.Post(ts.URL+"/user/keys", "", nil)
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, string(content), "POST:/user/keys")

	// put
	client := &http.Client{}
	req, err := http.NewRequest("PUT", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, string(content), "PUT:/user/keys/:id,testing")

	// delete
	req, err = http.NewRequest("DELETE", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, string(content), "DELETE:/user/keys/:id,testing")

	req, err = http.NewRequest("HEAD", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	// http seems not sending content
	assert.Equal(t, int(res.ContentLength), len([]byte("HEAD:/user/keys/:id,testing")))
	
	req, err = http.NewRequest("PATCH", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	// http seems not sending content
	assert.Equal(t, int(res.ContentLength), len([]byte("PATCH:/user/keys/:id,testing")))
}
