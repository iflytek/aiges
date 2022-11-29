package mnist

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/xfyun/aiges/env"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var defaultMnistUrl = "https://ghproxy.com/?q=https%3A%2F%2Fgithub.com%2Fiflytek%2Faiges_demo%2Farchive%2Frefs%2Fheads%2Fmain.zip"
var md5sum = "938f93048ee06eae8eeed9ce32627111"
var zipFile = "aiges_demo.zip"

func getFileMd5(filename string) string {
	// 文件全路径名
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		fmt.Errorf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}
func InitMnistPythonWrapper() (err error) {
	if env.AIGES_PLUGIN_MODE != "python" {
		fmt.Println(fmt.Sprintf("Not support this mode. %s ", env.AIGES_PLUGIN_MODE))
		fmt.Println(fmt.Sprintf("Please use `export AIGES_PLUGIN_MODE=python ` "))
		os.Exit(0)
	}
	// 判断当前是否存在
	var found = false
	_, err = os.Stat(zipFile)
	if err == nil {
		m := getFileMd5(zipFile)
		if m != md5sum {
			fmt.Println("md5 check failed")
		} else {
			fmt.Println("found exists aiges_demo.zip ...")
			found = true
		}
	}
	if !found {
		// Get the data
		resp, err := http.Get(defaultMnistUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// 创建一个文件用于保存
		out, err := os.Create(zipFile)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		// 然后将响应流和文件流对接起来
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}

	dst := "aiges_demo"
	prefix := "aiges_demo-1.0.0"
	archive, err := zip.OpenReader(zipFile)
	defer archive.Close()
	fmt.Println("解压中demo压缩包...")
	for _, f := range archive.File {
		p := strings.TrimPrefix(f.Name, prefix)
		filePath := filepath.Join(dst, p)
		if filePath == "aiges_demo" {
			filePath += string(os.PathSeparator)
		}
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return
}
