package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func zipFiles(paths []string, w io.Writer) error {
	zw := zip.NewWriter(w)
	for _, p := range paths {
		w, err := zw.Create(p)
		if err != nil {
			return err
		}
		file, err := os.Open(p)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
		if err != nil {
			return err
		}
	}
	return zw.Close()
}

func zipDir(file string, w io.Writer) error {
	//Traverse directory for all files
	var paths []string
	err := filepath.Walk(file, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	//Zip all files
	err = zipFiles(paths, w)

	if err != nil {
		return err
	}

	return nil
}
