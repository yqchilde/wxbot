package heisiwu

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
)

const (
	UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

func WriteFile(filePath, content string) {
	err := os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		log.Errorf("向文件 %s 写入 %s 失败", filePath, content)
	}
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
	err := os.Mkdir(folderPath, os.ModePerm)
	if err != nil {
		log.Errorf("创建目录 %s 失败, err: %v", folderPath, err)
	}
	return err == nil
}

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

	if !Exist(folderPath) && !MakeDir(folderPath) {
		return
	}
	filePath := folderPath + string(os.PathSeparator) + fileName
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

func GetPath(strs ...string) string {
	return strings.Join(strs, string(os.PathSeparator))
}
