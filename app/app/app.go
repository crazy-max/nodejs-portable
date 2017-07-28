package app

const (
	ID      = "nodejs-portable"
	NAME    = "Node.js Portable"
	VERSION = "2.0.0"
	PACKAGE = "github.com/crazy-max/nodejs-portable"
	URL     = "https://" + PACKAGE
)

// Lib structure
type Lib struct {
	Url        string `json:"url"`
	Dest       string
	OutputPath string
	Exe        string
}
