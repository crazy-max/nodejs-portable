//go:generate go get -v github.com/josephspurrier/goversioninfo/...
//go:generate goversioninfo -icon=res/app.ico
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"github.com/crazy-max/nodejs-portable/app/app"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/menu"
	"github.com/crazy-max/nodejs-portable/app/nodejs"
	"github.com/crazy-max/nodejs-portable/app/util"
	"github.com/fatih/color"
	version "github.com/mcuadros/go-version"
	"github.com/op/go-logging"
	"golang.org/x/sys/windows/registry"
)

const (
	NODEJS_VERSION_URL = "http://nodejs.org/dist/v%s/"
	NODEJS_MSI_PATTERN = "node-v%s-%s.msi"
)

// logger
var (
	log       = logging.MustGetLogger(app.ID)
	logFormat = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05} %{level:.4s} - %{message}`)
)

// libs
var (
	libLessmsi app.Lib
)

// paths
var (
	CurrentPath, _      = filepath.Abs(filepath.Dir(os.Args[0]))
	WorkPath            = fs.Join(CurrentPath, "work")
	NpmPath             = fs.Join(CurrentPath, "node_modules", "npm")
	TmpPath             = fs.Join(os.Getenv("USERPROFILE"), ".nodejs-portable")
	LibsPath            = fs.Join(TmpPath, "libs")
	ExtractPath         = fs.Join(TmpPath, "extract")
	NpmGlobalConfigPath = fs.Join(NpmPath, "npmrc")
)

func init() {
	// set window title
	exec.Command("cmd", "/c", fmt.Sprintf("title %s %s", app.NAME, app.VERSION)).Run()

	// log file
	CurrentPath = fs.FormatWinPath(CurrentPath)
	logfile, err := fs.OpenFile(path.Join(CurrentPath, app.ID+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("Log file:", err)
	}

	// init logger
	logBackendFile := logging.NewBackendFormatter(logging.NewLogBackend(logfile, "", 0), logFormat)
	logging.SetBackend(logBackendFile)

	// libs
	libLessmsi = app.Lib{
		Url:        "https://github.com/activescott/lessmsi/releases/download/v1.6.1/lessmsi-v1.6.1.zip",
		Dest:       fs.RemoveUnc(fs.Join(LibsPath, "lessmsi.zip")),
		OutputPath: fs.RemoveUnc(fs.Join(LibsPath, "lessmsi")),
		Exe:        fs.RemoveUnc(fs.Join(LibsPath, "lessmsi", "lessmsi.exe")),
	}

	log.Info("--------")
	log.Infof("Starting %s %s...", app.NAME, app.VERSION)
	log.Info("Current path:", CurrentPath)
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
	if _, err := fs.Stat(fs.Join(CurrentPath, "node.exe")); err == nil {
		util.PrintErrorStr("Node.js is already installed...")
		return nil
	}

	// Clean tmp folder
	fs.CreateSubfolder(TmpPath)

	// seek latest node.js version
	latestNodejs, err := nodejs.GetLatestVersion()
	if err != nil {
		util.PrintErrorStr(fmt.Sprintf("Seeking latest Node.js version: %s...", err.Error()))
		return nil
	}

	// input version and arch
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

	// check version
	fmt.Println()
	util.Print(fmt.Sprintf("Checking Node.js version '%s'... ", nodejsVersion))
	versionUrl := fmt.Sprintf(NODEJS_VERSION_URL, nodejsVersion)
	_, err = http.Get(versionUrl)
	if err != nil {
		util.PrintError(fmt.Errorf("Version %s does not exist in %s", nodejsVersion, versionUrl))
		return nil
	}
	util.PrintOk()

	// check if msi exists
	msiFilename := fmt.Sprintf(NODEJS_MSI_PATTERN, nodejsVersion, nodejsArch)
	msiFile := fs.Join(TmpPath, msiFilename)
	if _, err := fs.Stat(msiFile); err == nil {
		log.Info("Removing existing msi:", msiFile)
		fs.Remove(msiFile)
	}

	// download msi
	msiUrl := versionUrl + msiFilename
	util.Print(fmt.Sprintf("Downloading %s...", msiUrl))
	if err := util.DownloadFile(TmpPath, msiUrl); err != nil {
		fmt.Print(" ")
		util.PrintError(err)
		return nil
	}
	fmt.Print(" ")
	util.PrintOk()

	// extract msi
	if err := util.DownloadLib(libLessmsi); err != nil {
		return nil
	}
	util.Print(fmt.Sprintf("Extracting %s... ", msiFilename))
	fs.RemoveAll(ExtractPath)
	lessmsi := exec.Command(libLessmsi.Exe)
	lessmsi.Dir = fs.RemoveUnc(TmpPath)
	lessmsi.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CmdLine:       fmt.Sprintf(` x "%v" "%v\"`, msiFilename, fs.RemoveUnc(ExtractPath)),
		CreationFlags: 0,
	}
	if err := lessmsi.Run(); err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// copy nodejs folder
	util.Print("Moving nodejs folder... ")
	err = fs.CopyDir(fs.Join(ExtractPath, "SourceDir", "nodejs"), CurrentPath)
	if err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// prepare
	prepare()

	return nil
}

func shell(args ...string) error {
	log.Info(">> SHELL")
	fmt.Println()

	// prepare
	prepare()

	// check if installed
	util.Print("Checking if Node.js installed... ")
	if _, err := fs.Stat(path.Join(CurrentPath, "node.exe")); err != nil {
		util.PrintErrorStr("Not installed...")
		return nil
	}
	util.PrintOk()

	// create launch script
	util.Print("Creating launch script... ")
	launchScript := fs.Join(TmpPath, "launch.bat")
	launchScriptFile, err := fs.Create(launchScript)
	if err != nil {
		util.PrintError(err)
		return nil
	}
	defer launchScriptFile.Close()
	_, err = launchScriptFile.WriteString(util.GetLaunchScriptContent(CurrentPath))
	if err = launchScriptFile.Sync(); err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	// seek Git path
	util.Print("Seeking Git path... ")
	gitPath := ""
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Git_is1`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		gitPath, _, err = key.GetStringValue("InstallLocation")
		if err == nil {
			log.Info("Git installation found:", gitPath)
			os.Setenv("PATH", fmt.Sprintf("%s;%s", fs.FormatWinPath(gitPath+"cmd"), os.Getenv("PATH")))
			util.PrintOk()
		} else {
			util.PrintError(err)
		}
	} else {
		util.PrintError(err)
	}

	// add node to path
	util.Print("Adding node to path... ")
	if err := os.Setenv("PATH", fmt.Sprintf("%s;%s", CurrentPath, os.Getenv("PATH"))); err != nil {
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
		Dir:   fs.RemoveUnc(WorkPath),
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

func prepare() {
	// relocate and edit NPM global config file
	util.Print("Editing NPM global config file... ")
	npmGlobalConfigFile, err := fs.Create(NpmGlobalConfigPath)
	if err != nil {
		util.PrintError(err)
		return
	}
	defer npmGlobalConfigFile.Close()
	util.PrintOk()

	// write
	util.Print("Writing NPM global config file... ")
	globalconfig := fs.FormatWinPath(path.Join(NpmPath, "npmrc"))
	globalignorefile := fs.FormatWinPath(path.Join(NpmPath, "npmignore"))
	initmodule := fs.FormatWinPath(path.Join(NpmPath, "init.js"))
	cache := fs.FormatWinPath(path.Join(NpmPath, "cache"))
	_, err = npmGlobalConfigFile.WriteString(fmt.Sprintf("prefix = %s\n", CurrentPath))
	_, err = npmGlobalConfigFile.WriteString(fmt.Sprintf("globalconfig = %s\n", globalconfig))
	_, err = npmGlobalConfigFile.WriteString(fmt.Sprintf("globalignorefile = %s\n", globalignorefile))
	_, err = npmGlobalConfigFile.WriteString(fmt.Sprintf("initmodule = %s\n", initmodule))
	_, err = npmGlobalConfigFile.WriteString(fmt.Sprintf("cache = %s", cache))
	if err = npmGlobalConfigFile.Sync(); err != nil {
		util.PrintError(err)
		return
	}
	util.PrintOk()

	// create misc files and dir
	util.Print("Creating misc files and dirs... ")
	fs.CreateSubfolder(WorkPath)
	fs.CreateSubfolder(cache)
	fs.OpenFile(globalignorefile, os.O_RDONLY|os.O_CREATE, 0666)
	util.PrintOk()
}
