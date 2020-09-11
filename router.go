package lionrouter

import (
	"context"
	"net/http"
)

const (
	// context iota keys
	contextParams = iota // RouterParams context key
)

// RouterParam
// map[string]string holds the parameter from each route if set
type RouterParam map[string]string

// Router
// struct represents router instance which implements the http.Handler interface
type Router struct {
	// trie holds the routes, a trie is a special string search tree
	trie *trie

	// middleware stack
	middleware []func(http.Handler) http.Handler

	// 404 and 500 handler which called if no handler for route assigned or error occurs
	NotFoundHandler http.Handler
}

// New creates a new router instance
// returns router instance
func New() *Router {
	return &Router{
		trie: newTrie(),
	}
}

// Get assigns the handler as HTTP-GET Route for given path
func (r *Router) Get(path string, handler http.Handler) {
	err := r.handle(http.MethodGet, path, handler)

	if err != nil {
		panic(err)
	}
}

// Post assigns the handler as HTTP-POST Route for given path
func (r *Router) Post(path string, handler http.Handler) {
	err := r.handle(http.MethodPost, path, handler)

	if err != nil {
		panic(err)
	}
}

// Put assigns the handler as HTTP-PUT Route for given path
func (r *Router) Put(path string, handler http.Handler) {
	err := r.handle(http.MethodPut, path, handler)

	if err != nil {
		panic(err)
	}
}

// Patch assigns the handler as HTTP-PATCH Route for given path
func (r *Router) Patch(path string, handler http.Handler) {
	err := r.handle(http.MethodPatch, path, handler)

	if err != nil {
		panic(err)
	}
}

// Delete assigns the handler as HTTP-DELETE Route for given path
func (r *Router) Delete(path string, handler http.Handler) {
	err := r.handle(http.MethodDelete, path, handler)

	if err != nil {
		panic(err)
	}
}

// Head assigns the handler as HTTP-HEAD Route for given path
func (r *Router) Head(path string, handler http.Handler) {
	err := r.handle(http.MethodHead, path, handler)

	if err != nil {
		panic(err)
	}
}

// Options assigns the handler as HTTP-OPTIONS Route for given path
func (r *Router) Options(path string, handler http.Handler) {
	err := r.handle(http.MethodOptions, path, handler)

	if err != nil {
		panic(err)
	}
}

// handler assigns the handler for given method and route
// returns error if assignment fails
func (r *Router) handle(method string, path string, handler http.Handler) error {
	return r.trie.addHandler(method, path, handler)
}

// Route assigns the given handler for given path, if route is called the request is passed through this handler
// Used for sub-routing, path will be stripped through http.StripPrefix middleware before request passed to handler
func (r *Router) Route(path string, handler http.Handler) {
	err := r.trie.addRouter(path, http.StripPrefix(path, handler))

	if err != nil {
		panic(err)
	}
}

// Use assigns a middleware stack to the whole router instance
func (r *Router) Use(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		r.middleware = append(r.middleware, m)
	}
}

// middlewareChain builds the middleware stack
// returns the http.Handler stack
func (r *Router) middlewareChain(next http.Handler) http.Handler {
	for _, m := range r.middleware {
		next = m(next)
	}

	return next
}

// ServeHTTP handles a given request
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get path of current request
	path := req.URL.Path

	// get current context
	ctx := req.Context()

	// get handler from trie
	handler, params := r.trie.get(req.Method, path)

	if handler != nil {
		// check if params are given and add them to the current context
		if params != nil {
			ctx = context.WithValue(ctx, contextParams, params)
		}

		// get middleware chain
		handler = r.middlewareChain(handler)
	} else if r.NotFoundHandler == nil {
		// use default fallback handler for 404 response
		handler = func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Page not found!"))
			})
		}()
	} else {
		// use custom 404 handler
		handler = r.NotFoundHandler
	}

	// pass through the request to the given http.Handler
	handler.ServeHTTP(w, req.WithContext(ctx))
}

// Params extracts the params map from the given context
// It returns the params string map or nil
func Params(ctx context.Context) (RouterParam, bool) {
	// get params from context
	params, ok := ctx.Value(contextParams).(RouterParam)

	// check if params exist in given context and return them
	if ok {
		return params, true
	}

	// return error because no params exist
	return nil, false
}

// Param extracts the single param value from the given context and key
// It returns the param as string or an empty string
func Param(ctx context.Context, key string) string {
	// get params from context
	params, ok := ctx.Value(contextParams).(RouterParam)

	// check if params exist in given context and return them
	if ok {
		return params[key]
	}

	// return empty string because no param exists
	return ""
}
