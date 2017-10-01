package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/crazy-max/nodejs-portable/app/bindata"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	"github.com/crazy-max/nodejs-portable/app/util"
)

// App infos
const (
	ID      = "nodejs-portable"
	NAME    = "Node.js Portable"
	VERSION = "2.2.1"
	PACKAGE = "github.com/crazy-max/nodejs-portable"
	URL     = "https://" + PACKAGE
)

// Config
var (
	Conf ConfStruct
)

// ConfStruct the conf structure
type ConfStruct struct {
	Version   	  	string   `json:"version"`
	ImmediateMode 	bool `json:"immediateMode"`
	CustomPaths 	[]string `json:"customPaths"`
}

func init() {
	var err error
	var oldConf ConfStruct

	cfgPath := path.Join(pathu.CurrentPath, "nodejs-portable.conf")

	// Load default config
	defaultConf, err := bindata.Asset("nodejs-portable.conf")
	if err != nil {
		err = fmt.Errorf("Cannot load asset nodejs-portable.conf: %s", err.Error())
		util.QuitFatal(err)
	}
	err = json.Unmarshal(defaultConf, &Conf)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall defaultConf: %s", err.Error())
		util.QuitFatal(err)
	}

	// Create conf if not exists
	if _, err := os.Stat(cfgPath); err != nil {
		err = ioutil.WriteFile(cfgPath, defaultConf, 0644)
		if err != nil {
			err = fmt.Errorf("Cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
			util.QuitFatal(err)
		}
		util.Print("Created a new config file. Have a look, please.\n")
	}

	// Load current config
	raw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		err = fmt.Errorf("Cannot read %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}
	err = json.Unmarshal(raw, &oldConf)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}

	// Merge config
	err = json.Unmarshal(raw, &Conf)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}

	// Write config
	cfgJSON, _ := json.MarshalIndent(Conf, "", "  ")
	if err != nil {
		err = fmt.Errorf("Cannot marshal config: %s", err.Error())
		util.QuitFatal(err)
	}
	err = ioutil.WriteFile(cfgPath, cfgJSON, 0644)
	if err != nil {
		err = fmt.Errorf("Cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}
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
