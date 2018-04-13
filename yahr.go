/*
Package yahr implements a simple trie accepting named values
ex.:
rt := yahr.New()
rt.Handler("GET", "/myapp/:id", yahr.Handler)
rt.Handler("POST", "/myapp/:id", yahr.Handler)
rt.Handler("PUT", "/myapp/:id/rating", yahr.Handler)

If you want to understand how trie works:
https://en.wikipedia.org/wiki/Trie

A more complex and elegant implementation
https://github.com/julienschmidt/httprouter

Another example more flexible:
https://github.com/dimfeld/httptreemux
*/
package yahr

import (
	"net/http"
	"net/url"
	"strings"
)

// HandlerFunc defines the type of handlers this router accepts
type HandlerFunc func(http.ResponseWriter, *http.Request, url.Values)

// YAHR must be initialized with root "/" already defined
type YAHR struct {
	root *node
}

// New returns YAHR initialized
func New() *YAHR {
	return &YAHR{
		&node{
			path:     "/",
			methods:  make(map[string]HandlerFunc),
			children: make(map[string]*node),
		},
	}
}

// Node is the tree element
// path is always without slash
// ex.: myapp or :id
type node struct {
	path string

	methods  map[string]HandlerFunc
	children map[string]*node

	// true if path == :id
	isParam bool
}

func (r *YAHR) insert(nd *node, path []string) {
	length := len(path)

	for level := 0; level < length; level++ {
		index := path[level]

		if nd.children[index] == nil {
			var isParam bool
			if index[0] == ':' {
				isParam = true
			}
			nd.children[index] = &node{
				path:     index,
				isParam:  isParam,
				methods:  make(map[string]HandlerFunc),
				children: make(map[string]*node),
			}
		}

		nd = nd.children[index]
	}

}

func (r *YAHR) search(nd *node, path []string, params url.Values) *node {
	length := len(path)

	for level := 0; level < length; level++ {
		index := path[level]
		if nd.children[index] == nil {
			for _, v := range nd.children {
				if v.isParam {
					nd = v
					if params != nil {
						params.Add(v.path[1:], index)
					}
				}
			}
			continue
		}
		nd = nd.children[index]
	}

	return nd
}

// Handler add the user defined HandlerFunc
// for the path and method given
func (r *YAHR) Handler(method, path string, h HandlerFunc) {

	paths := strings.Split(path, "/")[1:]

	r.insert(r.root, paths)

	node := r.search(r.root, paths, nil)

	node.methods[method] = h

}

// ServeHTTP is the function to be called by http.Server
// We must implement this function to route the requests
func (r *YAHR) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// let's populate our reuqest.Form
	req.ParseForm()
	params := req.Form

	// the first element of the slice is empty
	// [1:] removes the first element
	node := r.search(r.root, strings.Split(req.URL.Path, "/")[1:], params)

	handler := node.methods[req.Method]

	if handler != nil {
		handler(w, req, params)
		return
	}

	// The default status is 404 Not Found
	w.WriteHeader(http.StatusNotFound)
}
