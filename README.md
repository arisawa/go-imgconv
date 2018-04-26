# imgconv
--
    import "github.com/arisawa/go-imgconv/imgconv"


## Usage

```go
var DestFormats = Formats{"png", "jpg", "gif"}
```
DestFormats is the list of supported destination formats.

```go
var SourceFormats = Formats{"png", "jpg", "gif", "webp"}
```
SourceFormats is the list of supported source formats.

#### func  Convert

```go
func Convert(src, dest string) error
```
Convert executes image conversion a source file to the destination file.

#### type Formats

```go
type Formats []string
```

Formats is the list of registered image formats.

#### func (*Formats) Inspect

```go
func (f *Formats) Inspect(file string) bool
```
Inspect returns true value when image format is supported.

#### type RecursiveConverter

```go
type RecursiveConverter struct {
}
```

RecursiveConverter converts target images recursively.

#### func  NewRecursiveConverter

```go
func NewRecursiveConverter(in, out, srcFormat, destFormat string) (*RecursiveConverter, error)
```
NewRecursiveConverter allocates a new RecursiveConverter struct and detect
error.

#### func (*RecursiveConverter) Convert

```go
func (rc *RecursiveConverter) Convert() error
```
Convert executes image conversion for target files.

#### func (*RecursiveConverter) GetTargets

```go
func (rc *RecursiveConverter) GetTargets() []*target
```
GetTargets returns property of targets.
