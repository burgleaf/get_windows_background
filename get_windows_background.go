package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	var localappdata string
	localappdata = os.Getenv("localappdata")
	userProfile := os.Getenv("UserProfile")
	myfolder := localappdata + `\Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets`
	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		srcName := myfolder + "/" + file.Name()
		//验证图片是否可以为背景
		imgConf := getImageWidthAndHeight(srcName)
		if imgConf == nil || imgConf.width != 1920 {
			fmt.Printf("该文件不合法,不作为桌面主题%s \n ", srcName)
			continue
		}
		copyFile(userProfile+`\Pictures\customdesktopbackground\`+file.Name()+`.jpg`, srcName)
	}
}

//copy文件
func copyFile(dstName string, srcName string) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	fmt.Printf("copy文件:%s到%s \n", srcName, dstName)
	io.Copy(dst, src)
}

//获取图片的长和宽
func getImageWidthAndHeight(imgPath string) *imgConf {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil
	}
	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y
	defer imgFile.Close()
	return &imgConf{width, height}
}

type imgConf struct {
	width, height int
}
