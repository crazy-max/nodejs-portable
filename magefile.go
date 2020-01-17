// +build mage

package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default mage target
var Default = Build

var (
	binPath   = path.Join("bin")
	buildPath = path.Join(binPath, "nodejs-portable.exe")
	buildEnv  = map[string]string{
		"GO111MODULE": "on",
		"GOPROXY":     "https://goproxy.io",
		"GOOS":        "windows",
		"GOARCH":      "386",
		"CGO_ENABLED": "0",
	}
)

// Build Run go build
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Generate)

	var args []string
	args = append(args, "build", "-o", buildPath, "-v")
	args = append(args, "-ldflags", flags())

	fmt.Println("⚙️ Go build...")
	if err := sh.RunWith(buildEnv, mg.GoCmd(), args...); err != nil {
		return err
	}

	return nil
}

// Clean Remove files generated at build-time
func Clean() error {
	if err := createDir(binPath); err != nil {
		return err
	}
	if err := cleanDir(binPath); err != nil {
		return err
	}
	return nil
}

// Download Run go mod download
func Download() error {

	fmt.Println("⚙️ Go mod download...")
	if err := sh.RunWith(buildEnv, mg.GoCmd(), "mod", "download"); err != nil {
		return err
	}

	return nil
}

// Generate Run go generate
func Generate() error {
	mg.Deps(Download)
	mg.Deps(appConf)
	mg.Deps(versionInfo)

	fmt.Println("⚙️ Go generate...")
	if err := sh.RunV(mg.GoCmd(), "generate", "-v"); err != nil {
		return err
	}

	return nil
}

// flags returns ldflags
func flags() string {
	//hash := hash()
	tag := tag()
	mod := mod()
	return fmt.Sprintf(`-s -w -X "main.version=%s" -X "main.module=%s"`, tag, mod)
}

// mod returns module name
func mod() string {
	f, err := os.Open("go.mod")
	if err == nil {
		reader := bufio.NewReader(f)
		line, _, _ := reader.ReadLine()
		return strings.Replace(string(line), "module ", "", 1)
	}
	return ""
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("bash", "-c", "git describe --abbrev=0 --tags 2> /dev/null")
	if s == "" {
		return "0.0.0"
	}
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// appConf generates app.conf file
func appConf() error {
	fmt.Println("🔨 Generating nodejs-portable.conf...")

	var tpl = template.Must(template.New("").Parse(`{
  "version": "{{ .Version }}",
  "immediateMode": false,
  "shell": "cmd",
  "workPath": "./work",
  "customPaths": [
    "C:/Program Files (x86)/Git/cmd",
    "D:/another_path",
    "../a_relative_path"
  ],
  "clearScreen": true
}`))

	f, err := os.Create("nodejs-portable.conf")
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, struct {
		Version string
	}{
		Version: tag(),
	})
}

// versionInfo generates versioninfo.json
func versionInfo() error {
	fmt.Println("🔨 Generating versioninfo.json...")

	var tpl = template.Must(template.New("").Parse(`{
	"FixedFileInfo":
	{
		"FileFlagsMask": "3f",
		"FileFlags ": "00",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":
	{
		"Comments": "",
		"CompanyName": "",
		"FileDescription": "Make Node.js portable on Windows",
		"FileVersion": "{{ .Version }}.0",
		"InternalName": "",
		"LegalCopyright": "https://{{ .Package }}",
		"LegalTrademarks": "",
		"OriginalFilename": "nodejs-portable.exe",
		"PrivateBuild": "",
		"ProductName": "Node.js Portable",
		"ProductVersion": "{{ .Version }}.0",
		"SpecialBuild": ""
	},
	"VarFileInfo":
	{
		"Translation": {
			"LangID": "0409",
			"CharsetID": "04B0"
		}
	}
}`))

	f, err := os.Create("versioninfo.json")
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, struct {
		Package string
		Version string
	}{
		Package: mod(),
		Version: tag(),
	})
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 777)
	}
	return nil
}

func cleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
