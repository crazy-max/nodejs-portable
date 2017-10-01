package nodejs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	"github.com/crazy-max/nodejs-portable/app/util"
)

const (
	distJSON = "https://nodejs.org/dist/index.json"
	zipURL   = "http://nodejs.org/dist/v%s/node-v%s-win-%s.zip"
	msiURL   = "http://nodejs.org/dist/v%s/node-v%s-%s.msi"
)

var (
	libLessmsi util.Lib
)

// Version structure
type Version struct {
	Version string `json:"version"`
}

// Versions structure
type Versions []Version

func init() {
	libLessmsi = util.Lib{
		URL:        "https://github.com/activescott/lessmsi/releases/download/v1.6.1/lessmsi-v1.6.1.zip",
		Dest:       fs.RemoveUnc(fs.Join(pathu.LibsPath, "lessmsi.zip")),
		OutputPath: fs.RemoveUnc(fs.Join(pathu.LibsPath, "lessmsi")),
		Exe:        fs.RemoveUnc(fs.Join(pathu.LibsPath, "lessmsi", "lessmsi.exe")),
	}
}

// GetLatestVersion returns the latest version of Node.js
func GetLatestVersion() (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(distJSON)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		var versions Versions
		json.Unmarshal(bodyBytes, &versions)
		return strings.TrimLeft(versions[0].Version, "v"), nil
	}

	return "", fmt.Errorf("Status code %d", response.StatusCode)
}

// GetDistURL returns the node.js distribution binary url
func GetDistURL(version string, arch string) (string, string, error) {
	zipFullURL := fmt.Sprintf(zipURL, version, version, arch)
	msiFullURL := fmt.Sprintf(msiURL, version, version, arch)
	if util.URLValid(zipFullURL) {
		return zipFullURL, filepath.Base(zipFullURL), nil
	}
	if util.URLValid(msiFullURL) {
		return msiFullURL, filepath.Base(msiFullURL), nil
	}
	return "", "", fmt.Errorf("Version %s (%s) does not exist", version, arch)
}

// ExtractDist extracts the node.js distribution binary msi or zip file in the extract folder
func ExtractDist(distPath string) (string, error) {
	if filepath.Ext(distPath) == ".zip" {
		return extractZipDist(distPath)
	}
	return extractMsiDist(distPath)
}

func extractZipDist(distPath string) (string, error) {
	fs.RemoveAll(pathu.ExtractPath)
	if err := fs.Unzip(fs.RemoveUnc(distPath), fs.RemoveUnc(pathu.ExtractPath)); err != nil {
		return "", err
	}
	subfolder := strings.TrimSuffix(filepath.Base(distPath), filepath.Ext(filepath.Base(distPath)))
	return fs.Join(pathu.ExtractPath, subfolder), nil
}

func extractMsiDist(distPath string) (string, error) {
	if err := util.DownloadLib(libLessmsi); err != nil {
		return "", err
	}
	fs.RemoveAll(pathu.ExtractPath)
	lessmsi := exec.Command(libLessmsi.Exe)
	lessmsi.Dir = fs.RemoveUnc(filepath.Dir(distPath))
	lessmsi.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       fmt.Sprintf(` x "%v" "%v\"`, filepath.Base(distPath), fs.RemoveUnc(pathu.ExtractPath)),
		CreationFlags: 0,
	}
	if err := lessmsi.Run(); err != nil {
		return "", err
	}
	return fs.Join(pathu.ExtractPath, "SourceDir", "nodejs"), nil
}

// CreateConfig init node.js configuration and folders
func CreateConfig() error {
	nodePath := fs.FormatWinPath(fs.RemoveUnc(pathu.CurrentPath))
	nodeCachePath := fs.Join(pathu.CurrentPath, "cache")
	nodeEtcPath := fs.Join(pathu.CurrentPath, "etc")
	npmrcPath := fs.Join(nodeEtcPath, "npmrc")
	npmIgnorePath := fs.Join(nodeEtcPath, ".npmignore")
	npmPath := fs.Join(pathu.CurrentPath, "node_modules", "npm")
	npmrcGlobalPath := fs.Join(npmPath, "npmrc")

	fs.CreateSubfolder(nodeCachePath)
	fs.CreateSubfolder(nodeEtcPath)
	fs.CreateSubfolder(pathu.WorkPath)

	err := util.CreateFile(npmrcGlobalPath, strings.Replace(`prefix = @CURRENT_PATH@\`, "@CURRENT_PATH@", nodePath, -1))
	if err != nil {
		return err
	}

	err = util.CreateFile(npmrcPath, strings.Replace(`cache = @CURRENT_PATH@\cache
globalconfig = @CURRENT_PATH@\etc\npmrc
globalignorefile = @CURRENT_PATH@\etc\.npmignore
init-module = @CURRENT_PATH@\etc\.npm-init.js
userconfig = @CURRENT_PATH@\etc\npmrc`, "@CURRENT_PATH@", nodePath, -1))
	if err != nil {
		return err
	}

	_, err = fs.OpenFile(npmIgnorePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	return nil
}
