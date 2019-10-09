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

var _xretryLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x5d\x4f\xe3\x38\x14\x7d\xcf\xaf\x38\x4a\x1f\x4a\x45\x58\xb6\x7c\x3d\x20\x82\x84\x58\x84\xaa\x85\x2e\x6a\xd9\xd1\xa0\xc8\x0f\xa6\xb9\x6d\x2d\x12\xa7\x72\x9c\xe1\x63\x34\xf3\xdb\x47\x8e\xed\x7c\xb4\x83\x34\x0f\x50\x7f\x9c\x7b\xee\xbd\xe7\xdc\xf8\xe0\x20\x49\x02\x60\x46\x5a\xbd\xa3\xd4\x3c\xcb\x28\x45\x51\x69\xe4\x54\x96\x7c\x45\x25\x0a\x89\x52\x2b\xe2\x79\x79\x1e\x00\x01\x00\x4c\xe4\xa6\xd2\xe7\xf5\x12\xf8\xf7\xe6\x69\x9e\x8c\x19\x86\x2b\x55\x54\x9b\x21\x0e\x0f\xf1\xb8\x26\xd4\x3b\xe8\x02\x8a\x78\xda\xb0\x75\x63\x8e\x18\x86\x8b\x42\x96\x55\x4e\x6a\x92\x36\x81\xfe\x08\x22\x85\x5e\x73\x0d\x51\xd6\x1c\xef\x86\x6c\x91\x71\x91\xff\x96\xed\x98\x61\xa8\x4c\x17\x77\x22\x17\xba\x66\xbb\xe7\x6f\x22\xaf\x72\xc8\x2a\x7f\x26\x85\x62\x09\x03\x10\xfd\xb8\x93\xba\x8a\x4a\x7e\x16\xd2\xe8\xa0\x0b\x3c\x93\x2d\x80\xd2\x2e\xc3\x29\xc3\x70\xc3\x57\x34\x17\x1f\x54\x93\x3c\xf0\x95\x90\x5c\x0b\xa3\x9c\xf8\x20\xbc\xae\x49\x42\x68\x52\x5c\x0b\xb9\x42\xf1\x8d\x14\x36\x24\x53\xb3\x69\xe8\x33\x51\xea\x2e\xeb\x19\xc3\x30\x25\x9e\x66\x42\x5a\xd6\x7b\x21\xeb\xd2\xb4\xc8\x09\x7b\x79\x39\xb2\xea\x70\x4f\x81\xbc\x2a\xb5\x29\x51\xc8\x86\xdd\x90\xfa\xc2\x0b\x59\x8a\x94\x14\xa5\x5d\x9b\x83\x3a\xe5\xd5\xec\xf6\x4b\x32\x65\x26\xcb\xdc\x5a\x8d\x17\x7a\x0f\x18\x0b\xb2\x62\xc1\x33\xa8\x05\xcf\x32\xc4\x50\x94\x8a\xf2\x2f\xb3\x71\x17\xba\xb0\x83\x13\xe3\xfb\x0f\x77\xe4\x04\xea\x1e\xb5\xb6\x20\x86\x2e\xac\xb6\x7b\xce\xb5\x91\x8f\x33\x1e\xec\xdc\x9f\x34\xf7\x5e\xe1\x1d\xc8\x69\x03\xf1\x72\xed\x40\xce\xd8\x28\x08\x96\x85\x82\x40\x8c\x71\x84\x81\x69\x18\x69\x11\x00\x8e\xdc\xea\x35\xaf\xf2\x9c\xd7\xed\xd4\x1d\xef\x85\x5f\x1f\x6e\xa6\xff\x4c\xa6\xb7\x61\x64\x35\x12\x2c\xf2\xf3\x3e\x0a\x00\xb1\xdc\x8a\x34\xdf\xc1\x25\xfe\x86\x5e\x93\xac\x95\xed\xd1\x4f\x64\x4a\x6f\x88\xb7\x63\x8e\x58\x0d\x7d\x5d\x8b\x8c\xfa\xd0\x9f\x31\xc2\x10\x5c\xa6\x18\x78\xa5\x2f\x9c\x50\x75\xf1\x3b\x19\xfe\xa8\xf2\xa8\x97\x24\x42\xb8\x1f\x46\x8d\xbc\x23\x47\xbb\x55\x72\x18\xba\x73\xa3\xa2\x76\x2a\xfa\xa4\x4d\x2d\x5d\x45\x12\xcd\xcc\xb7\x75\xd1\x75\xdf\x74\xd2\xb9\x3e\x66\xb8\x8c\x5b\xd7\x1a\xd1\x7a\x15\x18\xe0\xa0\x5d\xef\x8f\x19\x62\xdf\x52\x07\xed\xf4\x49\xbc\x50\x16\xd7\xc6\x35\x50\x92\x69\xb0\x95\x62\xcb\x17\x93\x71\xec\x03\x5a\xb8\x58\xb6\x1d\x5f\xb4\xe3\xd8\xab\xfa\x33\xd5\x3c\x8b\xfd\xb5\xff\x0d\x9f\x77\xf5\x32\x76\xb6\x36\x6c\xcf\x8a\xf8\x8b\xc3\x9a\xbf\xde\xf8\xfa\xb0\xce\x04\xdb\x87\xb1\xb5\xff\xfa\xee\x6a\x72\x1f\x46\x8d\x2e\x82\x25\xa7\xdd\x11\x70\x2f\x70\xd4\xc8\xdf\x83\x1a\x44\x38\xbb\x79\x9c\x3d\x5d\xff\xf7\xff\xf4\xb1\xcf\x73\xc2\xf6\xc7\x6e\xfa\x07\x36\x6d\x6f\xe4\x3b\xc8\x41\xbb\xb6\x7e\xd4\xf0\x64\xcc\xfc\xcc\xbb\xc7\x22\x19\xb8\x85\x45\xb5\x51\x8d\x58\x56\x03\x45\xba\x52\xd2\x47\xfd\x0a\x00\x00\xff\xff\x7e\x4c\xe1\x53\xbb\x06\x00\x00")

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

	info := bindataFileInfo{name: "xretry.lua", size: 1723, mode: os.FileMode(436), modTime: time.Unix(1570354580, 0)}
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
