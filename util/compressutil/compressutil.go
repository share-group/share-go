package compressutil

import (
	"archive/tar"
	"fmt"
	"github.com/dsnet/compress/bzip2"
	"io"
	"os"
)

func Bzip2Compress(zipFileName string, files ...string) error {
	// 创建tar.bz2文件
	outputFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", zipFileName, err)
	}
	defer outputFile.Close()

	// 使用BZIP2压缩流来写入tar文件
	bzipWriter, err := bzip2.NewWriter(outputFile, &bzip2.WriterConfig{Level: bzip2.BestCompression})
	if err != nil {
		return fmt.Errorf("error creating BZIP2 writer: %w", err)
	}
	defer bzipWriter.Close()

	// 创建tar写入器
	tarWriter := tar.NewWriter(bzipWriter)
	defer tarWriter.Close()

	// 遍历所有文件，添加到tar中
	for _, fileName := range files {
		// 打开要添加的文件
		file, err2 := os.Open(fileName)
		if err2 != nil {
			return fmt.Errorf("error opening file %s: %w", fileName, err2)
		}
		defer file.Close()

		// 获取文件信息
		fileInfo, err3 := file.Stat()
		if err3 != nil {
			return fmt.Errorf("error getting file info for %s: %w", fileName, err3)
		}

		// 创建tar头
		header := &tar.Header{
			Name: fileInfo.Name(),
			Mode: int64(fileInfo.Mode()),
			Size: fileInfo.Size(),
		}

		// 写入tar头
		err4 := tarWriter.WriteHeader(header)
		if err4 != nil {
			return fmt.Errorf("error writing tar header for %s: %w", fileName, err4)
		}

		// 将文件内容写入tar包
		_, err5 := io.Copy(tarWriter, file)
		if err5 != nil {
			return fmt.Errorf("error writing file %s to tar: %w", fileName, err5)
		}
	}

	fmt.Printf("Successfully created %s\n", zipFileName)
	return nil
}
