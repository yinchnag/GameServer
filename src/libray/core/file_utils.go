package core

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 获取当前执行文件绝对路径
func GetExecutableAbsPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	res = strings.Replace(res, "\\", "/", -1)
	return res
}

// 获取GOMOD路径
func getModPath() string {
	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	path := string(bytes.TrimSpace(stdout))
	path = strings.Replace(path, "\\", "/", -1)
	return path
}

// 获取基础路径
func GetBasePath() string {
	var dir string
	if runtime.GOOS != "windows" {
		dir, _ = os.Getwd()
	}
	if dir == "NUL" || dir == "" {
		dir = getModPath()
	}
	if dir == "NUL" || dir == "" {
		dir = GetExecutableAbsPath()
	}
	if !strings.Contains(dir, "/src") {
		return dir
	}
	arr := strings.Split(dir, "/src")
	return arr[0] + "/bin"
}

// 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
