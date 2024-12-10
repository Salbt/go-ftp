package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Driver interface {
	Create() error
	Delete() error
	Rename(name string) error
	List()
}

type FileDriver struct {
}

type FileInfo struct {
}

func (driver *FileDriver) List(path string) error {
	// 使用 filepath.Walk 遍历目录
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径，并构造树状打印的前缀
		relativePath, _ := filepath.Rel(path, walkPath)
		depth := len(strings.Split(relativePath, string(os.PathSeparator)))

		// 计算前缀
		prefix := strings.Repeat("│   ", depth-1)
		if depth > 0 {
			if info.IsDir() {
				prefix += "├── " // 中间节点使用 ├──
			} else {
				prefix += "└── " // 叶子节点使用 └──
			}
		}

		// 打印当前路径
		fmt.Printf("%s%s\n", prefix, info.Name())
		return nil
	})
}
