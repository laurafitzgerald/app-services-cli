package generate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

// InternalError wraps errors that are development issues and unrelated to user
// input.
type InternalError string

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", string(e))
}

// File wraps os.File. Use this type when generating files that may already
// exist on disk and should be overwritten.
type File struct {
	*os.File
}

// Open first creates dir then opens <dir>/<fileName> for reading and writing,
// creating the file if it does not exist.
func Open(dir, fileName string) (*File, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dir, fileName), os.O_RDWR|os.O_CREATE, 0666)
	return &File{f}, err
}

// WriteObject writes any object to w.
func WriteYAML(w io.Writer, obj interface{}) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return write(w, b)
}

// write writes b to w. If w is a File, its contents will be cleared and w
// will be closed following the write.
func write(w io.Writer, b []byte) error {
	if f, isFile := w.(*File); isFile {
		if err := f.Truncate(0); err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
	}
	_, err := w.Write(b)
	return err
}
