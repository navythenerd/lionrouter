package lionrouter

import "net/http"

type node struct {
	key *pathKey

	leaf     *leaf
	wildcard *node
	children map[string]*node

	router http.Handler
}

func newNode(key *pathKey) *node {
	return &node{
		key: key,
	}
}
