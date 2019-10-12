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

var _xretryLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\x5b\x4f\xf3\x38\x10\x7d\xcf\xaf\x38\x4a\x1f\x4a\x45\x58\xb6\xdc\x1e\x10\x41\x42\x2c\x42\x15\xd0\x45\x2d\xbb\x5a\x14\xf9\xc1\x34\xd3\xd6\x22\x71\x2a\xc7\x59\x2e\xab\xe5\xb7\xaf\x1c\xdb\xb9\xb4\xcb\xa7\xef\x01\xea\xd8\x67\xce\xcc\x9c\x33\xf6\xc1\x41\x92\x04\xc0\x8c\xb4\xfa\x40\xa9\x79\x96\x51\x8a\xa2\xd2\xc8\xa9\x2c\xf9\x8a\x4a\x14\x12\xa5\x56\xc4\xf3\xf2\x3c\x00\x02\x00\x98\xc8\x4d\xa5\xcf\xeb\x25\x70\x77\xf3\x3c\x4f\xc6\x0c\xc3\x95\x2a\xaa\xcd\x10\x87\x87\x78\x5a\x13\xea\x2f\xe8\x02\x8a\x78\xda\xb0\x75\x63\x8e\x18\x86\x8b\x42\x96\x55\x4e\x6a\x92\x36\x81\x7e\x0b\x22\x85\x5e\x73\x0d\x51\xd6\x1c\x1f\x86\x6c\x91\x71\x91\xff\x2f\xdb\x31\xc3\x50\x99\x2e\xee\x45\x2e\x74\xcd\xf6\xc0\xdf\x45\x5e\xe5\x90\x55\xfe\x42\x0a\xc5\x12\x06\x20\xfa\x71\x27\x75\x15\x95\xfc\x2e\xa4\xd1\x41\x17\x78\x21\x5b\x00\xa5\x5d\x86\x53\x86\xe1\x86\xaf\x68\x2e\x3e\xa9\x26\x79\xe4\x2b\x21\xb9\x16\x46\x39\xf1\x49\x78\x5b\x93\x84\xd0\xa4\xb8\x16\x72\x85\xe2\x6f\x52\xd8\x90\x4c\xcd\x47\x43\x9f\x89\x52\x77\x59\xcf\x18\x86\x29\xf1\x34\x13\xd2\xb2\x3e\x08\x59\x97\xa6\x45\x4e\xd8\xcb\xcb\x91\x55\x87\x7b\x0a\xe4\x55\xa9\x4d\x89\x42\x36\xec\x86\xd4\x17\x5e\xc8\x52\xa4\xa4\x28\xed\xda\x1c\xd4\x29\xaf\x66\xb7\x7f\x26\x53\x66\xb2\xcc\xad\xd5\x78\xa5\x8f\x80\xb1\x20\x2b\x16\x3c\x83\x5a\xf0\x2c\x43\x0c\x45\xa9\x28\x7f\x31\x1f\xee\x40\x17\x76\x70\x62\xfc\xf3\xaf\xdb\x72\x02\x75\xb7\x5a\x5b\x10\x43\x17\x56\xdb\x3d\xe7\xda\xc8\xc7\x19\x0f\x76\xce\x4f\x9a\x73\xaf\xf0\x0e\xe4\xb4\x81\x78\xb9\x76\x20\x67\x6c\x14\x04\xcb\x42\x41\x20\xc6\x38\xc2\xc0\x34\x8c\xb4\x08\x00\x47\x6e\xf5\x9a\x57\x79\xce\xeb\x76\xea\x8e\xf7\xc2\xbf\x1e\x6f\xa6\xbf\x4d\xa6\xb7\x61\x64\x35\x12\x2c\xf2\xf3\x3e\x0a\x00\xb1\xdc\x8a\x34\xf7\xe0\x12\xbf\x42\xaf\x49\xd6\xca\xf6\xe8\x27\x32\xa5\x77\xc4\xdb\x31\x47\xac\x86\xbe\xad\x45\x46\x7d\xe8\x57\x8c\x30\x04\x97\x29\x06\x5e\xe9\x0b\x27\x54\x5d\xfc\x4e\x86\x9f\xaa\x3c\xea\x25\x89\x10\xee\x87\x51\x23\xef\xc8\xd1\x6e\x95\x1c\x86\x6e\xdf\xa8\xa8\x9d\x8a\x3e\x69\x53\x4b\x57\x91\x44\x33\x73\xb7\x2e\xba\xee\x9b\x4e\x3a\xc7\xc7\x0c\x97\x71\xeb\x5a\x23\x5a\xaf\x02\x03\x1c\xb4\xeb\xfd\x31\x43\xec\x5b\xea\xa0\x9d\x3e\x89\x17\xca\xe2\xda\xb8\x06\x4a\x32\x0d\xb6\x52\x6c\xf9\x62\x32\x8e\x7d\x40\x0b\x17\xcb\xb6\xe3\x8b\x76\x1c\x7b\x55\x7f\xa7\x9a\x67\xb1\xbf\xf6\xbf\xe1\xf3\xae\x5e\xc6\xce\xd6\x86\xed\x45\x11\x7f\x75\x58\xf3\xd7\x1b\x5f\x1f\xd6\x99\x60\xfb\x30\xb6\xf6\x5f\xdf\x5f\x4d\x1e\xc2\xa8\xd1\x45\xb0\xe4\xb4\x3b\x02\xee\x05\x8e\x1a\xf9\x7b\x50\x83\x08\x67\x37\x4f\xb3\xe7\xeb\xdf\xff\x98\x3e\xf5\x79\x4e\xd8\xfe\xd8\x4d\xff\xc0\xa6\xed\x8d\xbc\x58\xda\x62\xcc\x5d\xf8\x8a\xb1\xe4\x59\xd9\x53\xa9\xc3\x34\x68\xd7\xd6\x2f\x1f\xe8\xef\x04\xfc\x73\x92\x0c\xdc\xc2\xe2\xda\x38\x2b\x6b\x56\x92\xc3\xfb\xfe\xaf\xae\xef\x7e\xd0\x7d\xaf\xd5\xd1\x96\x35\xb5\xdc\x8a\x74\xa5\xa4\x4f\xff\x5f\x00\x00\x00\xff\xff\xdf\xfb\xdd\xfb\x26\x07\x00\x00")

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

	info := bindataFileInfo{name: "xretry.lua", size: 1830, mode: os.FileMode(436), modTime: time.Unix(1570900571, 0)}
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
