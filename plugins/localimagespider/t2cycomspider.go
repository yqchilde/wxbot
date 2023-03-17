package localimagespider

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/plugins/localimage"
)

const (
	PageInfoFile         = "pageinfo"
	T2cyComUrl           = "https://t2cy.com"
	CosplayUri           = "/acg/cos"
	CosplayPagerTemplate = "/index_%v.html"
)

func crawlCosplay(storageFolder string) {
	// 文件夹路径
	folderPath := filepath.Join(storageFolder, "cosplay")
	if !localimage.Exist(folderPath) && !localimage.MakeDir(folderPath) {
		return
	}

	currentPageNum, err := getCurrentPageNum(folderPath)
	if err != nil {
		return
	}
	currentPageNum += 1

	url := T2cyComUrl + CosplayUri
	if currentPageNum > 1 {
		url += fmt.Sprintf(CosplayPagerTemplate, currentPageNum)
	}
	linkMap := GetTextLinkInContainer(url, "ul[class=\"cy2-coslist clr\"]")
	if len(linkMap) == 0 {
		// 为空说明当前分类爬完了，重置页码，这样新增数据之后，可以重新爬到
		localimage.WriteFile(filepath.Join(folderPath, PageInfoFile), "0")
		return
	}

	dirEntries, err := localimage.GetSubFolder(folderPath)
	if err != nil {
		return
	}
	existTopicMap := convert(dirEntries)

	for link, title := range linkMap {
		strs := strings.Split(link, "/")
		topicId := strings.TrimSuffix(strs[len(strs)-1], ".html")
		// TODO 这样不严谨，进一步应该比较一下元素数量
		if existTopicMap[topicId] != "" {
			continue
		}

		replaces := []string{"*", "_", ".", "_", "\"", "_", "/", "_", "\"", "_", "[", "_", "]", "_", ":", "_", ";", "_", "|", "_", ",", "_"}
		title = strings.NewReplacer(replaces...).Replace(title)
		topicFolderPath := filepath.Join(folderPath, topicId+"-"+title)
		if !localimage.MakeDir(topicFolderPath) {
			return
		}
		crawlTopic(T2cyComUrl+link, topicFolderPath)
	}

	localimage.WriteFile(filepath.Join(folderPath, PageInfoFile), strconv.Itoa(currentPageNum))
}

func convert(dirEntries []os.DirEntry) map[string]string {
	existTopicMap := make(map[string]string)
	for _, entry := range dirEntries {
		topicInfo := strings.Split(entry.Name(), "-")
		if len(topicInfo) != 2 {
			continue
		}

		existTopicMap[topicInfo[0]] = topicInfo[1]
	}
	return existTopicMap
}

func crawlTopic(link, topicFolderPath string) {
	imageLinks := GetImageLinkInContainer(link, "div[class=\"w maxImg tc\"]")
	for _, imageLink := range imageLinks {
		DownloadImage(T2cyComUrl+imageLink, link, topicFolderPath)
	}
}

func getCurrentPageNum(folderPath string) (int, error) {
	pageInfoFilePath := filepath.Join(folderPath, PageInfoFile)
	if !localimage.Exist(pageInfoFilePath) {
		return 0, nil
	}

	bytes, err := os.ReadFile(pageInfoFilePath)
	if err != nil {
		log.Errorf("读取文件 %s 失败, err: %v", pageInfoFilePath, err)
		return 0, err
	}

	pageNum, err := strconv.Atoi(string(bytes))
	if err != nil {
		log.Errorf("文件 %s 内容非数字, err: %v", pageInfoFilePath, err)
		return 0, nil
	}

	return pageNum, nil
}
