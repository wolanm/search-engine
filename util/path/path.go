package path

import (
	"fmt"
	"runtime"
)

func getRepositoryDir() string {
	switch runtime.GOOS {
	case "linux":
		return "/var/lib/search_engine"
	case "darwin":
		return "/var/lib/search_engine"
	case "windows":
		return "C:\\ProgramData\\SearchEngine"
	default:
		return ""
	}
}

func GetInvertedDBPath() string {
	return fmt.Sprintf("%s/inverted.db", getRepositoryDir())
}

func GetTrieDBPath() string {
	return fmt.Sprintf("%s/trietree.db", getRepositoryDir())
}
