package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

var gopath string

func init() {
	readGOPath()
}

func readGOPath() {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	goabspath, err := filepath.Abs(gopath)
	if err != nil {
		panic(err)
	}
	gopath = goabspath
}

func importPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	dirpath := filepath.Dir(path)

	prefixPath := filepath.Join(gopath, "/src")

	if !strings.HasPrefix(dirpath, prefixPath) {
		return "", fmt.Errorf("The pkg path doesn't contain $GOPATH:%v", gopath)

	}
	pkgpath := strings.TrimPrefix(dirpath, prefixPath)
	pkgpath = strings.TrimPrefix(pkgpath, "/")
	return pkgpath, nil
}
