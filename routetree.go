package r2router

import (
	"fmt"
	"strings"
)

type routeNode struct {
	paramNode bool
	paramName string
	path      string
	cchildren []*routeNode
	children  []*routeNode
	handler   Handler
	routePath string
}

func newRouteNode() *routeNode {
	r := &routeNode{}
	r.children = make([]*routeNode, 0)
	r.cchildren = make([]*routeNode, 0)
	return r
}

// findChild finds a child node that matches the given node
// It returns nil if no node found. This is to see if
// we already have a similar node registered
func (n *routeNode) findChild(nn *routeNode) *routeNode {
	for _, c := range n.children {
		if c.paramNode && nn.paramNode {
			// both are paramNode
			return c
		} else if c.path == nn.path {
			// same node
			return c
		}
	}
	return nil
}

func (n *routeNode) swapChild(i, j int) {
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

// insertChild registers given node in the route node tree
// If there is already a similar node it will not insert new node
// The returned node is always the registered one ie either
// newly registered or the old one
func (n *routeNode) insertChild(nn *routeNode) *routeNode {
	if child := n.findChild(nn); child != nil {
		if child.paramNode && nn.paramNode {
			// only allow one param child, unique param name
			if child.paramName != nn.paramName {
				panic("Param name must be same for")
			}
		}

		return child
	}

	n.children = append(n.children, nn)
	if len(n.children) > 1 {
		if n.children[len(n.children)-2].paramNode {
			n.swapChild(len(n.children)-2, len(n.children)-1)
		}
	}
	return nn
}

func (n *routeNode) insertCChild(nn *routeNode) *routeNode {
	for _, child := range n.cchildren {
		if child.path == nn.path {
			return child
		}
	}
	n.cchildren = append(n.cchildren, nn)
	return nn
}

type rootNode struct {
	root    *routeNode
	handler Handler
}

func newRouteTree() *rootNode {
	r := &rootNode{}
	r.root = newRouteNode()
	return r
}

func nextPath(path string) (string, string) {
	i := strings.Index(path, "/")
	if i == -1 {
		return path, ""
	}
	return path[:i], path[i+1:]
}

func (n *rootNode) addRoute(path string, handler Handler) {
	path = strings.Trim(path, "/")
	if path != "" {
		// Start with the roots
		parent := n.root
		var token string
		for {
			child := newRouteNode()

			if !strings.Contains(path, ":") {
				child.path = path
				child.routePath = fmt.Sprintf("%s/%s", parent.routePath, path)
				parent = parent.insertCChild(child)
				break
			}

			token, path = nextPath(path)

			if token[:1] == ":" {
				// param type
				child.paramName = strings.TrimSpace(token[1:])
				if child.paramName == "" {
					panic("Param name can not be empty")
				}
				child.paramNode = true
			} else {
				child.path = token
			}
			// will be parent for the next path token
			child.routePath = fmt.Sprintf("%s/%s", parent.routePath, token)
			parent = parent.insertChild(child)
			if path == "" {
				break
			}
		}
		// adding handler
		if parent.handler != nil {
			panic(fmt.Sprintf("'%s' has already a handler", path))
		}
		parent.handler = handler

	} else {
		if n.handler != nil {
			panic("There is already a handler for /")
		}
		n.handler = handler
	}
}

func (n *rootNode) match(path string) (Handler, Params, string) {
	path = strings.Trim(path, "/")
	params := &params_{}
	params.appData = make(map[string]interface{})
	params.requestParams = make(map[string]string)
	if path != "" {
		// can be better by getting one at the time
		var matched bool
		var token string
		route := n.root
		for {
			// comparing constant paths
			for _, c := range route.cchildren {
				if c.path == path {
					return c.handler, params, c.routePath
				}
			}
			token, path = nextPath(path)
			matched = false
			for _, c := range route.children {
				if c.path == token {
					route = c
					matched = true
					break
				}
				if c.paramNode {
					route = c
					matched = true
					params.requestParams[c.paramName] = token
					break
				}
			}
			if !matched {
				return nil, nil, ""
			}
			if path == "" {
				break
			}
		}
		return route.handler, params, route.routePath
	} else {
		return n.handler, params, "/"
	}
}

func (n *rootNode) dump() string {
	var dumNode func(node *routeNode, ident int) string

	dumNode = func(node *routeNode, ident int) string {
		s := ""
		identing := ""
		for i := 0; i < ident; i++ {
			identing += " "
		}
		s += identing + " |\n"
		identing += "  "
		if node.paramNode {
			s += identing + "-- :" + node.paramName
		} else {
			s += identing + "-- " + node.path
		}
		if node.handler != nil {
			s += fmt.Sprintf(" (<%p>)", node.handler)
		}
		s += "\n"
		for _, c := range node.cchildren {
			s += dumNode(c, ident+1)
		}

		for _, c := range node.children {
			s += dumNode(c, ident+1)
		}

		return s
	}

	return dumNode(n.root, 0)
}
