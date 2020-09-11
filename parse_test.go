package lionrouter

import "testing"

func TestSplit(t *testing.T) {
	split := splitPath("")

	if len(split) != 1 {
		t.Errorf("Wrong length, length is: %d", len(split))
	}

	split = splitPath("/")

	if len(split) != 1 {
		t.Errorf("Wrong length, length is: %d", len(split))
	}

	split = splitPath("/foo/bar/:hello/*world")

	if len(split) != 4 {
		t.Errorf("Wrong length, length is: %d", len(split))

		for _, entry := range split {
			t.Log(entry)
		}
	}

	split = splitPath("/foo/bar/:hello/*world/")

	if len(split) != 4 {
		t.Errorf("Wrong length, length is: %d", len(split))

		for _, entry := range split {
			t.Log(entry)
		}
	}

	split = splitPath("/foo/bar/:hello/*world/")

	if len(split) != 4 {
		t.Errorf("Wrong length, length is: %d", len(split))

		for _, entry := range split {
			t.Log(entry)
		}
	}
}

func TestParse(t *testing.T) {
	parsed := parsePath("/")

	if len(parsed) != 1 {
		t.Errorf("Parsing error, length is %d", len(parsed))
	} else if parsed[0].name != "" {
		t.Errorf("Parsing error, key name is %s, should be empty", parsed[0].name)
	}

	parsed = parsePath("/hello/world/:foo/*bar")

	if len(parsed) != 4 {
		t.Errorf("Parsing error, length is %d", len(parsed))
	} else {
		expected := [4]pathKey{
			{"hello", false, false},
			{"world", false, false},
			{"foo", true, false},
			{"bar", false, true},
		}

		for i, k := range parsed {
			if k.name != expected[i].name || k.wildcard != expected[i].wildcard || k.wildcardPath != expected[i].wildcardPath {
				t.Errorf("Parsing error, name: %s, wildcard: %t, wildcardPath: %t", k.name, k.wildcard, k.wildcardPath)
			}
		}
	}
}
