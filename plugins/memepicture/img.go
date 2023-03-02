package memepicture

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
)

// 将gif图片转换为png图片
func gif2Png(src, dst string) error {
	// 读入gif图片
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 解码gif图片
	gifImg, err := gif.DecodeAll(srcFile)
	if err != nil {
		return err
	}

	// 创建一张新的png图片
	pngImg := image.NewRGBA(image.Rect(0, 0, gifImg.Config.Width, gifImg.Config.Height))
	for y := pngImg.Bounds().Min.Y; y < pngImg.Bounds().Max.Y; y++ {
		for x := pngImg.Bounds().Min.X; x < pngImg.Bounds().Max.X; x++ {
			pngImg.Set(x, y, color.RGBA{})
		}
	}

	// 将gif图片中第一帧的像素值复制到png图片中
	for y := 0; y < gifImg.Config.Height; y++ {
		for x := 0; x < gifImg.Config.Width; x++ {
			pngImg.Set(x, y, gifImg.Image[0].At(x, y))
		}
	}

	// 保存png图片
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if err := png.Encode(dstFile, pngImg); err != nil {
		return err
	}
	return nil
}
