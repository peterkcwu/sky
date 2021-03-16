package config

import (
	"os"
	"path/filepath"
	"strings"
)

var ApiConf *Config
var LaunchDir string

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return "/tmp/"
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func init() {
	LaunchDir = GetCurrentDirectory()
}

