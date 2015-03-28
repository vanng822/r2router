package r2router

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"fmt"
)

type Counter struct {
	Count int64
	Tot   time.Duration
	Max   time.Duration
	Min   time.Duration
	Avg   time.Duration
}

func (c *Counter) Accumulate(start time.Time, end time.Time) {
	d := int64(end.Sub(start))
	tot := atomic.AddInt64((*int64)(&c.Tot), d)
	count := atomic.AddInt64((*int64)(&c.Count), 1)
	atomic.StoreInt64((*int64)(&c.Avg), tot/count)
	max := int64(c.Max)
	if d > max {
		atomic.CompareAndSwapInt64((*int64)(&c.Max), max, d)
	}
	min := int64(c.Min)
	if d < min {
		atomic.CompareAndSwapInt64((*int64)(&c.Min), min, d)
	}

}

type Timer struct {
	Since  time.Time
	routes map[string]*Counter
	mux    sync.Mutex
}

func NewTimer() *Timer {
	t := &Timer{}
	t.Since = time.Now()
	t.routes = make(map[string]*Counter)
	return t
}

func (t *Timer) Get(name string) *Counter {
	if c, exist := t.routes[name]; exist {
		return c
	}
	t.mux.Lock()
	t.routes[name] = &Counter{}
	t.mux.Unlock()
	return t.routes[name]
}

type Stat struct {
	Route string
	Count int64
	Tot   time.Duration
	Max   time.Duration
	Min   time.Duration
	Avg   time.Duration
}

type Stats struct {
	Generated time.Time
	UpTime    string
	Result    []*Stat
}

// For serving statistics
func (t *Timer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Read access OK map?
	stats := &Stats{}
	stats.Generated = time.Now()
	stats.UpTime = fmt.Sprintf("%s", stats.Generated.Sub(t.Since))
	stats.Result = make([]*Stat, 0, len(t.routes))
	for k, v := range t.routes {
		stat := &Stat{}
		stat.Route = k
		stat.Count = v.Count
		stat.Tot = v.Tot
		stat.Avg = v.Avg
		stat.Max = v.Max
		stat.Min = v.Min
		stats.Result = append(stats.Result, stat)
	}
	jsonData, _ := json.Marshal(stats)
	w.Write(jsonData)
}
