package goloader

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 專案建構方法2（舊）
func Load(rootPath string) {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	dirs := []os.FileInfo{}
	goFiles := []os.FileInfo{}
	otherFiles := []os.FileInfo{}

	isGoProjectRoot := false

	// root file collection
	for _, f := range files {
		fileName := f.Name()

		if f.IsDir() {
			dirs = append(dirs, f)
		} else if IsGoFile(fileName) {
			if fileName == "main.go" {
				isGoProjectRoot = true
			}
			goFiles = append(goFiles, f)
		} else {
			otherFiles = append(otherFiles, f)
		}
	}

	// project type check
	if !isGoProjectRoot {
		return
	}
	option := Option{
		MaxDepth: 0,
		SkipFile: []string{".vscode", ".git", "vendor"},
	}
	newDirs, newGoFiles, newOtherFiles := loop(rootPath, dirs, 0, option)
	dirs = append(dirs, newDirs...)
	goFiles = append(goFiles, newGoFiles...)
	otherFiles = append(otherFiles, newOtherFiles...)
}

// 取得子節點
func loop(path string, dirs []os.FileInfo, depth int, option Option) (findDirs, goFiles, otherFiles []os.FileInfo) {
	newPath := path
	depth++
	tab := ""
	for i := 0; i < depth; i++ {
		tab += " "
	}
	isSkip := false
	for _, dir := range dirs {
		dirName := dir.Name()

		isSkip = false
		for _, skip := range option.SkipFile {
			if skip == dirName {
				isSkip = true
			}
		}
		if isSkip {
			continue
		}
		// fmt.Printf("%s- %s\n", tab, dirName)

		newPath = filepath.Join(path, dirName)
		files, err := ioutil.ReadDir(newPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fileName := f.Name()
			if f.IsDir() {

				findDirs = append(findDirs, f)

				if option.MaxDepth != 0 && option.MaxDepth >= depth {
					continue
				}

				newDirs, newGoFiles, newOtherFile := loop(newPath, []os.FileInfo{f}, depth, option)
				findDirs = append(findDirs, newDirs...)

				goFiles = append(goFiles, newGoFiles...)
				otherFiles = append(otherFiles, newOtherFile...)
			} else if IsGoFile(fileName) {
				goFiles = append(goFiles, f)
			} else {
				otherFiles = append(otherFiles, f)
			}
		}
	}

	return
}
