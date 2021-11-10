package goanalysis

import (
	"bytes"
	"fmt"
)

// 語言分析
func AnalysisStyle1(code string) {
	// var package map[string]*PackageInfo
	buf := bytes.NewBuffer([]byte(" "))
	buf.WriteString(code)
	s := source{
		buf: buf.Bytes(),
	}
	s.start()
	for {
		if s.r == s.e {
			break
		}

		s.next()
		switch s.rangeStr() {
		case "package":
			s.next()
			s.PackageInfo.Name = s.rangeStr()
		case "import":
			s.importSpec()
		case "type":
			s.TypeDeclarations()
		case "func":
			s.FunctionDeclarations()
		case "const":
			s.ConstantDeclarations()
		case "var":
			s.VariableDeclarations()
		default:
			fmt.Println("--- ", s.rangeStr())
		}
	}
}

func AnalysisStyle2(packageInfo *PackageInfo, code string) {
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
			previrwCh := string(s.buf[s.r+1 : s.r+3])
			s.OnComments(previrwCh)
		} else {

			s.next()
			switch s.rangeStr() {
			case "package":
				s.next()
				s.PackageInfo.Name = s.rangeStr()
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
}
