package lionrouter

import "testing"

func TestNewNode(t *testing.T) {
	parsed := parsePath("/foo/:bar/*world")

	for _, k := range parsed {
		n := newNode(k)

		if n == nil {
			t.Error("Error creating new node.")
		} else if n.key != k {
			t.Error("Key mismatch")
		}
	}
}
