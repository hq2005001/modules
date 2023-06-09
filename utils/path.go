package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetAbsDir 得到程序运行的绝对路径
func GetAbsDir() string {
	workingDir, _ := os.Getwd()
	binPath, err := filepath.Abs(workingDir)
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Dir(binPath)
}

// RuntimeDir 运行时目录
func RuntimeDir(path ...string) string {
	dirArr := make([]string, 0)
	dirArr = append(dirArr, GetAbsDir())
	dirArr = append(dirArr, "runtime")
	if len(path) > 0 {
		dirArr = append(dirArr, strings.Join(path, string(os.PathSeparator)))
	}
	dir := strings.Join(dirArr, string(os.PathSeparator))
	if !IsDir(dir) {
		log.Println(dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return dir
}

// IsDir 是否是文件夹
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Mkdir(path []string) (err error) {
	dirArr := make([]string, 0)
	dirArr = append(dirArr, GetAbsDir())
	dirArr = append(dirArr, "runtime")
	if len(path) > 0 {
		dirArr = append(dirArr, strings.Join(path, string(os.PathSeparator)))
	}
	dir := strings.Join(dirArr, string(os.PathSeparator))
	if !IsDir(dir) {
		log.Println(dir)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return err
}
