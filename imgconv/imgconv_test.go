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
		if tc.err && err == nil {
			t.Fatalf("should be error but not")
		}
		if _, err := os.Stat(dest); !tc.err && os.IsNotExist(err) {
			t.Fatalf("dest file: %v should be created but not", dest)
		}
		os.Remove(dest)
	}
}

type testTarget struct {
	src, dest string
}

// tp returns path joined from testdata
func tp(path ...string) string {
	root := []string{"..", "testdata"}
	return filepath.Join(append(root, path...)...)
}

func TestNewRecurciveConverter(t *testing.T) {
	t.Helper()

	testCase := []struct {
		in, out, srcFormat, destFormat string
		wantTargets                    []testTarget
		err                            bool
	}{
		{
			in:         tp(""),
			out:        tp("tmp"),
			srcFormat:  "png",
			destFormat: "jpg",
			wantTargets: []testTarget{
				{tp("gopher.png"), tp("tmp", "gopher.jpg")},
				{tp("subdir", "gopher.png"), tp("tmp", "subdir", "gopher.jpg")},
			},
			err: false,
		},
	}

	for _, tc := range testCase {
		rc, err := imgconv.NewRecursiveConverter(tc.in, tc.out, tc.srcFormat, tc.destFormat)
		if !tc.err && err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
		if tc.err && err == nil {
			t.Fatalf("should be error but not")
		}
		for i, target := range rc.GetTargets() {
			wt := tc.wantTargets[i]
			if wt.src != target.GetSrc() {
				t.Fatalf("src file is not match. want: %v, got: %v", wt.src, target.GetSrc())
			}
			if wt.dest != target.GetDest() {
				t.Fatalf("dest file is not match. want: %v, got: %v", wt.dest, target.GetDest())
			}
		}
	}
}
