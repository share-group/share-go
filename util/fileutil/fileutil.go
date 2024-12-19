package fileutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ListDir(dir string) []string {
	var files []string

	// 打开目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 处理错误
		if err != nil {
			return err
		}

		// 只添加文件，不添加目录
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Println(fmt.Sprintf("error reading directory %s: %v", dir, err))
		return make([]string, 0)
	}

	return files
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
