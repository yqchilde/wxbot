package localimage

import (
	"image"
	"os"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var (
	imageSuffixes = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".webp"}
)

func WriteFile(filePath, content string) {
	err := os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		log.Errorf("向文件 %s 写入 %s 失败", filePath, content)
	}
}

func GetSubFolder(folderPath string) ([]os.DirEntry, error) {
	return ReadDir(folderPath, func(dirEntry os.DirEntry) bool {
		return dirEntry.IsDir()
	})
}

func GetFile(folderPath string) ([]os.DirEntry, error) {
	return ReadDir(folderPath, func(dirEntry os.DirEntry) bool {
		return !dirEntry.IsDir()
	})
}

func ReadDir(folderPath string, filters ...func(dirEntry os.DirEntry) bool) ([]os.DirEntry, error) {
	if !Exist(folderPath) {
		return nil, nil
	}

	// 获取给定文件夹下的所有子文件夹及文件
	dirEntries, err := os.ReadDir(folderPath)
	if err != nil {
		log.Errorf("读取目录 %s 失败, err: %v", folderPath, err)
		return nil, err
	}
	if filters == nil {
		return dirEntries, nil
	}

	filteredDirEntries := make([]os.DirEntry, 0)
	for _, entry := range dirEntries {
		for _, filter := range filters {
			if filter(entry) {
				filteredDirEntries = append(filteredDirEntries, entry)
			}
		}
	}

	return filteredDirEntries, nil
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

func IsImageFile(dirEntry os.DirEntry) bool {
	if dirEntry.IsDir() {
		return false
	}

	for _, suffix := range imageSuffixes {
		if strings.HasSuffix(strings.ToLower(dirEntry.Name()), suffix) {
			return true
		}
	}

	return false
}

func GeneratePdfFile(filePath string, imagePaths []string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	// 添加图片
	for _, imagePath := range imagePaths {
		file, err := os.Open(imagePath) // 打开图片文件
		if err != nil {
			return err
		}
		img, _, err := image.Decode(file) // 解码图片文件的配置，包括宽高信息
		_ = file.Close()                  // 关闭文件
		// 等比例缩放
		ratio := float64(img.Bounds().Size().X) / float64(img.Bounds().Size().Y)
		width := float64(595)
		size := &gopdf.Rect{W: width, H: width / ratio}

		pdf.AddPage()
		err = pdf.Image(imagePath, 0, 0, size)
		if err != nil {
			return err
		}
	}
	return pdf.WritePdf(filePath)
}
