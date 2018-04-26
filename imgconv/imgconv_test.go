package imgconv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arisawa/go-imgconv/imgconv"
)

func TestInspectFormat(t *testing.T) {
	t.Helper()
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

func TestConvert(t *testing.T) {
	t.Helper()
	testCase := []struct {
		srcFormat  string
		destFormat string
		err        bool
	}{
		{"png", "jpg", false},
		{"png", "gif", false},
		{"png", "webp", true},
	}

	for _, tc := range testCase {
		src := filepath.Join("..", "testdata", "gopher."+tc.srcFormat)
		dest := filepath.Join("..", "testdata", "gopher."+tc.destFormat)
		err := imgconv.Convert(src, dest)
		if !tc.err && err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
		if _, err := os.Stat(dest);  !tc.err && os.IsNotExist(err) {
			t.Fatalf("dest file: %v should be created but not", dest)
		}
		if tc.err && err == nil {
			t.Fatalf("should be error but not")
		}
		os.Remove(dest)
	}
}
