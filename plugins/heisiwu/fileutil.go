package heisiwu

import (
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"os"
)

func WriteFile(filePath, content string) {
	err := os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		log.Errorf("向文件 %s 写入 %s 失败", filePath, content)
	}
}

func GetSubFolder(folderPath string) ([]os.DirEntry, error) {
	entries, err := ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	subFolders := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			subFolders = append(subFolders, entry)
		}
	}
	return subFolders, nil
}

func ReadDir(folderPath string) ([]os.DirEntry, error) {
	if !Exist(folderPath) {
		return nil, nil
	}

	// 获取给定文件夹下的所有子文件夹
	subFolders, err := os.ReadDir(folderPath)
	if err != nil {
		log.Errorf("读取目录 %s 失败, err: %v", folderPath, err)
	}
	return subFolders, err
}

func Exist(folderPath string) bool {
	// 判断该文件夹是否存在
	_, err := os.Stat(folderPath)
	return !os.IsNotExist(err)
}

func MakeDir(folderPath string) bool {
	// 创建该文件夹
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Errorf("创建目录 %s 失败, err: %v", folderPath, err)
	}
	return err == nil
}
