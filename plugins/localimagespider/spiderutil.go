package localimagespider

import (
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/plugins/localimage"
)

const (
	UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

func DownloadImage(url, referer, folderPath string) {
	if !IsURL(url) || folderPath == "" {
		return
	}

	strs := strings.Split(url, "/")
	fileName := strs[len(strs)-1]
	bytes := req.C().SetBaseURL(url).
		SetCommonHeader("Referer", referer).
		SetCommonHeader("User-Agent", UA).
		Get().Do().Bytes()

	if !localimage.Exist(folderPath) && !localimage.MakeDir(folderPath) {
		return
	}
	filePath := filepath.Join(folderPath, fileName)
	err := os.WriteFile(filePath, bytes, os.ModePerm)
	if err != nil {
		log.Errorf("从 %s 下载图片并写入到 %s 失败, err: %v", url, filePath, err)
	}
	return
}

func GetImageLink(url string, selector string) []string {
	doc := ReadHtml(url)
	if doc == nil {
		return nil
	}

	links := make([]string, 0, 100)
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		links = append(links, src)
	})

	return links
}

func GetImageLinkInContainer(url, containerSelector string) []string {
	doc := ReadHtml(url)
	if doc == nil {
		return nil
	}

	links := make([]string, 0, 100)
	doc.Find(containerSelector).Each(func(i int, s *goquery.Selection) {
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			if src == "" {
				src, _ = s.Attr("data-loadsrc")
			}
			links = append(links, src)
		})
	})

	return links
}

func GetTextLinkInContainer(url, containerSelector string) map[string]string {
	doc := ReadHtml(url)
	if doc == nil {
		return nil
	}

	// 匹配带后缀的链接
	linkMap := make(map[string]string)
	doc.Find(containerSelector).Each(func(i int, s *goquery.Selection) {
		s.Find("a[href$='.html']").Each(func(i int, s *goquery.Selection) {
			if s.Text() != "" {
				// 获取 href
				href, _ := s.Attr("href")
				linkMap[href] = s.Text()
			}
		})
	})

	return linkMap
}

func GetTextLink(url string) map[string]string {
	doc := ReadHtml(url)
	if doc == nil {
		return nil
	}

	// 匹配带后缀的链接
	linkMap := make(map[string]string)
	doc.Find("a[href$='.html']").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "" {
			// 获取 href
			href, _ := s.Attr("href")
			linkMap[href] = s.Text()
		}
	})

	return linkMap
}

func ReadHtml(url string) *goquery.Document {
	if !IsURL(url) {
		return nil
	}
	// 爬取网页
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("访问 %s 失败, err: %v", url, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Errorf("访问 %s 失败, resp: %v", url, resp)
		return nil
	}

	// 使用 goquery 解析网页
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("解析 %s 内容失败, err: %v", url, err)
		return nil
	}
	return doc
}

func IsURL(url string) bool {
	if url == "" {
		return false
	}

	if _, err := neturl.Parse(url); err != nil {
		log.Errorf("%s 不是一个合法的 url", url)
		return false
	}

	return true
}
