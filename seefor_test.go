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
	router := NewSeeforRouter()
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

func TestSeeforMiddleware(t *testing.T) {
	router := NewSeeforRouter()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id") + p.AppGet("middleware").(string)))
	})
	router.After(Wrap(func(w http.ResponseWriter, r *http.Request, p Params) {
		p.AppSet("middleware", "Test Middleware")
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(content), "GET:/user/keys/:id,testingTest Middleware")
}

func TestSeeforMultiMiddleware(t *testing.T) {
	router := NewSeeforRouter()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id") + p.AppGet("middleware").(string) + p.AppGet("hello").(string)))
	})
	router.After(Wrap(func(w http.ResponseWriter, r *http.Request, p Params) {
		p.AppSet("middleware", "Test Middleware")
		p.AppSet("hello", "World")
	}))

	router.After(Wrap(func(w http.ResponseWriter, r *http.Request, p Params) {
		p.AppSet("middleware", "Middleware")
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(content), "GET:/user/keys/:id,testingMiddlewareWorld")
}

func TestSeeforMiddlewareStop(t *testing.T) {
	router := NewSeeforRouter()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id") + p.AppGet("middleware").(string) + p.AppGet("hello").(string)))
	})
	router.After(func(next Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request, p Params) {
			w.Write([]byte("Hello"))
		})
	})

	router.After(Wrap(func(w http.ResponseWriter, r *http.Request, p Params) {
		p.AppSet("middleware", "Middleware")
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(content), "Hello")
}

func TestSeeforMiddlewareContinue(t *testing.T) {
	router := NewSeeforRouter()

	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	router.After(WrapHandler(handler))

	router.After(Wrap(func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("World"))
	}))
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(content), "HelloWorldGET:/user/keys/:id,testing")
}

func TestSeeforTimer(t *testing.T) {
	router := NewSeeforRouter()

	router.Get("/hello", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("world"))
	})

	timer := router.UseTimer(nil)
	timer2 := router.UseTimer(nil)
	assert.Exactly(t, timer, timer2)
	assert.True(t, assert.ObjectsAreEqual(timer, timer2))

	ts := httptest.NewServer(router)
	defer ts.Close()

	timers := httptest.NewServer(timer)
	defer timers.Close()

	res, err := http.Get(ts.URL + "/hello")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, content, []byte("world"))

	res, err = http.Get(timers.URL)
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, string(content), "\"result\":[{\"route\":\"/hello\",\"count\":1,\"tot\":")
}

func TestMiddlewareBefore(t *testing.T) {
	router := NewSeeforRouter()

	// never call
	router.Get("/user/keys/:id", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("GET:/user/keys/:id," + p.Get("id")))
	})

	// wrapping which always call next handler
	router.Before(WrapBeforeHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})))

	// customized which choose not to call next handler
	router.Before(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(" World"))
		})
	})
	// never call
	router.Before(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ignore"))
		})
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/keys/testing")
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(content), "Hello World")
}
