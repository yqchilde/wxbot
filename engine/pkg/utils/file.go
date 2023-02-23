package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
)

// CheckFolderExists 检查文件夹是否存在，如果不存在则创建
func CheckFolderExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// CheckPathExists 判断文件/文件夹是否存在
func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Base64ToImage base64转图片
// b64Str base64字符串
// dst 图片保存路径
func Base64ToImage(b64Str, dst string) error {
	imgByte, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return err
	}
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(imgByte))
	if err != nil {
		return err
	}
	return nil
}
