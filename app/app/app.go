package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"path/filepath"

	"github.com/crazy-max/nodejs-portable/app/bindata"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	"github.com/crazy-max/nodejs-portable/app/util"
)

// Config
var (
	Conf ConfStruct
)

// ConfStruct the conf structure
type ConfStruct struct {
	Version       string   `json:"version"`
	ImmediateMode bool     `json:"immediateMode"`
	Shell         string   `json:"shell"`
	WorkPath      string   `json:"workPath"`
	CustomPaths   []string `json:"customPaths"`
	ClearScreen   bool     `json:"clearScreen"`
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

	// Check work path
	newWorkPath, err := filepath.Abs(Conf.WorkPath)
	if err != nil {
		err = fmt.Errorf("cannot get absolute path of work path '%s': %s", Conf.WorkPath, err.Error())
		util.QuitFatal(err)
	} else {
		Conf.WorkPath = newWorkPath
		if err := os.MkdirAll(Conf.WorkPath, 777); err != nil {
			err = fmt.Errorf("cannot create work path folder in '%s': %s", Conf.WorkPath, err.Error())
			util.QuitFatal(err)
		}
	}
}

// GetLaunchScriptContent is executed while launching shell
func GetLaunchScriptContent() string {
	workPath := fs.FormatWinPath(fs.RemoveUnc(Conf.WorkPath))
	nodePath := fs.FormatWinPath(fs.RemoveUnc(pathu.AppPath))

	launchScriptTpl := `@ECHO OFF
SETLOCAL EnableDelayedExpansion

SET nodejsPath=@NODEJS_PATH@
SET nodejsWork=@WORK_PATH@
SET npmPath=%nodejsPath%\node_modules\npm
SET npmGlobalConfigFilePath=%npmPath%\npmrc

SET PATH=%nodejsPath%;%PATH%
cd "%nodejsWork%"
"%nodejsPath%\nodevars.bat"
"%nodejsPath%\npm.cmd" config set globalconfig "%npmGlobalConfigFilePath%" --global`

	launchScript := strings.Replace(launchScriptTpl, "@WORK_PATH@", workPath, -1)
	launchScript = strings.Replace(launchScript, "@NODEJS_PATH@", nodePath, -1)

	return launchScript
}
