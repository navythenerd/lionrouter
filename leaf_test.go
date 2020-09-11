package lionrouter

import (
	"net/http"
	"testing"
)

func TestNewLeaf(t *testing.T) {
	testMethods := [7]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	testLeaf := newLeaf()

	for _, m := range testMethods {
		handler, err := testLeaf.getHandler(m)

		if err != nil {
			t.Error(err)
		} else if handler != nil {
			t.Errorf("Handler for method '%s' should be nil.", m)
		}
	}
}

func TestLeafHandler(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing
		})
	}

	testMethods := [7]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	testLeaf := newLeaf()

	for _, m := range testMethods {
		err := testLeaf.addHandler(m, testHandler())

		if err != nil {
			t.Error(err)
		} else {
			handler, err := testLeaf.getHandler(m)

			if err != nil {
				t.Error(err)
			} else if handler == nil {
				t.Errorf("Handler for method '%s' should not be nil.", m)
			}
		}
	}

	testLeaf = newLeaf()
	err := testLeaf.addHandler(testMethods[6], testHandler())

	if err != nil {
		t.Error(err)
	} else {
		handler, err := testLeaf.getHandler(testMethods[0])

		if err != nil {
			t.Error(err)
		} else if handler != nil {
			t.Errorf("Handler for method '%s' should be nil.", testMethods[0])
		}

		handler, err = testLeaf.getHandler(testMethods[6])

		if err != nil {
			t.Error(err)
		} else if handler == nil {
			t.Errorf("Handler for method '%s' should not be nil.", testMethods[6])
		}
	}
}
