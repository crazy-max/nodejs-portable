package app

import (
	"strings"

	"github.com/crazy-max/nodejs-portable/app/fs"
)

// App infos
const (
	ID      = "nodejs-portable"
	NAME    = "Node.js Portable"
	VERSION = "2.1.0"
	PACKAGE = "github.com/crazy-max/nodejs-portable"
	URL     = "https://" + PACKAGE
)

// Lib structure
type Lib struct {
	URL        string
	Dest       string
	OutputPath string
	Exe        string
}

// GetLaunchScriptContent is executed while launching shell
func GetLaunchScriptContent(nodejsPath string) string {
	return strings.Replace(`@ECHO OFF
SETLOCAL EnableDelayedExpansion

SET nodejsPath=@CURRENT_PATH@
SET nodejsWork=%nodejsPath%\work
SET npmPath=%nodejsPath%\node_modules\npm
SET npmGlobalConfigFilePath=%npmPath%\npmrc

SET PATH=%nodejsPath%;%PATH%
cd "%nodejsWork%"
"%nodejsPath%\nodevars.bat"
"%nodejsPath%\npm.cmd" config set globalconfig "%npmGlobalConfigFilePath%" --global`,
		"@CURRENT_PATH@", fs.FormatWinPath(fs.RemoveUnc(nodejsPath)), -1)
}
