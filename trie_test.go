package lionrouter

import (
	"net/http"
	"testing"
)

func TestTrieAdd(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing
		})
	}

	trie := newTrie()

	routes := [6]string{
		"/",
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/foo",
		"/foo/bar/world",
	}

	// test if first assign works
	for _, r := range routes {
		err := trie.addHandler(http.MethodGet, r, testHandler())

		if err != nil {
			t.Errorf("tree: failed adding handler for '%s' --> %s", r, err.Error())
		}
	}

	// test if reassign fails
	for _, r := range routes {
		err := trie.addHandler(http.MethodGet, r, testHandler())

		if err == nil {
			t.Errorf("tree: success reassign handler for '%s' --> %s", r, err.Error())
		}
	}
}

func TestTrieAddRouter(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing
		})
	}

	trie := newTrie()

	routes := [5]string{
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/super/:world",
		"/world/:super/:foo",
	}

	// test if first assign works
	for _, r := range routes {
		err := trie.addRouter(r, testHandler())

		if err != nil {
			t.Errorf("tree: failed adding router for '%s' --> %s", r, err.Error())
		}
	}

	// test if reassign fails
	for _, r := range routes {
		err := trie.addRouter(r, testHandler())

		if err == nil {
			t.Errorf("tree: success reassign router for '%s' --> %s", r, err.Error())
		}
	}
}

func TestTrieGet(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing
		})
	}

	trie := newTrie()

	routes := [9]string{
		"/",
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/foo",
		"/bar/foo/world",
		"/foo/bar/:user/:domain",
		"/foo/bar/:user/:domain/world",
		"/wild/:domain/*world",
	}

	routesRetr := [9]string{
		"/",
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/foo",
		"/bar/foo/world",
		"/foo/bar/peter/test.de",
		"/foo/bar/max/google.de/world",
		"/wild/foo.de/bar/hello.jpg",
	}

	routeParam := make([]RouterParam, 3)

	for i := range routeParam {
		routeParam[i] = make(RouterParam)
	}

	routeParam[0]["user"] = "peter"
	routeParam[0]["domain"] = "test.de"
	routeParam[1]["user"] = "max"
	routeParam[1]["domain"] = "google.de"
	routeParam[2]["domain"] = "foo.de"
	routeParam[2]["world"] = "/bar/hello.jpg"

	for _, r := range routes {
		err := trie.addHandler(http.MethodGet, r, testHandler())

		if err != nil {
			t.Errorf("tree: failed adding handler for '%s' --> %s", r, err.Error())
		}
	}

	for i, r := range routesRetr {
		handler, param := trie.get(http.MethodGet, r)

		if handler == nil {
			t.Errorf("tree: error while retrieving handler '%s'", r)
		}

		if param != nil && i > 5 && i < 8 {
			if routeParam[i-6]["user"] != param["user"] {
				t.Errorf("wrong param 'user': param is '%s' and should be '%s'", param["user"], routeParam[i-6]["user"])
			}

			if routeParam[i-6]["domain"] != param["domain"] {
				t.Errorf("wrong param 'domain': param is '%s' and should be '%s'", param["domain"], routeParam[i-6]["domain"])
			}
		}

		if param != nil && i == 8 {
			if routeParam[2]["domain"] != param["domain"] {
				t.Errorf("wrong param 'domain': param is '%s' and should be '%s'", param["domain"], routeParam[2]["domain"])
			}

			if routeParam[2]["wild"] != param["wild"] {
				t.Errorf("wrong param 'world': param is '%s' and should be '%s'", param["world"], routeParam[2]["world"])
			}
		}
	}
}

func TestTrieGetRouter(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing
		})
	}

	trie := newTrie()

	routes := [4]string{
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/world/:foo/:bar",
	}

	routesRetr := [4]string{
		"/foo/bar",
		"/bar/foo/",
		"/hello",
		"/world/test1/test2",
	}

	routeParam := make([]RouterParam, 3)

	for i := range routeParam {
		routeParam[i] = make(RouterParam)
	}

	foo := "test1"
	bar := "test2"

	for _, r := range routes {
		err := trie.addRouter(r, testHandler())

		if err != nil {
			t.Errorf("tree: failed adding router for '%s' --> %s", r, err.Error())
		}
	}

	for i, r := range routesRetr {
		handler, param := trie.get(http.MethodGet, r)

		if handler == nil {
			t.Errorf("tree: error while retrieving handler '%s'", r)
		}

		if param != nil && i == 4 {
			if param["foo"] != foo {
				t.Errorf("wrong param 'foo': param is '%s' and should be '%s'", param["foo"], foo)
			}

			if param["bar"] != bar {
				t.Errorf("wrong param 'bar': param is '%s' and should be '%s'", param["bar"], bar)
			}
		}
	}
}
