package lionrouter

import (
	"github.com/oxtoacart/bpool"
	"net/http"
)

type trie struct {
	root  *node
	bpool *bpool.BufferPool
}

func newTrie() *trie {
	return &trie{
		root:  newNode(&pathKey{"", false, false}),
		bpool: bpool.NewBufferPool(48),
	}
}

func (t *trie) walkAdd(path string) (*node, error) {
	// parse path
	parsed := parsePath(path)

	// set cursor tree root
	cursor := t.root

	// walk tree
	for _, key := range parsed {
		// if router is set, don't walk tree
		if cursor.router != nil {
			return nil, ErrAssignment
		}

		// check if key is wildcard or wildcard path
		if key.wildcard || key.wildcardPath {
			// create wildcard node if nil
			if cursor.wildcard == nil {
				cursor.wildcard = newNode(key)
			}

			// walkAdd tree using wildcard node
			cursor = cursor.wildcard

			if key.wildcardPath {
				break
			}
		} else {
			// create node map for children if nil
			if cursor.children == nil {
				cursor.children = make(map[string]*node)
			}

			// create node entry in map if nil
			if cursor.children[key.name] == nil {
				cursor.children[key.name] = newNode(key)
			}

			// walkAdd tree using node map
			cursor = cursor.children[key.name]
		}
	}

	return cursor, nil
}

func (t *trie) walkGet(path string) (*node, RouterParam) {
	// split path
	parsed := splitPath(path)

	// set cursor tree root
	cursor := t.root

	param := make(RouterParam)

	for i, key := range parsed {
		// if router is set, don't walk tree
		if cursor.router != nil {
			break
		}

		// check for wildcard child
		if cursor.wildcard != nil {
			cursor = cursor.wildcard

			// check whether wildcard is single key or whole path
			if cursor.key.wildcardPath {
				// get buffer from pool
				buffer := t.bpool.Get()

				// build wildcard path
				for _, s := range parsed[i:] {
					buffer.WriteRune('/')
					buffer.WriteString(s)
				}

				// set wildcard path in param map
				param[cursor.key.name] = buffer.String()

				// put back buffer to pool
				t.bpool.Put(buffer)

				// break walking the tree
				break
			} else {
				// set current wildcard key in param map
				param[cursor.key.name] = key
			}
		} else if cursor.children != nil && cursor.children[key] != nil {
			// access child through key name
			cursor = cursor.children[key]
		} else {
			// return because no node found
			return nil, nil
		}
	}

	// return the node and param map
	return cursor, param
}

func (t *trie) addHandler(method string, path string, handler http.Handler) error {
	// check if handler not nil
	if handler == nil {
		return ErrNilHandler
	}

	// walk tree
	cursor, err := t.walkAdd(path)

	// check for errors
	if err != nil {
		return err
	}

	// check again if router is set and return if set
	if cursor.router != nil {
		return ErrAssignment
	}

	// create leaf if nil
	if cursor.leaf == nil {
		cursor.leaf = newLeaf()
	}

	// add handler to leaf
	err = cursor.leaf.addHandler(method, handler)

	// check for errors
	if err != nil {
		return err
	}

	// handler added, no errors
	return nil
}

func (t *trie) addRouter(path string, router http.Handler) error {
	// check if handler not nil
	if router == nil {
		return ErrNilHandler
	}

	// walk tree
	cursor, err := t.walkAdd(path)

	// check for errors
	if err != nil {
		return err
	}

	// check if no leaf/wildcard/children/router are set and return if set
	if cursor.leaf != nil || cursor.wildcard != nil || cursor.children != nil || cursor.router != nil {
		return ErrAssignment
	}

	// add router
	cursor.router = router

	// router added, no errors
	return nil
}

func (t *trie) get(method string, path string) (http.Handler, RouterParam) {
	// retrieve node from trie
	node, param := t.walkGet(path)

	// check if node not nil
	if node == nil {
		return nil, nil
	}

	// check if sub-router exists
	if node.router != nil {
		return node.router, param
	}

	// check if node has leaf and retrieve handler from leaf
	if node.leaf != nil {
		h, err := node.leaf.getHandler(method)

		if err != nil {
			return nil, nil
		}

		return h, param
	}

	// return node
	return nil, nil
}
