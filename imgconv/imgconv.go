package imgconv

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// formatInspecter inspects supported image format.
type formatInspecter interface {
	Inspect(string) bool
}

// Formats is the list of registered image formats.
type Formats []string

// Inspect returns true value when image format is supported.
func (f *Formats) Inspect(file string) bool {
	for _, format := range *f {
		if format == strings.TrimLeft(filepath.Ext(file), ".") {
			return true
		}
	}
	return false
}

// SourceFormats is the list of supported source formats.
var SourceFormats = Formats{"png", "jpg", "gif"}

// DestFormats is the list of supported destination formats.
var DestFormats = Formats{"png", "jpg", "gif"}

// Convert executes image conversion a source file to the destination file.
func Convert(src, dest string) error {
	if !SourceFormats.Inspect(src) {
		return fmt.Errorf("src:%v is not supported", src)
	}
	if !DestFormats.Inspect(dest) {
		return fmt.Errorf("dest:%s is not supported", dest)
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	switch filepath.Ext(dest) {
	case ".png":
		err = png.Encode(w, img)
	case ".jpg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	case ".gif":
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	}
	if err != nil {
		return err
	}
	fmt.Printf("convert %v to %v\n", src, dest)
	return nil
}

// Imgconv is used to store options of CLI.
type Imgconv struct {
	// from is image format before conversion
	from string

	// to is image format after conversion
	to string
}

// NewImgconv allocates a new Imgconv struct and detect error.
func NewImgconv(from, to string) (*Imgconv, error) {
	if !SourceFormats.Inspect(from) {
		return &Imgconv{}, fmt.Errorf("from:%s is not supported", from)
	}
	if !DestFormats.Inspect(to) {
		return &Imgconv{}, fmt.Errorf("to:%s is not supported", to)
	}
	if from == to {
		return &Imgconv{}, fmt.Errorf("same formats are specified")
	}
	return &Imgconv{from, to}, nil
}

// Do executes image conversion for target files.
func (c *Imgconv) ConvertRecursively(in, out string) error {
	for _, dir := range []string{in, out} {
		stat, err := os.Stat(dir)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("%s is not directory", dir)
		}
	}

	err := filepath.Walk(in, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if err = Convert(src, c.buildDestPath(src, out)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// buildDestPath creates the destination file path.
func (c *Imgconv) buildDestPath(src, out string) string {
	destPath := strings.Split(src, "/")
	basename := filepath.Base(src)
	destPath[0] = out
	destPath[len(destPath)-1] = strings.TrimSuffix(basename, filepath.Ext(basename)) + "." + c.to
	dest := filepath.Join(destPath...)

	destDir := filepath.Dir(dest)
	if _, err := os.Stat(destDir); err != nil {
		os.MkdirAll(destDir, os.ModePerm)
	}
	return dest
}
