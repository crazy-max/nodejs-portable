//go:generate go get -v github.com/jteeuwen/go-bindata/go-bindata/...
//go:generate go-bindata -pkg bindata -o app/bindata/bindata.go nodejs-portable.conf
//go:generate go get -v github.com/josephspurrier/goversioninfo/...
//go:generate goversioninfo -icon=res/app.ico
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/crazy-max/nodejs-portable/app/app"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/log"
	"github.com/crazy-max/nodejs-portable/app/menu"
	"github.com/crazy-max/nodejs-portable/app/nodejs"
	"github.com/crazy-max/nodejs-portable/app/pathu"
	"github.com/crazy-max/nodejs-portable/app/util"
	"github.com/fatih/color"
	version "github.com/mcuadros/go-version"
)

func init() {
	// set window title
	exec.Command("cmd", "/c", fmt.Sprintf("title %s %s", app.Name, app.Version)).Run()

	log.Logger.Info("--------")
	log.Logger.Info(fmt.Sprintf("Starting %s %s...", app.Name, app.Version))
	log.Logger.Info("Current path:", pathu.CurrentPath)
}

func main() {
	color.New(color.FgHiWhite).Println(app.Name + " " + app.Version)
	color.New(color.FgHiWhite).Println(app.Url)

	// check for update
	latestVersion, err := util.GetLatestVersion()
	if err != nil {
		log.Logger.Error("Cannot contact the update server:", err.Error())
		color.New(color.FgYellow).Printf("\n%s can't contact the update server: %s\n", app.Name, err.Error())
	} else if version.Compare(app.Version, latestVersion, "<") {
		log.Logger.Info("New release available:", latestVersion)
		color.New(color.FgHiGreen).Print("\nA new release is available : ")
		color.New(color.FgHiGreen, color.Bold).Print(latestVersion)
		color.New(color.FgHiGreen).Print("\nDownload : ")
		color.New(color.FgHiGreen, color.Bold).Print(app.Url + "/releases/latest\n")
	}

	// open shell on immediate mode
	if app.Conf.ImmediateMode == true {
		shell()
		return
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
	log.Logger.Info(">> INSTALL")
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
	log.Logger.Info("nodejsVersion:", nodejsVersion)

	// input arch
	nodejsArch := util.ReadLine("  Architecture (default x86): ")
	if nodejsArch == "" {
		nodejsArch = "x86"
	}
	log.Logger.Info("nodejsArch:", nodejsArch)

	// check dist
	fmt.Println()
	util.Print(fmt.Sprintf("Checking Node.js version %s... ", nodejsVersion))
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
	log.Logger.Info(">> SHELL")
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

	// check custom paths
	customPaths := ""
	util.Println("Checking custom paths...")
	for _, customPath := range app.Conf.CustomPaths {
		if customPath == "" {
			continue
		}
		tmpCustomPath, _ := filepath.Abs(customPath)
		util.Print(fmt.Sprintf("# %s ", tmpCustomPath))
		if _, err := os.Stat(tmpCustomPath); err == nil {
			if customPaths != "" {
				customPaths = customPaths + ";"
			}
			customPaths = customPaths + strings.TrimRight(fs.FormatWinPath(tmpCustomPath), `\`)
			util.PrintOk()
		} else {
			util.PrintWarningStr("Not found...")
		}
	}

	// add custom paths
	if customPaths != "" {
		util.Print("Adding customs path to PATH... ")
		os.Setenv("PATH", fmt.Sprintf("%s;%s", customPaths, os.Getenv("PATH")))
		util.PrintOk()
	}

	// add Node to path
	util.Print("Adding node to PATH... ")
	if err := os.Setenv("PATH", fmt.Sprintf("%s;%s", pathu.CurrentPath, os.Getenv("PATH"))); err != nil {
		util.PrintError(err)
	}
	util.PrintOk()

	// set NODE_PATH
	util.Print("Setting NODE_PATH... ")
	if err := os.Setenv("NODE_PATH", fs.FormatWinPath(path.Join(pathu.CurrentPath, "node_modules"))); err != nil {
		util.PrintError(err)
	}
	util.PrintOk()

	// create launch script
	util.Print("Creating launch script... ")
	launchScript := fs.Join(pathu.TmpPath, "launch.bat")
	if err := util.CreateFile(launchScript, app.GetLaunchScriptContent(pathu.CurrentPath)); err != nil {
		util.PrintError(err)
		return nil
	}
	util.PrintOk()

	if app.Conf.ImmediateMode == false {
		// wait user input to open the shell
		fmt.Print("\nPress Enter to open the shell...")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}

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
	log.Logger.Info("Starting up the shell... ")
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
