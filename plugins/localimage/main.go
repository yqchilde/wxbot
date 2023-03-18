package localimage

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"modernc.org/mathutil"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	storageFolder = ""
	pdfSuffix     = ".pdf"
)

func init() {
	engine := control.Register("localimage", &control.Options{
		Alias: "读取本地图片",
		Help: "指令:\n" +
			"* 列出图片目录\n" +
			"* 来点[目录名称]图片\n" +
			"* 来[数量]张[目录名称]图片\n" +
			"* 来点[搜索词]图片\n" +
			"* 来[数量]张[搜索词]图片\n" +
			"* [目录名称]\n" +
			"* [目录名称] [数量]\n",
		DataFolder: "localimage",
	})

	storageFolder = engine.GetCacheFolder()

	engine.OnFullMatch("列出图片目录").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		folders, _ := GetSubFolder(storageFolder)
		folderInfos := dir(folders)
		if len(folderInfos) == 0 {
			ctx.ReplyTextAndAt("目录内容为空，请管理员添加")
			return
		}

		ctx.ReplyTextAndAt(strings.Join(folderInfos, "\n"))
	})

	engine.OnRegex(`^来点(.+)图片$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if hitDir(storageFolder, ctx.State["regex_matched"].([]string)[1]) {
			replyFixedCategory(ctx, storageFolder, ctx.State["regex_matched"].([]string)[1], smallAmount())
		} else {
			replySearch(ctx, storageFolder, ctx.State["regex_matched"].([]string)[1], smallAmount())
		}
	})

	engine.OnRegex(`^来(\d+)张(.+)图片$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		words := ctx.State["regex_matched"].([]string)
		num, _ := strconv.Atoi(words[1])
		if hitDir(storageFolder, words[2]) {
			replyFixedCategory(ctx, storageFolder, words[2], num)
		} else {
			replySearch(ctx, storageFolder, words[2], num)
		}
	})

	folders, _ := GetSubFolder(storageFolder)
	if len(folders) == 0 {
		return
	}

	folderNames := make([]string, 0, len(folders))
	for _, folder := range folders {
		folderNames = append(folderNames, folder.Name())
	}
	engine.OnFullMatchGroup(folderNames).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		replyFixedCategory(ctx, storageFolder, ctx.State["matched"].(string), smallAmount())
	})

	engine.OnRegex(fmt.Sprintf(`^(%s) ?(\d+)$`, strings.Join(folderNames, "|"))).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		words := ctx.State["regex_matched"].([]string)
		if num, err := strconv.Atoi(words[2]); err == nil {
			replyFixedCategory(ctx, storageFolder, words[1], num)
		}
	})
}

func dir(folders []os.DirEntry) []string {
	if len(folders) == 0 {
		return nil
	}
	folderInfos := make([]string, 0, len(folders))

	//遍历文件
	for _, folder := range folders {
		//获取文件名
		name := folder.Name()

		//跳过 . 和 ..
		if name == "." || name == ".." {
			continue
		}

		info, _ := folder.Info()
		//获取文件修改时间
		modTime := info.ModTime().Local()

		//格式化输出结果
		fileInfo := fmt.Sprintf("%s %s", name, modTime.Format("2006-01-02 15:04"))
		folderInfos = append(folderInfos, fileInfo)
	}
	return folderInfos
}

func replySearch(ctx *robot.Ctx, storageFolder, keywords string, num int) {
	if num <= 0 || keywords == "" {
		return
	}

	title, imageUrls := searchImage(storageFolder, keywords, num)
	reply(ctx, title, keywords, imageUrls)
}

func replyFixedCategory(ctx *robot.Ctx, storageFolder, category string, num int) {
	if num <= 0 || category == "" {
		return
	}

	title, imageUrls := getImage(storageFolder, category, num)
	reply(ctx, title, category, imageUrls)
}

func reply(ctx *robot.Ctx, title string, target string, imageUrls []string) {
	if title == "" {
		ctx.ReplyTextAndAt(fmt.Sprintf("获取%s失败，请稍后重试", target))
		return
	}

	ctx.ReplyTextAndAt(title)
	for _, url := range imageUrls {
		if strings.HasSuffix(url, pdfSuffix) {
			ctx.ReplyFile(url)
		} else {
			ctx.ReplyImage(url)
			if dur, err := time.ParseDuration(fmt.Sprintf("%vms", int(rand.Float64()*2000))); err == nil {
				// 等待 2s 内的一个随机数
				time.Sleep(dur)
			}
		}
	}
}

func hitDir(storageFolder, keywords string) bool {
	folders, _ := GetSubFolder(storageFolder)
	if len(folders) == 0 {
		return false
	}

	for _, folder := range folders {
		if folder.Name() == keywords {
			return true
		}
	}

	return false
}

func getImage(storageFolder, category string, num int) (string, []string) {
	categoryPath := filepath.Join(storageFolder, category)
	folders, _ := GetSubFolder(categoryPath)
	if len(folders) == 0 {
		return "", nil
	}
	title, imageFiles, err := get(categoryPath, folders[rand.Intn(len(folders))], num)
	if err != nil {
		return "", nil
	}

	if len(imageFiles) == 0 {
		// empty, then retry
		return getImage(storageFolder, category, num)
	}

	return title, imageFiles
}

func searchImage(storageFolder string, keywords string, num int) (string, []string) {
	folders, _ := GetSubFolder(storageFolder)
	if len(folders) == 0 {
		return "", nil
	}

	dirEntryVos := make([]DirEntryVo, 0)
	for _, folder := range folders {
		categoryPath := filepath.Join(storageFolder, folder.Name())
		subFolders, _ := ReadDir(categoryPath, func(dirEntry os.DirEntry) bool {
			return dirEntry.IsDir() && strings.Contains(strings.ToLower(dirEntry.Name()), strings.ToLower(keywords))
		})
		for _, subFolder := range subFolders {
			dirEntryVos = append(dirEntryVos, DirEntryVo{subFolder, folder})
		}
	}
	if len(dirEntryVos) == 0 {
		return "", nil
	}

	// 有可能目录下为空，最多查三次
	for i := 0; i < 3; i++ {
		dirEntryVo := dirEntryVos[rand.Intn(len(dirEntryVos))]
		categoryPath := filepath.Join(storageFolder, dirEntryVo.parentFolder.Name())
		title, imageFiles, err := get(categoryPath, dirEntryVo.folder, num)
		if err != nil {
			return "", nil
		}

		if len(imageFiles) != 0 {
			return title, imageFiles
		}
	}

	return "", nil
}

func get(targetPath string, entry os.DirEntry, num int) (string, []string, error) {
	title := entry.Name()
	topicPath := filepath.Join(targetPath, title)
	imageFiles, err := ReadDir(topicPath, IsImageFile)
	if err != nil {
		return "", nil, err
	}

	sort.Slice(imageFiles, func(i, j int) bool {
		left, _ := imageFiles[i].Info()
		right, _ := imageFiles[j].Info()
		return left.ModTime().Before(right.ModTime())
	})

	num = mathutil.Min(num, len(imageFiles))
	images := make([]string, 0, num)
	// 小于 10 张，直接返回
	if num <= 10 {
		for i := 0; i < num; i++ {
			images = append(images, "local://"+filepath.Join(topicPath, imageFiles[i].Name()))
		}
		return title, images, nil
	}

	// 大于 10 张，返回文件
	for i := 0; i < num; i++ {
		images = append(images, filepath.Join(topicPath, imageFiles[i].Name()))
	}
	pdfFilePath := filepath.Join(topicPath, title+"-"+strconv.Itoa(num)+pdfSuffix)
	if !Exist(pdfFilePath) {
		err := GeneratePdfFile(pdfFilePath, images)
		if err != nil {
			log.Errorf("生成 pdf 文件 %s 失败, err: %v", pdfFilePath, err)
			return "", nil, err
		}
	}
	pwd, _ := os.Getwd()
	return title, []string{filepath.Join(pwd, pdfFilePath)}, nil
}

func smallAmount() int {
	return rand.Intn(5) + 1
}

func GetStorageFolder() string {
	return storageFolder
}
