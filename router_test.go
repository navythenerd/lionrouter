package lionrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	handlerBody := []byte("TEST_HANDLER")
	errBody := []byte("404_NOT_FOUND")

	// create router instance
	router := New()

	// define handler
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(handlerBody)
		})
	}

	errHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(errBody)
		})
	}

	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next.ServeHTTP(w, r)
		})
	}

	router.NotFoundHandler = errHandler()

	errRoutes := [5]string{
		"/abc/xyz",
		"/xyz",
		"/test.png",
		"/static/css/app.css",
		"/static/img/logo.png",
	}

	// define some routes
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

	router.Use(testMiddleware)

	for _, r := range routes {
		router.Get(r, testHandler())
	}

	// start server
	startServer := func() *http.Server {
		srv := &http.Server{Addr: ":8080", Handler: router}

		go func() {
			// returns ErrServerClosed on graceful close
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop. don't use
				// code with race conditions like these for production. see post
				// comments below on more discussion on how to handle this.
				t.Errorf("ListenAndServe(): %s", err)
			}
		}()

		// returning reference so caller can call Shutdown()
		return srv
	}

	srv := startServer()

	for _, r := range errRoutes {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080%s", r))

		if err != nil {
			t.Errorf("route '%s' call failed", r)
		}

		if status := resp.StatusCode; status != http.StatusNotFound {
			t.Errorf("route '%s' should be 404, server responded with: %d", r, status)
		}

		if middlewareVal := resp.Header.Get("x-test"); middlewareVal != "" {
			t.Errorf("middleware test failed, x-test header is: %s. Middleware should not be used here.", middlewareVal)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("error while reading body: %s", err.Error())
		}

		bodyMismatch := false

		for i, b := range body {
			if b != errBody[i] {
				bodyMismatch = true
			}
		}

		if bodyMismatch {
			t.Error("response mismatch")
		}

		resp.Body.Close()
	}

	for _, r := range routesRetr {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080%s", r))

		if err != nil {
			t.Errorf("route '%s' call failed", r)
		}

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("route '%s' should be 200, server responded with: %d", r, status)
		}

		if middlewareVal := resp.Header.Get("x-test"); middlewareVal != "test" {
			t.Errorf("middleware failed, x-test header is: %s", middlewareVal)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("error while reading body: %s", err.Error())
		}

		bodyMismatch := false

		for i, b := range body {
			if b != handlerBody[i] {
				bodyMismatch = true
			}
		}

		if bodyMismatch {
			t.Error("response mismatch")
		}

		resp.Body.Close()
	}

	srv.Close()
}

func TestSubRouter(t *testing.T) {
	handlerBody := []byte("TEST_HANDLER")

	// create router instance
	router := New()
	subRouter := New()

	// define handler
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(handlerBody)
		})
	}

	// define some routes
	routes := [2]string{
		"/foo/bar",
		"/hello/world/",
	}

	subRoutes := [4]string{
		"/",
		"/hello",
		"/foo/world",
		"/hello/:key",
	}

	routesRetr := [6]string{
		"/foo/bar",
		"/foo/bar/hello",
		"/hello/world/",
		"/hello/world/hello/",
		"/hello/world/foo/world",
		"/hello/world/hello/super",
	}

	for _, r := range subRoutes {
		subRouter.Get(r, testHandler())
	}

	for _, r := range routes {
		router.Route(r, subRouter)
	}

	// start server
	startServer := func() *http.Server {
		srv := &http.Server{Addr: ":8080", Handler: router}

		go func() {
			// returns ErrServerClosed on graceful close
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				// NOTE: there is a chance that next line won't have time to run,
				// as main() doesn't wait for this goroutine to stop. don't use
				// code with race conditions like these for production. see post
				// comments below on more discussion on how to handle this.
				t.Errorf("ListenAndServe(): %s", err)
			}
		}()

		// returning reference so caller can call Shutdown()
		return srv
	}

	srv := startServer()

	for _, r := range routesRetr {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080%s", r))

		if err != nil {
			t.Errorf("route '%s' call failed", r)
		}

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("route '%s' should be 200, server responded with: %d", r, status)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("error while reading body: %s", err.Error())
		}

		bodyMismatch := false

		for i, b := range body {
			if b != handlerBody[i] {
				bodyMismatch = true
			}
		}

		if bodyMismatch {
			t.Error("response mismatch")
		}

		resp.Body.Close()
	}

	srv.Close()
}
