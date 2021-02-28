// +build !live

package assets

import "embed"

//go:embed agent base jobs nodes single src tasks
var Content embed.FS
