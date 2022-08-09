package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"

	"wxBot/internal/config"
	"wxBot/internal/model"
	"wxBot/internal/pkg/download"
)

// GetShaoJiPhoto 获取烧鸡图片
func GetShaoJiPhoto(msg *openwechat.Message) {
	// 保证出图速度
	var isSend bool
	shaoJiConf := config.GetShaoJiConf()
	filepath.Walk(shaoJiConf.Dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			img, err := os.Open(path)
			if err != nil {
				log.Errorf("GetShaoJiPhoto open file error: %v", err)
				return err
			}
			defer img.Close()

			msg.ReplyImage(img)
			isSend = true
			_ = os.Remove(path)
			return io.EOF
		}
		return nil
	})

	// download 图片
	var imgInfo []model.ImgInfo
	for i := 0; i < 10; i++ {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, shaoJiConf.Url, nil)
		if err != nil {
			log.Errorf("GetShaoJiPhoto http new request error: %v", err)
			return
		}
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
		res, err := client.Do(req)
		if err != nil {
			log.Errorf("GetShaoJiPhoto client do error: %v", err)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Errorf("GetShaoJiPhoto read body error: %v", err)
			return
		}

		var resp model.ShaoJiApiResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Errorf("GetShaoJiPhoto unmarshal error: %v", err)
			return
		}
		if resp.Success != true {
			log.Errorf("GetShaoJiPhoto api error: %v", resp)
			return
		}

		imgInfo = append(imgInfo, model.ImgInfo{
			Url:  resp.Imgurl,
			Name: fmt.Sprintf("%s/%s.jpg", shaoJiConf.Dir, time.Now().Add(time.Duration(i)*time.Second).Format("20060102150405")),
		})
		res.Body.Close()
	}

	if err := download.BatchDownload(imgInfo); err != nil {
		log.Errorf("GetShaoJiPhoto batch download error: %v", err)
		return
	}

	if !isSend {
		filepath.Walk(shaoJiConf.Dir, func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				file, err := os.Open(path)
				if err != nil {
					log.Errorf("GetShaoJiPhoto open file error: %v", err)
					return err
				}
				defer file.Close()

				msg.ReplyImage(file)
				_ = os.Remove(path)
				return io.EOF
			}
			return nil
		})
	}
}
