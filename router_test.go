package router

import "testing"

func TestMatches(t *testing.T) {
	cases := map[string][]string{
		"/one/two/three": []string{"/one/two/three", "/one/two/*", "/one/*", Default, Error},
	}
	for value, expected := range cases {
		result := matches(value)
		if len(result) != len(expected) {
			t.Errorf("Expected %+v but got %+v for case %q", expected, result, value)
		}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Expected %q but got %q at position %d for case %q", expected[i], v, i, value)
			}
		}
	}
}

var r []string

func BenchmarkMatches(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r = matches("/one/two/three")
	}
}