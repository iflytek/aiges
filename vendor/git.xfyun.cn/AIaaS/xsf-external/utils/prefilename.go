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

func getParentDirectory(directory string) string {
	return substr(directory, 0, strings.LastIndex(directory, string(os.PathSeparator)))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if nil != err {
		log.Fatal(err)
	}
	return dir
}
func checkFile(fileName string) bool {
	f, e := os.Open(fileName)
	if nil != e {
		return false
	}
	_ = f.Close()
	return true
}
func fileNamePreProcessing(fileName string) string {
	if checkFile(fileName) {
		//return getCurrentDirectory() + string(os.PathSeparator) + fileName
		return fileName
	}
	tmp := getParentDirectory(getCurrentDirectory()) +
		string(os.PathSeparator) +
		"conf" +
		string(os.PathSeparator) +
		fileName
	if checkFile(tmp) {
		return tmp
	}
	log.Fatalf("can't open %v\n", fileName)
	return ""
}
func FileNamePreProcessing(fileName string) string {
	return fileNamePreProcessing(fileName)
}
