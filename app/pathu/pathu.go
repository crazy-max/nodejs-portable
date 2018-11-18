package pathu

import (
	"os"
	"path/filepath"

	"github.com/crazy-max/nodejs-portable/app/fs"
)

// Paths
var (
	CurrentPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	AppPath        = fs.Join(CurrentPath, "app")
	TmpPath        = fs.Join(os.Getenv("TEMP"), ".nodejs-portable")
	LibsPath       = fs.Join(TmpPath, "libs")
	ExtractPath    = fs.Join(TmpPath, "extract")
)
