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
var (
	Id      = "nodejs-portable"
	Name    = "Node.js Portable"
	Package = "github.com/crazy-max/nodejs-portable"
	Url     = "https://" + Package
	Version = "1.0"
)

// Config
var (
	Conf ConfStruct
)

// ConfStruct the conf structure
type ConfStruct struct {
	Version       string   `json:"version"`
	ImmediateMode bool     `json:"immediateMode"`
	CustomPaths   []string `json:"customPaths"`
}

func init() {
	var err error
	var oldConf ConfStruct

	cfgPath := path.Join(pathu.CurrentPath, "nodejs-portable.conf")

	// Load default config
	defaultConf, err := bindata.Asset("nodejs-portable.conf")
	if err != nil {
		err = fmt.Errorf("cannot load asset nodejs-portable.conf: %s", err.Error())
		util.QuitFatal(err)
	}
	err = json.Unmarshal(defaultConf, &Conf)
	if err != nil {
		err = fmt.Errorf("cannot unmarshall defaultConf: %s", err.Error())
		util.QuitFatal(err)
	}

	// Create conf if not exists
	if _, err := os.Stat(cfgPath); err != nil {
		err = ioutil.WriteFile(cfgPath, defaultConf, 0644)
		if err != nil {
			err = fmt.Errorf("cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
			util.QuitFatal(err)
		}
	}

	// Load current config
	raw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		err = fmt.Errorf("cannot read %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}
	err = json.Unmarshal(raw, &oldConf)
	if err != nil {
		err = fmt.Errorf("cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}

	// Merge config and force version
	oldConf.Version = Conf.Version
	err = json.Unmarshal(raw, &Conf)
	if err != nil {
		err = fmt.Errorf("cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
		util.QuitFatal(err)
	}

	// Write config
	cfgJSON, _ := json.MarshalIndent(Conf, "", "  ")
	if err != nil {
		err = fmt.Errorf("cannot marshal config: %s", err.Error())
		util.QuitFatal(err)
	}
	err = ioutil.WriteFile(cfgPath, cfgJSON, 0644)
	if err != nil {
		err = fmt.Errorf("cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.CurrentPath), err.Error())
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
