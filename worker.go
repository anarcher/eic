package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Worker struct {
	DryRun bool
}

func (w *Worker) WorkDir(path string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if filepath.Ext(filePath) == ".go" {
			if err := w.WorkFile(filePath); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (w *Worker) WorkFile(path string) error {
	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	importPath, err := importPath(path)
	if err != nil {
		return err
	}

	astFile := NewASTFile(src, importPath)
	if err := astFile.EnsureImportComment(); err != nil {
		return err
	}

	if w.DryRun == true {
		fmt.Println(astFile.String())
		return nil
	}

	if !astFile.IsChanged() {
		return nil
	}

	if err := ioutil.WriteFile(path, astFile.Bytes(), finfo.Mode()); err != nil {
		return err
	}
	fmt.Printf("Processed file: %v\n", path)

	return nil
}
