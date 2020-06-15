// +build live

package asset

import (
	//"github.com/omeid/go-resources/live"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// FS represents the virtual file system provided by go-resource.
var FS = dir("./assets")

// The following code is necessary until omeid/go-resources#20 is merged

// Resources describes an instance of the go-resources which is an extension of
// http.FileSystem
type Resources interface {
	http.FileSystem
	String(string) (string, bool)
}

func dir(dir string) Resources {

	filename, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir = filepath.Join(filepath.Dir(filename), dir)
	return &resources{http.Dir(dir)}
}

type resources struct {
	http.FileSystem
}

func (r *resources) String(name string) (string, bool) {

	file, err := r.Open(name)
	if err != nil {
		return "", false
	}

	content, err := ioutil.ReadAll(file)

	if err != nil {
		return "", false
	}

	return string(content), true
}
