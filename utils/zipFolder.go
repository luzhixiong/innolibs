package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipFolder 递归打包整个文件夹到 ZIP 文件
func ZipFolder(source, targetZip string) error {
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径，确保 ZIP 内部路径不包含绝对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// 如果是目录，直接返回
		if info.IsDir() {
			return nil
		}

		// 创建 ZIP 文件条目
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// 打开原始文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 复制文件内容到 ZIP 文件
		_, err = io.Copy(zipEntry, srcFile)
		return err
	})
}
