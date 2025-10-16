package log

import (
	"errors"
	"fmt"
	"runtime"
)

// RootLogDir 日志根目录
func RootLogDir() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "/var/log/search_engine", nil
	case "darwin":
		return "/var/log/search_engine", nil
	case "windows":
		return "C:\\ProgramData\\SearchEngine", nil
	default:
		return "", errors.New("unsupported operating system: " + runtime.GOOS)
	}
}

// LogDir 根据服务名获取对应的日志目录
func LogDir(serviceName string) string {
	rootLogDir, err := RootLogDir()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/", rootLogDir, serviceName)
}
