package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/crazy-max/WindowsSpyBlocker/app/bindata"
	"github.com/crazy-max/WindowsSpyBlocker/app/utils/file"
	"github.com/crazy-max/WindowsSpyBlocker/app/utils/pathu"
	"github.com/crazy-max/WindowsSpyBlocker/app/utils/print"
)

// App constants
const (
	NAME    = "WindowsSpyBlocker"
	VERSION = "4.8.2"
	PACKAGE = "github.com/crazy-max/WindowsSpyBlocker"
	URL     = "https://" + PACKAGE
)

// Config
var (
	App appConf
)

// Lib structure
type Lib struct {
	Url        string `json:"url"`
	Checksum   string `json:"checksum"`
	Dest       string
	OutputPath string
	Checkfile  string
}

type appConf struct {
	Version   string `json:"version"`
	Debug     bool   `json:"debug"`
	Proxifier struct {
		LogPath string `json:"logPath"`
	} `json:"proxifier"`
	Sysmon struct {
		EvtxPath string `json:"evtxPath"`
	} `json:"sysmon"`
	Wireshark struct {
		PcapngPath string `json:"pcapngPath"`
		Capture    struct {
			Interface int    `json:"interface"`
			Filter    string `json:"filter"`
		} `json:"capture"`
	} `json:"wireshark"`
	Exclude struct {
		Ips   []string `json:"ips"`
		Hosts []string `json:"hosts"`
		Orgs  []string `json:"orgs"`
	} `json:"exclude"`
}

func init() {
	var err error
	var appOld appConf

	cfgPath := path.Join(pathu.Current, "app.conf")

	// Load default config
	defaultConf, err := bindata.Asset("app.conf")
	if err != nil {
		err = fmt.Errorf("Cannot load asset app.conf: %s", err.Error())
		print.QuitFatal(err)
	}
	err = json.Unmarshal(defaultConf, &App)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall defaultConf: %s", err.Error())
		print.QuitFatal(err)
	}
	newVersion := App.Version

	// Create conf if not exists
	if _, err := os.Stat(cfgPath); err != nil {
		err = ioutil.WriteFile(cfgPath, defaultConf, 0644)
		if err != nil {
			err = fmt.Errorf("Cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.Current), err.Error())
			print.QuitFatal(err)
		}
	}

	// Load current config
	raw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		err = fmt.Errorf("Cannot read %s: %s", strings.TrimLeft(cfgPath, pathu.Current), err.Error())
		print.QuitFatal(err)
	}
	err = json.Unmarshal(raw, &appOld)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.Current), err.Error())
		print.QuitFatal(err)
	}
	oldVersion := appOld.Version

	// Perform upgrade if different version
	if newVersion != oldVersion {
		if err := performUpgrade(); err != nil {
			print.QuitFatal(err)
		}
	}

	// Merge config
	err = json.Unmarshal(raw, &App)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshall %s: %s", strings.TrimLeft(cfgPath, pathu.Current), err.Error())
		print.QuitFatal(err)
	}
	App.Version = newVersion

	// Write config
	cfgJson, _ := json.MarshalIndent(App, "", "  ")
	if err != nil {
		err = fmt.Errorf("Cannot marshal config: %s", err.Error())
		print.QuitFatal(err)
	}
	err = ioutil.WriteFile(cfgPath, cfgJson, 0644)
	if err != nil {
		err = fmt.Errorf("Cannot write file %s: %s", strings.TrimLeft(cfgPath, pathu.Current), err.Error())
		print.QuitFatal(err)
	}
}

func performUpgrade() error {
	// Remove content of libs folder
	if err := file.RemoveContents(pathu.Libs); err != nil {
		return err
	}

	return nil
}
