package r2router

import (
	"fmt"
	"strings"
)

type routeNode struct {
	paramNode  bool
	paramName  string
	path       string
	children   []*routeNode
	paramChild *routeNode
	handler    Handler
	routePath  string
}

func newRouteNode() *routeNode {
	r := &routeNode{}
	r.children = make([]*routeNode, 0)
	return r
}

// findChild finds a child node that matches the given node
// It returns nil if no node found. This is to see if
// we already have a similar node registered
func (n *routeNode) findChild(nn *routeNode) *routeNode {
	for _, c := range n.children {
		if c.path == nn.path {
			// same node
			return c
		}
	}
	return nil
}

// insertChild registers given node in the route node tree
// If there is already a similar node it will not insert new node
// The returned node is always the registered one ie either
// newly registered or the old one
func (n *routeNode) insertChild(nn *routeNode) *routeNode {
	if child := n.findChild(nn); child != nil {
		return child
	}

	if n.paramChild != nil && nn.paramNode {
		// only allow one param child, unique param name
		if n.paramChild.paramName != nn.paramName {
			panic("Param name must be same for")
		}
		return n.paramChild
	}
	if nn.paramNode {
		n.paramChild = nn
	} else {
		n.children = append(n.children, nn)
	}
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
			if path == "" {
				break
			}
			child := newRouteNode()
			token, path = nextPath(path)
			//fmt.Println(token, path)
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
	params.appData = make(map[interface{}]interface{})
	params.requestParams = make(map[string]string)
	if path != "" {
		// can be better by getting one at the time
		var matched bool
		var token string
		route := n.root
		for {
			if path == "" {
				break
			}
			token, path = nextPath(path)
			matched = false
			for _, c := range route.children {
				if c.path == token {
					route = c
					matched = true
					break
				}
			}

			if !matched {
				if route.paramChild != nil {
					route = route.paramChild
					matched = true
					params.requestParams[route.paramName] = token
					continue
				}
				return nil, nil, ""
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
			s += fmt.Sprintf(" (<%p>)", &node.handler)
		}
		s += "\n"

		for _, c := range node.children {
			s += dumNode(c, ident+1)
		}
		if node.paramChild != nil {
			s += dumNode(node.paramChild, ident+1)
		}
		return s
	}

	return dumNode(n.root, 0)
}
