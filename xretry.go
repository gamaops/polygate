// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// xretry.lua
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _xretryLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x5d\x4f\xe3\x38\x14\x7d\xcf\xaf\x38\x4a\x1f\x4a\x45\x58\xb6\x7c\x3d\x20\x82\x84\x58\x84\xaa\x85\x2e\x6a\xd9\xd5\xa2\xc8\x0f\xa6\xb9\x6d\x2d\x12\xa7\x72\x9c\xe1\x63\x34\xfc\xf6\x91\x63\x3b\x1f\xed\x20\xcd\x03\xd4\x1f\xe7\x9c\x7b\xef\xb9\xd7\x39\x38\x48\x92\x00\x98\x91\x56\xef\x28\x35\xcf\x32\x4a\x51\x54\x1a\x39\x95\x25\x5f\x51\x89\x42\xa2\xd4\x8a\x78\x5e\x9e\x07\x40\x00\x00\x13\xb9\xa9\xf4\x79\xbd\x04\xfe\xbe\x79\x9a\x27\x63\x86\xe1\x4a\x15\xd5\x66\x88\xc3\x43\x3c\xae\x09\xf5\x0e\xba\x80\x22\x9e\x36\x6a\x5d\xce\x11\xc3\x70\x51\xc8\xb2\xca\x49\x4d\xd2\x86\xe8\x8f\x20\x52\xe8\x35\xd7\x10\x65\xad\xf1\x6e\xc4\x16\x19\x17\xf9\x2f\xd5\x8e\x19\x86\xca\x54\x71\x27\x72\xa1\x6b\xb5\x7b\xfe\x26\xf2\x2a\x87\xac\xf2\x67\x52\x28\x96\x30\x00\xd1\xe7\x9d\xd4\x59\x54\xf2\x2b\x4a\xe3\x83\x2e\xf0\x4c\x36\x01\x4a\xbb\x0a\xa7\x0c\xc3\x0d\x5f\xd1\x5c\x7c\x50\x2d\xf2\xc0\x57\x42\x72\x2d\x8c\x73\xe2\x83\xf0\xba\x26\x09\xa1\x49\x71\x2d\xe4\x0a\xc5\x37\x52\xd8\x90\x4c\xcd\xa6\x91\xcf\x44\xa9\xbb\xaa\x67\x0c\xc3\x94\x78\x9a\x09\x69\x55\xef\x85\xac\x53\xd3\x22\x27\xec\xe5\xe5\xc8\xba\xc3\xbd\x04\xf2\xaa\xd4\x26\x45\x21\x1b\x75\x23\xea\x13\x2f\x64\x29\x52\x52\x94\x76\xdb\x1c\xd4\x21\xaf\x66\xb7\xff\x25\x53\x66\xa2\xcc\x6d\xab\xf1\x42\xef\x01\x63\x41\x56\x2c\x78\x06\xb5\xe0\x59\x86\x18\x8a\x52\x51\xfe\x61\x36\xee\x42\x17\x76\x70\x62\x7c\xff\xe1\x8e\x9c\x41\xdd\xa3\xb6\x2d\x88\xa1\x0b\xeb\xed\x9e\xeb\xda\xc8\xf3\x4c\x0f\x76\xee\x4f\x9a\x7b\xef\xf0\x0e\xe4\xb4\x81\x78\xbb\x76\x20\x67\x6c\x14\x04\xcb\x42\x41\x20\xc6\x38\xc2\xc0\x14\x8c\xb4\x08\x00\x27\x6e\xfd\x9a\x57\x79\xce\xeb\x72\xea\x8a\xf7\xc2\xff\x1f\x6e\xa6\x7f\x4d\xa6\xb7\x61\x64\x3d\x12\x2c\xf2\xf3\x3e\x0a\x00\xb1\xdc\x62\x9a\x77\x70\x89\x3f\xa1\xd7\x24\x6b\x67\x7b\xf2\x13\x99\xd2\x1b\xe2\x6d\xce\x11\xab\xa1\xaf\x6b\x91\x51\x1f\xfa\x19\x23\x0c\xc1\x65\x8a\x81\x77\xfa\xc2\x19\x55\x27\xbf\x13\xe1\xb7\x32\x8f\x7a\x41\x22\x84\xfb\x61\xd4\xd8\x3b\x72\xb2\x5b\x29\x87\xa1\x3b\x37\x2e\x6a\xe7\xa2\x0f\xda\xe4\xd2\x75\x24\xd1\xcc\xbc\xad\x8b\x6e\xf7\x4d\x25\x9d\xeb\x63\x86\xcb\xb8\xed\x5a\x63\x5a\x2f\x03\x03\x1c\xb4\xeb\xfd\x31\x43\xec\x4b\xea\xa0\x9d\x3f\x89\x37\xca\xe2\x5a\x5e\x03\x25\x99\x06\x5b\x21\xb6\xfa\x62\x22\x8e\x3d\xa1\x85\x8b\x65\x5b\xf1\x45\x3b\x8e\xbd\xac\xbf\x72\xcd\xab\xd8\x5f\xfb\xdf\xe8\xf9\xae\x5e\xc6\xae\xad\x8d\xda\xb3\x22\xfe\xe2\xb0\xe6\xaf\x37\xbe\x9e\xd6\x99\x60\xfb\x61\x6c\xdb\x7f\x7d\x77\x35\xb9\x0f\xa3\xc6\x17\xc1\x92\xd3\xee\x08\xb8\x2f\x70\xd4\xd8\xdf\x83\x1a\x44\x38\xbb\x79\x9c\x3d\x5d\xff\xf3\xef\xf4\xb1\xaf\x73\xc2\xf6\xc7\x6e\xfa\x07\x36\xac\x19\x79\xd3\xdb\x7a\x67\xde\xc0\x67\x8c\x25\xcf\xca\x8e\x3b\x1d\xfe\xa0\x5d\xdb\x2e\x79\x9a\x7f\x09\xee\x13\x92\x0c\xdc\xc2\xa2\x5a\x56\xd7\x16\x45\xba\x52\xd2\x53\x7e\x06\x00\x00\xff\xff\x0a\x58\x96\xb6\xce\x06\x00\x00")

func xretryLuaBytes() ([]byte, error) {
	return bindataRead(
		_xretryLua,
		"xretry.lua",
	)
}

func xretryLua() (*asset, error) {
	bytes, err := xretryLuaBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "xretry.lua", size: 1742, mode: os.FileMode(436), modTime: time.Unix(1570883249, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"xretry.lua": xretryLua,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"xretry.lua": &bintree{xretryLua, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
