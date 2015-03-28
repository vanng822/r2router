package r2router

import (
	"github.com/stretchr/testify/assert"
	//"fmt"
	"testing"
	"time"
	"sync"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func TestCounter(t *testing.T) {
	c := &Counter{}
	w := sync.WaitGroup{}
	
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func(start time.Time) {
			end := time.Now()
			c.Accumulate(start, end)
			w.Done()
		}(time.Now())
	}
	w.Wait()
	assert.NotNil(t, c.Count)
	assert.Equal(t, c.Count, 1000)
}

func TestTimer(t *testing.T) {
	timer := NewTimer()
	w := sync.WaitGroup{}
	
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func(start time.Time, name string) {
			end := time.Now()
			timer.Get(name).Accumulate(start, end)
			w.Done()
		}(time.Now(), string(i%25))
	}
	w.Wait()
	assert.NotNil(t, timer.routes)
	assert.Equal(t, len(timer.routes), 25)
}

func TestTimerStats(t *testing.T) {
	timer := NewTimer()
	w := sync.WaitGroup{}
	
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func(start time.Time, name string) {
			end := time.Now()
			timer.Get(name).Accumulate(start, end)
			w.Done()
		}(time.Now(), "r"+string(i%25))
	}
	w.Wait()
	assert.NotNil(t, timer.routes)
	assert.Equal(t, len(timer.routes), 25)
	
	ts := httptest.NewServer(timer)
	defer ts.Close()
	
	// get
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString := string(content)
	assert.Contains(t, contentString, "\"Result\":[{\"Route\":\"")
	assert.Contains(t, contentString, "\"Count\":40")
	
	
	// TODO, assert ordering
	res, err = http.Get(ts.URL+"/?sort=max")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"Result\":[{\"Route\":\"")
	assert.Contains(t, contentString, "\"Count\":40")
	
	res, err = http.Get(ts.URL+"/?sort=tot")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"Result\":[{\"Route\":\"")
	assert.Contains(t, contentString, "\"Count\":40")
	
	
	res, err = http.Get(ts.URL+"/?sort=count")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"Result\":[{\"Route\":\"")
	assert.Contains(t, contentString, "\"Count\":40")
}