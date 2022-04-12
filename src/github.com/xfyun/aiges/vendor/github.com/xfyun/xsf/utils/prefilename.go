package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, string(os.PathSeparator)))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
func checkFile(fileName string) bool {
	f, e := os.Open(fileName)
	f.Close()
	if e != nil {
		return false
	}
	return true
}
func fileNamePreProcessing(fileName string) string {
	if checkFile(fileName) {
		//return getCurrentDirectory() + string(os.PathSeparator) + fileName
		return fileName
	}
	tmp := getParentDirectory(getCurrentDirectory()) + string(os.PathSeparator) + "conf" + string(os.PathSeparator) + fileName
	if checkFile(tmp) {
		return tmp
	}
	log.Fatalf("can't open %v\n", fileName)
	return ""
}
func FileNamePreProcessing(fileName string) string {
	return fileNamePreProcessing(fileName)
}
