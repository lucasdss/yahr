package yahr

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request, p url.Values) {

}

func TestYAHR(t *testing.T) {

	rt := New()

	rt.Handler("GET", "/a/b/c", dummyHandler)

	rt.Handler("POST", "/a/b/c", dummyHandler)

	rt.Handler("POST", "/a/b/:id", dummyHandler)
	rt.Handler("POST", "/a/:id/c", dummyHandler)

	node := rt.search(rt.root, strings.Split("/a/b/c", "/")[1:], nil)
	if node == nil {
		t.Fatal("could not find any router")
	}
	if _, ok := node.methods["GET"]; !ok {
		t.Fatal("could not find GET router")
	}

	node = rt.search(rt.root, strings.Split("/a/b/123", "/")[1:], nil)
	if _, ok := node.methods["POST"]; !ok {
		t.Fatal("could not find POST router")
	}

	node = rt.search(rt.root, strings.Split("/a/123/c", "/")[1:], nil)
	if _, ok := node.methods["POST"]; !ok {
		t.Fatal("could not find POST router")
	}
}
