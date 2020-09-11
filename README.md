# LionRouter (Golang path-trie router)

This is a lightweight minimalistic implementation of a go router based on a path trie approach. Don't expect a fully featured router neither now nor in future. Only the most necessary features will be included, based on our needs.

If you need some features which are not included, feel free to add them in fork. 

# Supported functionality

- implements the default go `http.Handler` interface
- routes
    - static routing
    - parameter routes
        - Named keys
        - Named paths
- Sub-Routing using the go `http.Handler` interface
- Middleware stack using the default wrapper technique through `func(next http.Handler) http.Handler`

# How to use

The router implements the `http.Handler` interface, so it can be used alongside with other frameworks. Each route handler is using the `http.Handler` interface itself to maximize compatibility with other packages.

## New router instance
`New() *Router`
```
router := lionrouter.New()
```

## Define a route
Routes can easily be assigned by using one of the following functions.

`Get(path string, handler http.Handler)`

`Post(path string, handler http.Handler)`

`Put(path string, handler http.Handler)`

`Patch(path string, handler http.Handler)`

`Delete(path string, handler http.Handler)`

`Head(path string, handler http.Handler)`

`Options(path string, handler http.Handler)`

### Static route

```
router.Get("/contact", contactHandler())
```

Each request to `http://.../contact` will be passed through the contactHandler() method.

### Parameter routes

#### Named key

```
router.Get("/download/:file", downloadHandler())
```

Each request to `http://.../download/someFileName` will be passed through the downloadHandler() method where you can get the corresponding value passed by as `file` by extracing it from the `context`.

```
params, ok := lionrouter.Params(r.Context())

if ok {
    w.Write([]byte(params["file"]))
}
```
or
```
w.Write([]byte(lionrouter.Param(r.Context(), "file"))
```

This will generate the output: `someFileName`.
A request like `http://.../download/profile_picture.png` would get: `profile_picture.png`.

Multiple named keys per route are possible:

```
router.Handler(http.MethodGet, "/download/:param1/:param2", downloadHandler())
```

#### Named path

Named path is corresponding to the named key functionality except a whole path is read from request.

```
router.Get("/download/*file", downloadHandler())
```

Each request to `http://.../download/some/path/here/test.png` will be passed through the downloadHandler() method where you can get the corresponding value passed by as `file` by extracing it from the `context`.

```
params, ok := lionrouter.Params(r.GetContext())

if ok {
    w.Write([]byte(params["file"]))
}
```

This will generate the output: `/some/path/here/test.png`.

A named path can also be combined with named keys:

```
router.Get("/download/:param1/:param2/*param3", downloadHandler())
```
**Named paths are only possible at the end and cannot be used alongside with the subrouting funciotnality.**

## Sub-Routing

You can register any http.Handler to a given path. The request will then be passed through the registered Handler. Note that the http.StripPrefix() middleware will be chained before passing through the request to the sub-router handler.

`Router(path string, handler http.Handler)`

```
mainRouter := lionrouter.New()
staticRouter := lionrouter.New()

staticRouter.Get("/*file", staticHandler())
mainRouter.Route("/static", staticRouter)
```

**Path to sub-routers doesn't support named paths or keys at this moment**

## Middleware

To assign any given middleware of the type `func(http.Handler) http.Handler`, just use the `Middleware(...func(http.Handler) Handler)` method.

```
router.Use(someMiddleware)
router.Use(otherMiddleware)
```
this adds `someMiddleware` and `otherMiddleware`. You can also define the middleware at once.
```
router.Use(someMiddleware, otherMiddleware)
```

# License

MIT licensed 2017-2019 Cedrik Kaufmann. See the LICENSE file for further details.
