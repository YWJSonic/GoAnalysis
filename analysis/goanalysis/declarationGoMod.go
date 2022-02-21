package goanalysis

import (
	"bytes"
	"codeanalysis/analysis/dao"
	wwe "codeanalysis/load/project/goloader"
	"codeanalysis/util"
)

func GoAnalysisGoMod(node *wwe.GoFileNode) {
	// 初始化 module 資料
	moduleInfo := dao.NewModuleInfo()
	moduleInfo.CurrentFileNodes = node
	Instants.ModuleInfo = moduleInfo

	// 讀檔
	code := util.ReadFile(node.Path())
	buf := bytes.NewBuffer([]byte(" "))
	buf.WriteString(code)
	s := source{
		buf: buf.Bytes(),
	}

	s.start()
	for {
		s.toNextCh()
		if s.r+1 == s.e {
			break
		}

		// 最外層註解
		if s.CheckCommon() {
			s.OnComments(string(s.buf[s.r+1 : s.r+3]))
		} else {

			s.next()
			switch s.rangeStr() {
			case "module":
				// 解析 module name
				s.next()
				moduleInfo.ModuleName = s.rangeStr()
			case "go":
				// 解析 專案使用的 go 版本
				s.next()
				moduleInfo.GoVersion = s.rangeStr()
			case "require":
				// 解析使用到的 package
				if s.buf[s.r+1] == '(' {
					s.next()

					for {
						s.toNextCh()

						if s.buf[s.r+1] == ')' {
							break
						}

						// 取得 package path
						s.next()
						// path := s.rangeStr()

						// 取得 package version
						s.next()
						// packageVersion := s.rangeStr()

						// // 解析 預設名稱
						// // 判斷是否為特殊格式路徑
						// isgopkg := false
						// if strings.Index(path, "gopkg.in") == 0 {
						// 	isgopkg = true
						// }

						// // 封包名稱解析
						// var name string
						// if isgopkg {
						// 	// 版本控制格式-1
						// 	splitStr := strings.Split(path, "/")
						// 	name = splitStr[len(splitStr)-1]
						// 	name = name[:strings.LastIndexByte(name, '.')]

						// } else {
						// 	splitStr := strings.Split(path, "/")
						// 	name = splitStr[len(splitStr)-1]
						// 	isversiontFlag := 'v' == name[0]

						// 	// 版本控制格式-2
						// 	if vernumber, err := strconv.Atoi(name[1:]); isversiontFlag && err == nil && vernumber > 0 {
						// 		name = splitStr[len(splitStr)-2]
						// 	}

						// 	// 一般格式
						// 	if len(name) > 3 && name[:3] == "go-" {
						// 		name = name[3:]
						// 	} else if len(name) > 3 && name[len(name)-3:] == ".go" {
						// 		name = name[:len(name)-3]
						// 	}
						// }

						// if name == "" {
						// 	panic("import name error")
						// }

						// // 初始化 package 資料
						// packageInfo := dao.NewPackageInfo()
						// packageInfo.PackageVersion = packageVersion
						// packageInfo.SetName("")
						// moduleInfo.VendorMap[path] = packageInfo
					}
				}

			default:
			}
		}
	}
}
