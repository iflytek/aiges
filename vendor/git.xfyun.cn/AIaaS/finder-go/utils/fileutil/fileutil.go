package fileutil

import (
	"io/ioutil"
	"os"
	"strings"
	"fmt"
)

func ExistPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func GetSystemSeparator() string {
	s := "/"

	if os.IsPathSeparator('\\') {
		s = "\\"
	}

	return s
}

func WriteFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func IsTomlFile(name string) bool {
	return false
	//return strings.HasSuffix(name, ".toml")
}
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func ParseTomlFile(file []byte) map[string]interface{} {
	tomlConfig := make(map[string]interface{})
	content := string(file)

	var currentGroup string
	lineList := strings.Split(content, "\x0d\x0a")
	var fullLine []string
	for _, line := range lineList {
		fullLine = append(fullLine, strings.Split(line, "\n")...)
	}
	for _, line := range fullLine {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		} else if strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "[")  {
			if strings.Contains(line,"#"){
				line=strings.Split(line,"#")[0]
			}
			fmt.Println("line:",line)
			currentGroup = strings.TrimSpace(line[1 : len(line)-1])
			continue
		} else if !strings.Contains(line, "=") {
			continue
		} else {
			if strings.Contains(line, "#") {
				contentAndComment := strings.Split(line, "#")
				contentValue := strings.Split(strings.TrimSpace(contentAndComment[0]), "=")
				if len(currentGroup) == 0 {
					tomlConfig[strings.TrimSpace(contentValue[0])] = strings.TrimSpace(contentValue[1])
				} else {
					tomlConfig[currentGroup+"."+strings.TrimSpace(contentValue[0])] = strings.TrimSpace(contentValue[1])
				}
			} else {
				contentValue := strings.Split(strings.TrimSpace(line), "=")
				if len(currentGroup) == 0 {
					tomlConfig[strings.TrimSpace(contentValue[0])] = strings.TrimSpace(contentValue[1])
				} else {
					tomlConfig[currentGroup+"."+strings.TrimSpace(contentValue[0])] = strings.TrimSpace(contentValue[1])
				}
			}
		}
	}
	return tomlConfig
}
