package plmm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/util"
)

func getPlmmPhoto(msg *openwechat.Message) {
	// 保证出图速度
	var isSend bool
	var plmmConf Plmm
	plugin.RawConfig.Unmarshal(&plmmConf)
	filepath.Walk(plmmConf.Dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			img, err := os.Open(path)
			if err != nil {
				log.Errorf("getPlmmPhoto open file error: %v", err)
				return err
			}
			defer img.Close()

			if _, err := msg.ReplyImage(img); err != nil {
				if strings.Contains(err.Error(), "operate too often") {
					msg.ReplyText("Warn: 被微信ban了，请稍后再试")
				} else {
					log.Errorf("msg.ReplyImage reply image error: %v", err)
				}
			}
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
		log.Errorf("getPlmmPhoto http get error: %v", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("getPlmmPhoto read body error: %v", err)
		return
	}

	var resp PlmmApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Errorf("getPlmmPhoto unmarshal error: %v", err)
		return
	}
	if resp.Code != 1 {
		log.Errorf("getPlmmPhoto api error: %v", resp.Msg)
		return
	}

	var imgInfo []util.ImgInfo
	for i, v := range resp.Data {
		imgInfo = append(imgInfo, util.ImgInfo{
			Url:  v.ImageUrl,
			Name: fmt.Sprintf("%s/%s.jpg", plmmConf.Dir, time.Now().Add(time.Duration(i)*time.Second).Format("20060102150405")),
		})
	}
	if err := util.BatchDownload(imgInfo); err != nil {
		log.Errorf("getPlmmPhoto batch download error: %v", err)
		return
	}

	if !isSend {
		filepath.Walk(plmmConf.Dir, func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				file, err := os.Open(path)
				if err != nil {
					log.Errorf("getPlmmPhoto open file error: %v", err)
					return err
				}
				defer file.Close()

				if _, err := msg.ReplyImage(file); err != nil {
					if strings.Contains(err.Error(), "operate too often") {
						msg.ReplyText("Warn: 被微信ban了，请稍后再试")
					} else {
						log.Errorf("msg.ReplyImage reply image error: %v", err)
					}
				}
				_ = os.Remove(path)
				return io.EOF
			}
			return nil
		})
	}
}
