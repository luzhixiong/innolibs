package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压
func Unzip(zipFile string, destDir string) (paths []string, files []string, err error) {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if strings.Index(f.Name, "..") > -1 {
			err = fmt.Errorf("%s 文件名不合法", f.Name)
			return
		}
		fpath := filepath.Join(destDir, f.Name)
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return
			}

			inFile, err := f.Open()
			if err != nil {
				return nil, nil, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return nil, nil, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return nil, nil, err
			}
			files = append(files, fpath)
		}
	}
	return
}
