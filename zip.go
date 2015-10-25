package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func zipDir(file string, w io.Writer) error {
	zw := zip.NewWriter(w)
	err := filepath.Walk(file, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			w, err := zw.Create(path)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = zw.Close()
	if err != nil {
		return err
	}

	return nil
}
