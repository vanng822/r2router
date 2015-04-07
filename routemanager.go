package r2router

import (
	"fmt"
	"net/url"
	"strings"
)

// Shortcut for map[string][]string
type P map[string][]string

// For managing route and getting url
type RouteManager interface {
	// For setting baseurl if one needs full url
	SetBaseUrl(baseUrl string)
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
	UrlFor(routeName string, params map[string][]string) string
	// Returning url for given path and provided data
	// Will panic if missmatched
	UrlForPath(path string, params map[string][]string) string
}

type routeManager struct {
	baseUrl string
	routes  map[string]string
}

func NewRouteManager() RouteManager {
	m := &routeManager{}
	m.routes = make(map[string]string)

	return m
}

func (m *routeManager) SetBaseUrl(baseUrl string) {
	m.baseUrl = strings.TrimRight(baseUrl, "/")
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

func (m *routeManager) UrlFor(routeName string, params map[string][]string) string {
	return m.UrlForPath(m.PathFor(routeName), params)
}

func (m *routeManager) UrlForPath(path string, params map[string][]string) string {
	paths := strings.Split(path, "/")
	parts := make([]string, 0)

	urlParams := url.Values{}
	for k, v := range params {
		for _, vv := range v {
			urlParams.Add(k, vv)
		}
	}

	for _, p := range paths {
		if p == "" || p[:1] != ":" {
			parts = append(parts, p)
			continue
		}
		key := p[1:]
		if val, exist := params[key]; exist {
			if len(val) == 1 {
				parts = append(parts, val[0])
				urlParams.Del(key)
				continue
			}
		}

		panic(fmt.Sprintf("Param %s missing in provided data or has multiple values", key))
	}
	var query string
	if len(urlParams) > 0 {
		query = fmt.Sprintf("?%s", urlParams.Encode())
	}
	// if path without / in the beginning will not work but probably no one does
	return fmt.Sprintf("%s%s%s", m.baseUrl, strings.Join(parts, "/"), query)
	
}
