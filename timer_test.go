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
		before := time.Now()
		w.Add(1)
		go func(before, after time.Time) {
			end := time.Now()
			c.Accumulate(before, after, end, time.Now())
			w.Done()
		}(before, time.Now())
	}
	w.Wait()
	assert.NotNil(t, c.Count)
	assert.Equal(t, c.Count, int64(1000))
}

func TestTimer(t *testing.T) {
	timer := NewTimer()
	w := sync.WaitGroup{}
	
	for i := 0; i < 1000; i++ {
		before := time.Now()
		w.Add(1)
		go func(before, after time.Time, name string) {
			end := time.Now()
			timer.Get(name).Accumulate(before, after, end, time.Now())
			w.Done()
		}(before, time.Now(), string(i%25))
	}
	w.Wait()
	assert.NotNil(t, timer.routes)
	assert.Equal(t, len(timer.routes), 25)
}

func TestTimerStats(t *testing.T) {
	timer := NewTimer()
	w := sync.WaitGroup{}
	
	for i := 0; i < 1000; i++ {
		before := time.Now()
		w.Add(1)
		go func(before, after time.Time, name string) {
			end := time.Now()
			timer.Get(name).Accumulate(before, after, end, time.Now())
			w.Done()
		}(before, time.Now(), "r"+string(i%25))
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
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"\"")
	
	// TODO, assert ordering
	res, err = http.Get(ts.URL+"/?sort=max")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"max\"")
	
	res, err = http.Get(ts.URL+"/?sort=tot")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"tot\"")
	
	res, err = http.Get(ts.URL+"/?sort=count")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"count\"")
	
	
	res, err = http.Get(ts.URL+"/?sort=avg_before")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"avg_before\"")
	
	
	res, err = http.Get(ts.URL+"/?sort=avg_after")
	assert.Nil(t, err)
	content, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	contentString = string(content)
	assert.Contains(t, contentString, "\"result\":[{\"route\":\"")
	assert.Contains(t, contentString, "\"count\":40")
	assert.Contains(t, contentString, "\"sortBy\":\"avg_after\"")
}