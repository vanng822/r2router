package r2router

import (
	"fmt"
	"strings"
	"net/url"
)

// For managing route and getting url
type RouteManager interface {
	// Register a route and return the path
	// This can be good for adding and register handler at the same time
	// router.Get(rm.Add("user", "/user/:id"), handler)
	Add(routeName, path string) string
	// Return the path for a specific route name
	// Use for register handler
	// router.Delete(rm.PathFor("user"), handler)
	PathFor(routeName string) string
	// Returning url for given route name and provided data
	// Will panic if missmatched
	UrlFor(routeName string, params map[string]interface{}) string
}

type routeManager struct {
	routes map[string]string
}

func NewRouteManager() RouteManager {
	m := &routeManager{}
	m.routes = make(map[string]string)
	
	return m
}

func (m *routeManager) Add(routeName, path string) string {
	if _, exist := m.routes[routeName]; exist {
		panic("Route name must be unique")
	}

	m.routes[routeName] = path
	return path
}

func (m *routeManager) PathFor(routeName string) string {
	if path, exist := m.routes[routeName]; exist {
		return path
	}
	panic(fmt.Sprintf("Could not find any path for route name: %s", routeName))
}

func (m *routeManager) UrlFor(routeName string, params map[string]interface{}) string {
	path := m.PathFor(routeName)
	paths := strings.Split(path, "/")
	parts := make([]string, 0)
	data := make(map[string]string)
	for key, val := range params {
		// could use type switch here
		data[key] = fmt.Sprintf("%v", val)
	}
	for _, p := range paths {
		if p == "" || p[:1] != ":" {
			parts = append(parts, p)
			continue
		}
		key := p[1:]
		if val, exist := data[key]; exist {
			parts = append(parts, val)
			delete(data, key)
			continue
		}

		panic(fmt.Sprintf("Param %s missing in provided data", key))
	}
	
	if len(data) > 0 {
		urlParams := url.Values{}
		for key, val := range data {
			urlParams.Add(key, val)
		}
		return fmt.Sprintf("%s?%s", strings.Join(parts, "/"), urlParams.Encode())
	}
	return strings.Join(parts, "/")
}
