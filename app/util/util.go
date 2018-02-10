package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/cavaliercoder/grab"
	"github.com/crazy-max/nodejs-portable/app/fs"
	"github.com/crazy-max/nodejs-portable/app/log"
	"github.com/fatih/color"
)

// Lib structure
type Lib struct {
	URL        string
	Dest       string
	OutputPath string
	Exe        string
}

// DownloadFile downloads a file and display status
func DownloadFile(filename string, url string) error {
	req, err := grab.NewRequest(url)
	if err != nil {
		return err
	}
	req.Filename = filename

	respch := grab.DefaultClient.DoAsync(req)
	resp := <-respch

	ticker := time.NewTicker(200 * time.Millisecond)
	for range ticker.C {
		if resp.IsComplete() {
			if resp.Error != nil {
				return resp.Error
			}
			break
		}
		fmt.Print(".")
	}

	ticker.Stop()
	return nil
}

// GetLatestVersion returns the latest version from github
func GetLatestVersion() (string, error) {
	probeURL := "https://raw.githubusercontent.com/crazy-max/nodejs-portable/master/res/latest"

	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(probeURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	}

	return "", fmt.Errorf("Status code %d", response.StatusCode)
}

// ReadLine reads input from console
func ReadLine(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// PrintPretty Pretty print of struct or slice
func PrintPretty(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// Print print info
func Print(str string) {
	log.Logger.Info(str)
	fmt.Print(str)
}

// Println print info
func Println(str string) {
	log.Logger.Info(str)
	fmt.Println(str)
}

// PrintError print error in red color
func PrintError(err error) {
	log.Logger.Error(err.Error())
	color.New(color.FgRed).Printf("Error: %s\n", err.Error())
}

// PrintErrorStr printed in red color
func PrintErrorStr(str string) {
	log.Logger.Error(str)
	color.New(color.FgRed).Printf("Error: %s\n", str)
}

// PrintWarning print warning in yellow color
func PrintWarning(err error) {
	log.Logger.Warning(err.Error())
	color.New(color.FgYellow).Printf("Warning: %s\n", err.Error())
}

// PrintWarningStr printed in yellow color
func PrintWarningStr(str string) {
	log.Logger.Warning(str)
	color.New(color.FgYellow).Printf("Warning: %s\n", str)
}

// PrintOk printed in green color
func PrintOk() {
	color.New(color.FgGreen).Print("OK!\n")
}

// DownloadLib download an external library
func DownloadLib(lib Lib) error {
	if lib.OutputPath != "" {
		if err := os.MkdirAll(lib.OutputPath, 777); err != nil {
			PrintError(err)
			return err
		}
		if _, err := fs.Stat(lib.Exe); err != nil {
			Print(fmt.Sprintf("Downloading %s...", lib.URL))
			if err := DownloadFile(lib.Dest, lib.URL); err != nil {
				fmt.Print(" ")
				PrintError(err)
				return err
			}
			fmt.Print(" ")
			PrintOk()

			Print(fmt.Sprintf("Unzipping %s... ", lib.Dest))
			if err := fs.Unzip(lib.Dest, lib.OutputPath); err != nil {
				PrintError(err)
				return err
			}
			PrintOk()

			Print(fmt.Sprintf("Seeking checkfile %s... ", lib.Exe))
			if _, err := fs.Stat(lib.Exe); err != nil {
				PrintError(err)
				return err
			}
			PrintOk()
		}
	}
	return nil
}

// QuitFatal quit the app and wait for user input
func QuitFatal(err error) {
	log.Logger.Error(err.Error())
	color.New(color.FgHiRed, color.Bold).Printf("\nFatal: %s\n", err.Error())
	fmt.Print("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	os.Exit(1)
}

// URLValid makes a GET request to a given URL and checks whether or not the
// resulting status code is 200.
func URLValid(url string) bool {
	resp, err := http.Get(url)
	if (err == nil) && (resp.StatusCode == 200) {
		return true
	}
	return false
}

// CreateFile creates / overwrites a file with content
func CreateFile(path string, content string) error {
	file, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

// SetConsoleTitle sets windows console title
func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		PrintError(err)
		return 0, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		PrintError(err)
		return 0, err
	}

	rTitle, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		PrintError(err)
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(rTitle)), 0, 0)
	return int(r), err
}
