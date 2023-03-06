package heisiwu

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

const (
	HeisiwuURL          = "http://hs.heisiwu.com/"
	CategoryURLTemplate = HeisiwuURL + "%s/page/%v"
	PageInfoFile        = "pageinfo"
)

var (
	cats = []string{"heisi", "baisi", "juru", "jk", "mcn", "meizu"}
	job  = gocron.NewScheduler(time.Local)
)

func start(storageFolder string) {
	// 每个分类爬虫任务的执行状态，true 运行，false 未运行
	categorySpiderTaskState := make(map[string]bool)
	job.Tag("黑丝屋爬虫任务").Every(5).Minutes().Do(func() {
		category := cats[rand.Intn(len(cats))]
		if categorySpiderTaskState[category] {
			return
		}
		categorySpiderTaskState[category] = true
		crawlCategory(storageFolder, category)
		categorySpiderTaskState[category] = false
	})
	job.StartAsync()
}

func crawlCategory(storageFolder, category string) {
	// 文件夹路径
	folderPath := filepath.Join(storageFolder, category)
	if !Exist(folderPath) && !MakeDir(folderPath) {
		return
	}

	currentPageNum, err := getCurrentPageNum(folderPath)
	if err != nil {
		return
	}
	currentPageNum += 1

	url := fmt.Sprintf(CategoryURLTemplate, category, currentPageNum)
	linkMap := GetTextLink(url)
	if len(linkMap) == 0 {
		// 为空说明当前分类爬完了，重置页码，这样新增数据之后，可以重新爬到
		WriteFile(filepath.Join(folderPath, PageInfoFile), "0")
		return
	}

	dirEntries, err := GetSubFolder(folderPath)
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

		topicFolderPath := filepath.Join(folderPath, topicId+"-"+title)
		if !MakeDir(topicFolderPath) {
			return
		}
		crawlTopic(link, topicFolderPath)
	}

	WriteFile(filepath.Join(folderPath, PageInfoFile), strconv.Itoa(currentPageNum))
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
	imageLinks := GetImageLink(link, "img[loading=\"lazy\"]")
	for _, imageLink := range imageLinks {
		DownloadImage(imageLink, link, topicFolderPath)
	}
}

func getCurrentPageNum(folderPath string) (int, error) {
	pageInfoFilePath := filepath.Join(folderPath, PageInfoFile)
	if !Exist(pageInfoFilePath) {
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
