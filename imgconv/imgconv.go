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

// Imgconv is used to store options of CLI.
type Imgconv struct {
	// from is image format before conversion
	from string

	// to is image format after conversion
	to string
}

var supportedFormats = map[string]int{
	"png": 1,
	"jpg": 1,
	"gif": 1,
}

// SupportedFormats returns comma separated string of supported image formats.
func SupportedFormats() string {
	var formats []string
	for k, _ := range supportedFormats {
		formats = append(formats, k)
	}
	return strings.Join(formats, ", ")
}

// NewImgconv allocates a new Imgconv struct and detect error.
func NewImgconv(from, to string) (*Imgconv, error) {
	if _, ok := supportedFormats[from]; !ok {
		return &Imgconv{}, fmt.Errorf("from:%s is not supported", from)
	}
	if _, ok := supportedFormats[to]; !ok {
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

		if err = c.Convert(src, out); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Convert executes image conversion a source file to the dest file.
func (c *Imgconv) Convert(src, out string) error {
	dest := c.buildDestPath(src, out)

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

	switch c.to {
	case "png":
		err = png.Encode(w, img)
	case "jpg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	case "gif":
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	}
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

