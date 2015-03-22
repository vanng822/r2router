package r2router

import (
	//"fmt"
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
		if n.children[len(n.children)-1].paramNode {
			n.swapChild(len(n.children)-1, len(n.children))
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
	// /group/:id
	// /group/:id/action
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	var parent *routeNode
	if len(paths) > 0 {
		parent = n.root
		for _, p := range paths {
			//fmt.Println(p)
			child := newRouteNode()

			if strings.Contains(p, ":") {
				// param type
				child.paramName = p[1:]
				child.paramNode = true
			} else {
				child.path = p
			}
			parent = parent.insertChild(child)
		}
		// adding handler
		parent.handler = handler
	} else if path == "/" {
		if n.handler != nil {
			panic("There is already a handler for /")
		}
		n.handler = handler
	}
}

func (n *rootNode) match(path string) (Handler, Params) {
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	if len(paths) > 0 {
		var matched bool
		route := n.root
		params := make(Params)
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
					params[c.paramName] = p
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
	} else if path == "/" {
		if n.handler != nil {
			return n.handler, nil
		}
	}

	return nil, nil
}
