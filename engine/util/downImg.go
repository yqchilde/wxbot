package util

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"github.com/yqchilde/pkgs/log"
)

type ImgInfo struct {
	Url  string
	Name string
}

// SingleDownload 单次下载
func SingleDownload(imgInfo ImgInfo) (fileName string, err error) {
	resp, err := http.Get(imgInfo.Url)
	if err != nil {
		return "", errors.Wrap(err, "SingleDownload http get error")
	}
	defer resp.Body.Close()
	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "SingleDownload read body error")
	}

	mimeType := mimetype.Detect(readBytes[:1024])
	fileName = imgInfo.Name + mimeType.Extension()
	os.WriteFile(fileName, readBytes, 0666)
	return fileName, nil
}

// BatchDownload 批量下载
func BatchDownload(imgInfoList []ImgInfo) error {
	var wg sync.WaitGroup
	wg.Add(len(imgInfoList))
	for _, img := range imgInfoList {
		go func(img ImgInfo) {
			defer wg.Done()
			res, err := http.Get(img.Url)
			if err != nil {
				log.Errorf("BatchDownload http get error: %v", err)
				return
			}
			defer res.Body.Close()
			readBytes, err := io.ReadAll(res.Body)
			if err != nil {
				log.Errorf("BatchDownload read body error: %v", err)
				return
			}

			os.WriteFile(img.Name, readBytes, 0666)
		}(img)
	}
	wg.Wait()

	return nil
}
