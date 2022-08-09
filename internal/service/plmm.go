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

// GetPlmmPhoto 获取Plmm图片
func GetPlmmPhoto(msg *openwechat.Message) {
	// 保证出图速度
	var isSend bool
	plmmConf := config.GetPlmmConf()
	filepath.Walk(plmmConf.Dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			img, err := os.Open(path)
			if err != nil {
				log.Errorf("GetPlmmPhoto open file error: %v", err)
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
	apiUrl := fmt.Sprintf("%s?app_id=%s&app_secret=%s", plmmConf.Url, plmmConf.AppId, plmmConf.AppSecret)
	res, err := http.Get(apiUrl)
	if err != nil {
		log.Errorf("GetPlmmPhoto http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("GetPlmmPhoto read body error: %v", err)
		return
	}

	var resp model.PlmmApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("GetPlmmPhoto unmarshal error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("GetPlmmPhoto api error: %v", resp.Msg)
		return
	}

	var imgInfo []model.ImgInfo
	for i, v := range resp.Data {
		imgInfo = append(imgInfo, model.ImgInfo{
			Url:  v.ImageUrl,
			Name: fmt.Sprintf("%s/%s.jpg", plmmConf.Dir, time.Now().Add(time.Duration(i)*time.Second).Format("20060102150405")),
		})
	}
	if err := download.BatchDownload(imgInfo); err != nil {
		log.Errorf("GetPlmmPhoto batch download error: %v", err)
		return
	}

	if !isSend {
		filepath.Walk(plmmConf.Dir, func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				file, err := os.Open(path)
				if err != nil {
					log.Errorf("GetPlmmPhoto open file error: %v", err)
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
