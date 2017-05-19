package file

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
)

//WriteLines 按行写字符串
func WriteLines(filePath string, lines []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(lines) != 0 {
		for _, line := range lines {
			_, err := f.WriteString(strings.TrimSpace(line) + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//ReadAll 读取文件内容
func ReadAll(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

//ReadLines 按行读取文件
func ReadLines(filePath string) []string {
	all, _ := ReadAll(filePath)
	lines := strings.Split(string(all), "\n")
	return lines
}

//CopyFile 复制文件
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		glog.Error(err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		glog.Error(err)
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//IsExist 判断文件或文件夹是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
	// 或者
	//return err == nil || !os.IsNotExist(err)
	// 或者
	//return !os.IsNotExist(err)
}
