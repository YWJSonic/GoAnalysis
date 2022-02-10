package goanalysis

import (
	"bytes"
	"codeanalysis/analysis/dao"
	"fmt"
)

// 語言分析
// func AnalysisStyle1(code string) {
// 	// var package map[string]*PackageInfo
// 	buf := bytes.NewBuffer([]byte(" "))
// 	buf.WriteString(code)
// 	s := source{
// 		buf: buf.Bytes(),
// 	}
// 	s.start()
// 	for {
// 		if s.r == s.e {
// 			break
// 		}

// 		s.next()
// 		switch s.rangeStr() {
// 		case "package":
// 			s.next()
// 			s.PackageInfo.Name = s.rangeStr()
// 		case "import":
// 			s.importSpec()
// 		case "type":
// 			s.TypeDeclarations()
// 		case "func":
// 			s.FunctionDeclarations()
// 		case "const":
// 			s.ConstantDeclarations()
// 		case "var":
// 			s.VariableDeclarations()
// 		default:
// 			fmt.Println("--- ", s.rangeStr())
// 		}
// 	}
// }

func AnalysisStyle2(packageInfo *dao.PackageInfo, code string) {
	// var package map[string]*PackageInfo
	buf := bytes.NewBuffer([]byte(" "))
	buf.WriteString(code)
	s := source{
		buf:         buf.Bytes(),
		PackageInfo: packageInfo,
	}
	fmt.Println("FileName:", packageInfo.FileNodes.Path())
	if packageInfo.FileNodes.Path() == "analysis/goanalysis/token.go" {
		fmt.Println("")
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
			case "package":
				s.next()
				s.PackageInfo.SetName(s.rangeStr())
			case "import":
				s.ImportDeclarations()

			case "type":
				s.TypeDeclarations()

			case "func":
				if s.buf[s.r+1] == '(' {
					s.MethodDeclarations()
				} else {
					s.FunctionDeclarations()
				}
			case "const":
				s.ConstantDeclarations()
			case "var":
				s.VariableDeclarations()

			default:
			}
		}
	}
	// for name, info := range s.PackageInfo.AllImportLink {
	// 	fmt.Println(name, info.Path)
	// }
	// for name, info := range s.PackageInfo.AllConstInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllVarInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllTypeInfos {
	// 	fmt.Println(name, info.GetTypeName())
	// }
	// for name, info := range s.PackageInfo.AllFuncInfo {
	// 	fmt.Println(name, info.GetTypeName())
	// }
}
