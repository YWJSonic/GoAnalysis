package tool

import (
	"bytes"
	"io/ioutil"
)

func ReadFile(filePath string) string {
	dat, _ := ioutil.ReadFile(filePath)
	return string(dat)
}

func ReadFileToLineStr(filePath string) []string {

	var codeData []string
	dat, _ := ioutil.ReadFile(filePath)
	lines := bytes.Split(dat, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		codeData = append(codeData, string(line))
	}
	return codeData
}
