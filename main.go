//go:generate go get -v github.com/jteeuwen/go-bindata/go-bindata/...
//go:generate go-bindata -pkg bindata -o app/bindata/bindata.go nodejs-portable.conf
//go:generate go get -v github.com/josephspurrier/goversioninfo/...
//go:generate goversioninfo -icon=res/app.ico
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"errors"
	"strings"

	"github.com/crazy-max/nodejs-portable/app/app"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/menu"
	"github.com/crazy-max/nodejs-portable/app/nodejs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	"github.com/crazy-max/nodejs-portable/app/util"
	"github.com/fatih/color"
	version "github.com/mcuadros/go-version"
	"github.com/op/go-logging"
	"golang.org/x/sys/windows/registry"
)

// logger
var (
	log       = logging.MustGetLogger(app.ID)
	logFormat = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05} %{level:.4s} - %{message}`)
)

func init() {
	// set window title
	exec.Command("cmd", "/c", fmt.Sprintf("title %s %s", app.NAME, app.VERSION)).Run()

	// log file
	pathu.CurrentPath = fs.FormatWinPath(pathu.CurrentPath)
	logfile, err := fs.OpenFile(fs.Join(pathu.CurrentPath, app.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("Log file:", err)
	}

	// init logger
	logBackendFile := logging.NewBackendFormatter(logging.NewLogBackend(logfile, "", 0), logFormat)
	logging.SetBackend(logBackendFile)

	log.Info("--------")
	log.Infof("Starting %s %s...", app.NAME, app.VERSION)
	log.Info("Current path:", pathu.CurrentPath)
}

func main() {
	color.New(color.FgHiWhite).Println(app.NAME + " " + app.VERSION)
	color.New(color.FgHiWhite).Println(app.URL)

	// check for update
	latestVersion, err := util.GetLatestVersion()
	if err != nil {
		log.Error("Cannot contact the update server:", err.Error())
		color.New(color.FgRed).Printf("\n%s can't contact the update server: %s\n", app.NAME, err.Error())
	} else if version.Compare(app.VERSION, latestVersion, "<") {
		log.Info("New release available:", latestVersion)
		color.New(color.FgHiGreen).Print("\nA new release is available : ")
		color.New(color.FgHiGreen, color.Bold).Print(latestVersion)
		color.New(color.FgHiGreen).Print("\nDownload : ")
		color.New(color.FgHiGreen, color.Bold).Print(app.URL + "/releases/latest\n")
	}

	// build menu
	menuCommands := []menu.CommandOption{
		{
			Description: "Install",
			Function:    install,
		},
		{
			Description: "Shell",
			Function:    shell,
		},
	}

	menuOptions := menu.NewOptions("Menu", "'menu' for help> ", 0, "")

	menuN := menu.NewMenu(menuCommands, menuOptions)
	menuN.Start()
}

func install(args ...string) error {
	log.Info(">> INSTALL")
	fmt.Println()

	// check if already installed
	if _, err := fs.Stat(fs.Join(pathu.CurrentPath, "node.exe")); err == nil {
		util.PrintErrorStr("Node.js is already installed...")
		return nil
	}

	// create tmp folder
	fs.CreateSubfolder(pathu.TmpPath)

	// seek latest node.js version
	latestNodejs, err := nodejs.GetLatestVersion()
	if err != nil {
		util.PrintErrorStr(fmt.Sprintf("Seeking latest Node.js version: %s...", err.Error()))
		return nil
	}

	// input version
	nodejsVersion := util.ReadLine(fmt.Sprintf("  Version (default %s): ", latestNodejs))
	if nodejsVersion == "" {
		nodejsVersion = latestNodejs
	}
	log.Info("nodejsVersion:", nodejsVersion)

	// input arch
	nodejsArch := util.ReadLine("  Architecture (default x86): ")
	if nodejsArch == "" {
		nodejsArch = "x86"
	}
	log.Info("nodejsArch:", nodejsArch)

	// check dist
	fmt.Println()
	util.Print(fmt.Sprintf("Checking Node.js version '%s'... ", nodejsVersion))
	distURL, distFilename, err := nodejs.GetDistURL(nodejsVersion, nodejsArch)
	if err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// download dist
	util.Print(fmt.Sprintf("Downloading %s...", distURL))
	distPath := fs.Join(pathu.TmpPath, distFilename)
	if err := util.DownloadFile(pathu.TmpPath, distURL); err != nil {
		fmt.Print(" ")
		util.PrintError(err)
		return nil
	}
	fmt.Print(" ")
	util.PrintOk()

	// extract dist
	util.Print(fmt.Sprintf("Extracting %s... ", distFilename))
	extractPath, err := nodejs.ExtractDist(distPath)
	if err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// move nodejs folder
	util.Print("Moving nodejs folder... ")
	err = fs.CopyDir(extractPath, pathu.CurrentPath)
	if err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// create config
	util.Print("Creating Node.js configuration... ")
	if err := nodejs.CreateConfig(); err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	return nil
}

func shell(args ...string) error {
	log.Info(">> SHELL")
	fmt.Println()

	// check if installed
	util.Print("Checking if Node.js installed... ")
	if _, err := fs.Stat(path.Join(pathu.CurrentPath, "node.exe")); err != nil {
		util.PrintErrorStr("Not installed...")
		return nil
	}
	util.PrintOk()

	// create tmp folder
	fs.CreateSubfolder(pathu.TmpPath)

	// create config
	util.Print("Creating Node.js configuration... ")
	if err := nodejs.CreateConfig(); err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// create launch script
	util.Print("Creating launch script... ")
	launchScript := fs.Join(pathu.TmpPath, "launch.bat")
	err := util.CreateFile(launchScript, app.GetLaunchScriptContent(pathu.CurrentPath))
	if err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// Seeking git path
	gitPath, err := getGitPath()
	if err != nil {
		util.PrintError(err)
	} else {
		util.PrintOk()
	}

	// adding git to PATH if found
	if gitPath != "" {
		util.Println("Git found: " + gitPath)
		util.Print("Adding Git to PATH... ")
		os.Setenv("PATH", fmt.Sprintf("%s;%s", gitPath+`\cmd`, os.Getenv("PATH")))
		util.PrintOk()
	}

	// Seeking python path
	pythonPath, err := getPythonPath()
	if err != nil {
		util.PrintError(err)
	} else {
		util.PrintOk()
	}

	// adding python to PATH if found
	if pythonPath != "" {
		util.Println("Python found: " + gitPath)
		util.Print("Adding Python to PATH... ")
		os.Setenv("PATH", fmt.Sprintf("%s;%s", pythonPath, os.Getenv("PATH")))
		util.PrintOk()
	}

	// add node to path
	util.Print("Adding node to path... ")
	if err := os.Setenv("PATH", fmt.Sprintf("%s;%s;%s", fs.FormatWinPath(pathu.CurrentPath+`\node_modules\npm`), pathu.CurrentPath, os.Getenv("PATH"))); err != nil {
		util.PrintError(err)
	}
	util.PrintOk()

	// set NODE_PATH
	util.Print("Setting NODE_PATH... ")
	if err := os.Setenv("NODE_PATH", fs.FormatWinPath(path.Join(pathu.CurrentPath, "node_modules"))); err != nil {
		util.PrintError(err)
	}
	util.PrintOk()

	// clear screen
	util.Println("Clearing screen...")
	clear := exec.Command("cmd", "/c", "cls")
	clear.Stdout = os.Stdout
	clear.Run()

	// transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   fs.RemoveUnc(pathu.WorkPath),
		Sys: &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf(` /k "%s"`, fs.RemoveUnc(launchScript)),
		},
	}

	// start up a new shell.
	log.Info("Starting up the shell... ")
	proc, err := os.StartProcess(os.Getenv("COMSPEC"), []string{}, &pa)
	if err != nil {
		util.PrintError(err)
		return nil
	}

	// wait until user exits the shell
	if _, err = proc.Wait(); err != nil {
		util.PrintError(err)
		return nil
	}

	return nil
}

func getGitPath() (string, error) {
	gitPath := ""
	if _, err := os.Stat(app.Conf.GitPath); err == nil {
		gitPath = app.Conf.GitPath
	} else {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Git_is1`, registry.QUERY_VALUE)
		if err == nil {
			defer key.Close()
			gitRegPath, _, err := key.GetStringValue("InstallLocation")
			if err == nil {
				gitPath = gitRegPath
			}
		}
	}
	if gitPath != "" {
		if _, err := os.Stat(path.Join(gitPath, "cmd", "git.exe")); err != nil {
			return "", errors.New("git.exe not found in " + path.Join(gitPath, "cmd"))
		}
	}
	return gitPath, nil
}

func getPythonPath() (string, error) {
	pythonPath := ""
	if _, err := os.Stat(app.Conf.PythonPath); err == nil {
		pythonPath = app.Conf.PythonPath
	}
	if pythonPath != "" {
		if _, err := os.Stat(path.Join(pythonPath, "python.exe")); err != nil {
			return "", errors.New("python.exe not found in " + pythonPath)
		}
	}
	return strings.TrimRight(fs.FormatWinPath(pythonPath), `\`), nil
}
