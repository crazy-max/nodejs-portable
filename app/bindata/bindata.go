package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _nodejs_portable_conf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8c\xc1\x8a\xc2\x30\x14\x45\xf7\xf9\x8a\xcb\x5b\xcd\xc0\x90\x30\x22\x22\xd9\x2a\xba\x12\xba\x17\x29\x8f\xf6\x69\x83\x4d\x23\x49\xac\x82\xf8\xef\x92\xd6\xe5\x3d\x87\x73\x5f\x0a\xa0\x51\x62\x72\x61\x20\x0b\x5a\xe8\xa5\xfe\xa7\xbf\x42\x9d\xf7\xd2\x3a\xce\x72\x08\xad\x90\xc5\x99\xfb\x24\x93\x79\x84\x78\xad\x38\x77\x25\xd0\xa6\xac\xb9\x68\xee\x29\x07\x5f\x4c\x22\x8b\xa3\x02\x00\xda\x58\x53\xc5\x70\x89\xec\xb1\x73\xbd\x24\xfc\x3c\xd7\xab\x5f\xb3\x77\xd9\x34\xbe\x9d\x42\x80\xb6\xd6\xf0\x10\x72\x27\xb1\xbe\x95\xe7\x2f\xd6\xda\x70\x1d\xa5\xe7\xec\x46\x99\x8d\x02\x4e\xea\xad\x3e\x01\x00\x00\xff\xff\x1e\xeb\x9c\xb2\xb9\x00\x00\x00")

func nodejs_portable_conf() ([]byte, error) {
	return bindata_read(
		_nodejs_portable_conf,
		"nodejs-portable.conf",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"nodejs-portable.conf": nodejs_portable_conf,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"nodejs-portable.conf": &_bintree_t{nodejs_portable_conf, map[string]*_bintree_t{
	}},
}}
