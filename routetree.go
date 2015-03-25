package r2router

import (
	"fmt"
	"strings"
)

type routeNode struct {
	paramNode bool
	paramName string
	path      string
	children  []*routeNode
	handler   Handler
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

type rootNode struct {
	root    *routeNode
	handler Handler
}

func newRouteTree() *rootNode {
	r := &rootNode{}
	r.root = newRouteNode()
	return r
}

func (n *rootNode) addRoute(path string, handler Handler) {
	path = strings.Trim(path, "/")
	if path != "" {
		paths := strings.Split(path, "/")
		// Start with the roots
		parent := n.root
		for _, p := range paths {
			//fmt.Println(p)
			child := newRouteNode()

			if strings.Contains(p, ":") {
				// param type
				child.paramName = strings.TrimSpace(p[1:])
				if child.paramName == "" {
					panic("Param name can not be empty")
				}
				child.paramNode = true
			} else {
				child.path = p
			}
			// will be parent for the next path token
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

func (n *rootNode) match(path string) (Handler, Params) {
	path = strings.Trim(path, "/")
	if path != "" {
		// can be better by getting one at the time
		paths := strings.Split(path, "/")
		var matched bool
		route := n.root
		params := &params_{}
		params.appData = make(map[string]interface{})
		params.requestParams = make(map[string]string)
		for _, p := range paths {
			matched = false
			//fmt.Println("p:", p)
			for _, c := range route.children {
				if c.path == p {
					//fmt.Println("match:", c.path)
					route = c
					matched = true
					break
				}
				if c.paramNode {
					route = c
					matched = true
					params.requestParams[c.paramName] = p
					break
				}
			}
			if !matched {
				break
			}
		}
		if matched {
			return route.handler, params
		}

	} else {
		if n.handler != nil {
			return n.handler, nil
		}
	}

	return nil, nil
}
