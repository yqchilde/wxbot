package download

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/yqchilde/pkgs/log"

	"wxBot/internal/model"
)

// SingleDownload 单次下载
func SingleDownload(imgInfo model.ImgInfo) (fileName string, err error) {
	resp, err := http.Get(imgInfo.Url)
	if err != nil {
		log.Errorf("SingleDownload http get error: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	readBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("SingleDownload read body error: %v", err)
	}

	mimeType := mimetype.Detect(readBytes[:1024])
	fileName = imgInfo.Name + mimeType.Extension()
	ioutil.WriteFile(fileName, readBytes, 0666)
	return fileName, nil
}

// BatchDownload 批量下载
func BatchDownload(imgInfoList []model.ImgInfo) error {
	var wg sync.WaitGroup
	wg.Add(len(imgInfoList))
	for _, img := range imgInfoList {
		go func(img model.ImgInfo) {
			defer wg.Done()
			res, err := http.Get(img.Url)
			if err != nil {
				log.Errorf("BatchDownload http get error: %v", err)
				return
			}
			defer res.Body.Close()
			readBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Errorf("BatchDownload read body error: %v", err)
				return
			}

			ioutil.WriteFile(img.Name, readBytes, 0666)
		}(img)
	}
	wg.Wait()

	return nil
}
