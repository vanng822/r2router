package r2router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	Route string        `json:"route"`
	Count int64         `json:"count"`
	Tot   time.Duration `json:"tot"`
	Max   time.Duration `json:"max"`
	Min   time.Duration `json:"min"`
	Avg   time.Duration `json:"avg"`
}

type Stats struct {
	Generated time.Time `json:"generated"`
	UpTime    string    `json:"upTime"`
	Result    []*Stat   `json:"result"`
	SortBy    string    `json:"sortBy"`
}

// Implements sort interface
func (s *Stats) Len() int {
	return len(s.Result)
}

func (s *Stats) Swap(i, j int) {
	s.Result[j], s.Result[i] = s.Result[i], s.Result[j]
}

func (s *Stats) Less(i, j int) bool {
	switch s.SortBy {
	case "count":
		return s.Result[i].Count < s.Result[j].Count
	case "tot":
		return s.Result[i].Tot < s.Result[j].Tot
	case "max":
		return s.Result[i].Max < s.Result[j].Max
	default:
		return s.Result[i].Avg < s.Result[j].Avg
	}
}

// For serving statistics
func (t *Timer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	sortBy := req.Form.Get("sort")

	stats := &Stats{}
	stats.SortBy = strings.ToLower(sortBy)
	stats.Generated = time.Now()
	stats.UpTime = fmt.Sprintf("%s", stats.Generated.Sub(t.Since))
	// Read access OK for map?
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
	sort.Sort(sort.Reverse(stats))
	jsonData, _ := json.Marshal(stats)
	w.Write(jsonData)
}
