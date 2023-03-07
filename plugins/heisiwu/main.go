package heisiwu

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"modernc.org/mathutil"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	categoryMap = map[string]string{
		"黑丝": "heisi",
		"白丝": "baisi",
		"巨乳": "juru",
		"美足": "meizu",
		"网红": "mcn",
		"jk": "jk",
	}
	categoryKeys = func() []string {
		keys := make([]string, 0, len(categoryMap))
		for k := range categoryMap {
			keys = append(keys, k)
		}
		return keys
	}()
	categoryMatch = strings.Join(categoryKeys, "|")
	categoryRegex = fmt.Sprintf(`^(%s) ?(\d+)$`, categoryMatch)
	pdfSuffix     = ".pdf"
)

func init() {
	engine := control.Register("heisiwu", &control.Options{
		Alias: "黑丝屋",
		Help: "指令:\n" +
			"* {" + categoryMatch + "} => 获取 1 张作品\n" +
			"* {黑丝 5} => 获取 5 张黑丝作品，限制 10 张\n" +
			"* {巨乳 3} => 获取 3 张巨乳作品，依此类推",
		DataFolder: "heisiwu",
	})

	storageFolder := engine.GetCacheFolder()

	engine.OnFullMatchGroup(categoryKeys).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		reply(ctx, storageFolder, ctx.State["matched"].(string), 1)
	})

	engine.OnRegex(categoryRegex).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		words := ctx.State["regex_matched"].([]string)
		if num, err := strconv.Atoi(words[2]); err == nil {
			reply(ctx, storageFolder, words[1], num)
		}
	})

	// 启动黑丝屋爬虫
	start(storageFolder)
}

func reply(ctx *robot.Ctx, storageFolder, category string, num int) {
	if num <= 0 || category == "" {
		return
	}

	title, imageUrls := getSetu(storageFolder, categoryMap[category], num)
	if title == "" {
		ctx.ReplyTextAndAt(fmt.Sprintf("获取%s作品失败，请稍后重试", category))
		return
	}

	ctx.ReplyTextAndAt(title)
	for _, url := range imageUrls {
		if strings.HasSuffix(url, pdfSuffix) {
			ctx.ReplyFile(url)
		} else {
			ctx.ReplyImage(url)
		}
		if dur, err := time.ParseDuration(fmt.Sprintf("%vms", int(rand.Float64()*2000))); err == nil {
			// 等待 2s 内的一个随机数
			time.Sleep(dur)
		}
	}
}

func getSetu(storageFolder, category string, num int) (string, []string) {
	categoryPath := filepath.Join(storageFolder, category)
	entries, err := GetSubFolder(categoryPath)

	if err != nil || len(entries) == 0 {
		return "", nil
	}

	title := entries[rand.Intn(len(entries))].Name()
	topicPath := filepath.Join(categoryPath, title)
	files, err := ReadDir(topicPath)
	if err != nil {
		return "", nil
	}
	if len(files) == 0 {
		// empty, then retry
		return getSetu(storageFolder, category, num)
	}

	imageFiles := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), pdfSuffix) {
			imageFiles = append(imageFiles, file)
		}
	}

	sort.Slice(imageFiles, func(i, j int) bool {
		left, _ := imageFiles[i].Info()
		right, _ := imageFiles[j].Info()
		return left.ModTime().Before(right.ModTime())
	})

	num = mathutil.Min(num, len(imageFiles))
	setus := make([]string, 0, num)
	// 小于 10 张，直接返回
	if num <= 10 {
		for i := 0; i < num; i++ {
			setus = append(setus, "local://"+filepath.Join(topicPath, imageFiles[i].Name()))
		}
		return title, setus
	}

	// 大于 10 张，返回文件
	for i := 0; i < num; i++ {
		setus = append(setus, filepath.Join(topicPath, imageFiles[i].Name()))
	}
	pdfFilePath := filepath.Join(topicPath, title+"-"+strconv.Itoa(num)+pdfSuffix)
	if !Exist(pdfFilePath) {
		err := generatePdfFile(pdfFilePath, setus)
		if err != nil {
			log.Errorf("生成 pdf 文件 %s 失败, err: %v", pdfFilePath, err)
			return "", nil
		}
	}
	pwd, _ := os.Getwd()
	return title, []string{filepath.Join(pwd, pdfFilePath)}
}

func generatePdfFile(filePath string, imagePaths []string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	// 添加图片
	for _, imagePath := range imagePaths {
		file, err := os.Open(imagePath) // 打开图片文件
		if err != nil {
			return err
		}
		img, _, err := image.Decode(file) // 解码图片文件的配置，包括宽高信息
		_ = file.Close()                  // 关闭文件
		// 等比例缩放
		ratio := float64(img.Bounds().Size().X) / float64(img.Bounds().Size().Y)
		width := float64(595)
		size := &gopdf.Rect{W: width, H: width / ratio}

		pdf.AddPage()
		err = pdf.Image(imagePath, 0, 0, size)
		if err != nil {
			return err
		}
	}
	return pdf.WritePdf(filePath)
}
