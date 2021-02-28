// +build live

package assets

import (
	"os"
)

// Content is an fs.FS that is representative of the live file system.
var Content = os.DirFS("./assets")
