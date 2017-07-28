package nodejs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DIST_JSON = "https://nodejs.org/dist/index.json"
)

// Version structure
type Version struct {
	Version string `json:"version"`
}

// Versions structure
type Versions []Version

// GetLatestVersion returns the latest version of Node.js
func GetLatestVersion() (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(DIST_JSON)
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
