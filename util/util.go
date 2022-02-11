package util

import (
	"bytes"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
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

func IsLetter(ch rune) bool {
	return 'a' <= Lower(ch) && Lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func IsDigit(ch rune) bool {
	return IsDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func Lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func IsDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func IsHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= Lower(ch) && Lower(ch) <= 'f' }
