package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/lxn/win"
)

func main() {
	localappdata := os.Getenv("localappdata")
	userProfile := os.Getenv("UserProfile")
	myfolder := localappdata + `\Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets`
	files, _ := ioutil.ReadDir(myfolder)

	//自定义的主题图片目录
	customThemeFolder := userProfile + `\Pictures\customdesktopbackground\`
	if !createFolderIfNotExist(customThemeFolder) {
		fmt.Println("创建目录错误")
		return
	}

	var count = 0
	for _, file := range files {
		srcName := myfolder + "/" + file.Name()
		//验证图片是否可以为背景
		imgConf := getImageWidthAndHeight(srcName)
		if imgConf == nil || imgConf.width != 1920 {
			fmt.Printf("忽略宽度不合适的图片:%s \n ", srcName)
			continue
		}
		copyFile(customThemeFolder+file.Name()+`.jpg`, srcName)
		count++
	}

	title := "复制锁屏图片成功"
	content := "复制" + strconv.Itoa(count) + "张图片到[" + customThemeFolder + "]目录"
	popWindows(title, content)

}

func createFolderIfNotExist(folderPath string) bool {
	var _, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		var err = os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	return true
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
	fmt.Printf("复制文件:%s到%s \n", srcName, dstName)
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

func popWindows(title string, content string) {
	var hwnd win.HWND
	win.MessageBox(hwnd, _winText(content), _winText(title), win.MB_OK)
}

func _winText(_str string) *uint16 {
	return syscall.StringToUTF16Ptr(_str)
}

type imgConf struct {
	width, height int
}
