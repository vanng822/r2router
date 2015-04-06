package r2router

import (
	"github.com/stretchr/testify/assert"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()
	router.Get("/user/keys/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys"))
	})
	router.Post("/user/keys/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("POST:/user/keys"))
	})
	router.Put("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("PUT:/user/keys/:id," + p.Get("id")))
	})
	router.Delete("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("DELETE:/user/keys/:id," + p.Get("id")))
	})
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
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

	// options
	req, err = http.NewRequest("OPTIONS", ts.URL+"/user/keys/testing", nil)
	res, err = client.Do(req)
	res.Body.Close()
	assert.Contains(t, res.Header.Get("Allow"), "GET")
	assert.Contains(t, res.Header.Get("Allow"), "PUT")
	assert.Contains(t, res.Header.Get("Allow"), "DELETE")
}

func TestRouterMethodNotAllowed(t *testing.T) {
	router := NewRouter()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", ts.URL+"/user/keys/testing", nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusMethodNotAllowed)
}

func TestRouterMethodNotAllowedCustomized(t *testing.T) {
	router := NewRouter()
	router.MethodNotAllowed = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Hello"))
	}
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", ts.URL+"/user/keys/testing", nil)
	assert.Nil(t, err)

	res, err := client.Do(req)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, string(content), "Hello")
	assert.Equal(t, res.StatusCode, http.StatusMethodNotAllowed)

}

func TestRouterCustomizedOptions(t *testing.T) {
	router := NewRouter()

	router.AddHandler("OPTIONS", "/users/", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Header().Add("Allow", "GET, POST")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()
	client := &http.Client{}
	// custom
	req, _ := http.NewRequest("OPTIONS", ts.URL+"/users/", nil)
	res2, _ := client.Do(req)
	res2.Body.Close()
	assert.Equal(t, res2.Header.Get("Allow"), "GET, POST")
}

func TestRouterNotFound(t *testing.T) {
	router := NewRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}

func TestRouterNotFoundCustomized(t *testing.T) {
	router := NewRouter()
	router.NotFound = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Boo"))
	}
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
	assert.Equal(t, content, []byte("Boo"))
}

func TestRouterFirstNodeParamNode(t *testing.T) {
	router := NewRouter()

	router.Get("/:page", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/:page," + p.Get("page")))
	})

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", ts.URL+"/testing", nil)
	assert.Nil(t, err)
	res, err := client.Do(req)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, string(content), "GET:/:page,testing")

	req, err = http.NewRequest("GET", ts.URL+"/user/keys/testing", nil)
	assert.Nil(t, err)
	res, err = client.Do(req)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, string(content), "GET:/user/keys/:id,testing")
}

func TestRouterDump(t *testing.T) {
	router := NewRouter()
	router.Get("/:page", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/:page," + p.Get("page")))
	})
	assert.Contains(t, router.Dump(), " |\n  -- \n  |\n   -- :page (<")
}
