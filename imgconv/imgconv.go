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

type Imgconv struct {
	in      string
	out     string
	from    string
	to      string
	verbose bool
}

var Supported = map[string]int{
	"png": 1,
	"jpg": 1,
	"gif": 1,
}

func SupportedFormats() string {
	var formats []string
	for k, _ := range Supported {
		formats = append(formats, k)
	}
	return strings.Join(formats, ", ")
}

func NewImgconv(in, out, from, to string, verbose bool) (*Imgconv, error) {
	stat, err := os.Stat(in)
	if err != nil {
		return &Imgconv{}, err
	}
	if !stat.IsDir() {
		return &Imgconv{}, fmt.Errorf("in:%s is not directory", in)
	}
	stat, err = os.Stat(out)
	if err != nil {
		return &Imgconv{}, err
	}
	if !stat.IsDir() {
		return &Imgconv{}, fmt.Errorf("out:%s is not directory", out)
	}
	if _, ok := Supported[from]; !ok {
		return &Imgconv{}, fmt.Errorf("from:%s is not supported", from)
	}
	if _, ok := Supported[to]; !ok {
		return &Imgconv{}, fmt.Errorf("to:%s is not supported", to)
	}
	if from == to {
		return &Imgconv{}, fmt.Errorf("same formats are specified")
	}
	return &Imgconv{in, out, from, to, verbose}, nil
}

func (c *Imgconv) Do() error {
	err := filepath.Walk(c.in, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != fmt.Sprintf(".%s", c.from) {
			c.vLog("format is not match %s: %s", path, c.from)
			return nil
		}

		err = c.convert(path)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *Imgconv) convert(src string) error {
	destPath := strings.Split(src, "/")
	basename := filepath.Base(src)
	destPath[0] = c.out
	destPath[len(destPath)-1] = strings.Replace(basename, c.from, c.to, 1)
	dest := filepath.Join(destPath...)

	destDir := filepath.Dir(dest)
	if _, err := os.Stat(destDir); err != nil {
		os.MkdirAll(destDir, os.ModePerm)
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	switch c.to {
	case "png":
		err = png.Encode(out, img)
	case "jpg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
	case "gif":
		err = gif.Encode(out, img, &gif.Options{NumColors: 256})
	}
	if err != nil {
		return err
	}
	c.vLog("convert %s to %s", src, dest)
	return nil
}

func (c *Imgconv) vLog(format string, a ...interface{}) {
	if !c.verbose {
		return
	}

	s := fmt.Sprintf(format, a...)
	if strings.HasSuffix(s, "\n") {
		fmt.Print(s)
	} else {
		fmt.Println(s)
	}
}
