package imgconv_test

import (
	"testing"

	"github.com/arisawa/go-imgconv/imgconv"
)

func TestInspectFormat(t *testing.T) {
	var testFormats = imgconv.Formats{"png", "jpg"}

	testCase := []struct {
		input string
		want  bool
	}{
		{"/path/to/go.png", true},
		{"/path/to/go.jpg", true},
		{"/path/to/go.gif", false},
		{"/path/to/go.webp", false},
	}

	for _, tc := range testCase {
		ret := testFormats.Inspect(tc.input)
		if tc.want != ret {
			t.Errorf("input: %v, want: %v, got %v", tc.input, tc.want, ret)
		}
	}
}
