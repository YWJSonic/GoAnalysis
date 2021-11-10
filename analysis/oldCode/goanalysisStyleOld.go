package oldCode

// func GoCodeLineAnalysis(codeLine []string) {
// 	lineCount := -1
// 	funcBlock := false
// 	importBlock := false
// 	allPackageMap := map[string]*PackageInfo{}
// 	var currentPackage *PackageInfo
// 	var currentFunc *FuncInfo
// 	var funcBlockCount uint
// 	for {
// 		lineCount++
// 		if len(codeLine) == lineCount {
// 			break
// 		}
// 		currentLine := codeLine[lineCount]
// 		if importBlock { // 分析 import 檔
// 			if currentLine == ")" {
// 				importBlock = false
// 				continue
// 			}
// 			// 清除格式美化
// 			cleanReplacer := strings.NewReplacer("\"", "", "\t", "")
// 			currentLine = cleanReplacer.Replace(currentLine)
// 			// import 參數
// 			filePathSplit := strings.Split(currentLine, " ")
// 			newName := ""
// 			// package 有重新命名
// 			if len(filePathSplit) > 1 {
// 				newName = filePathSplit[0]
// 			}
// 			// 取得資料夾名稱
// 			filePathSplit = strings.Split(currentLine, string(os.PathSeparator))
// 			pacakageName := filePathSplit[len(filePathSplit)-1]
// 			// package 關聯建立
// 			if packageInfo, ok := allPackageMap[pacakageName]; ok {
// 				packageLink := NewPackageLink(newName, packageInfo)
// 				currentPackage.ImportLink[packageLink.Name()] = packageLink
// 				allPackageMap[packageInfo.Name] = packageInfo
// 			} else {
// 				packageInfo = NewPackageInfo(pacakageName)
// 				packageLink := NewPackageLink(newName, packageInfo)
// 				currentPackage.ImportLink[packageLink.Name()] = packageLink
// 				allPackageMap[packageInfo.Name] = packageInfo
// 			}
// 		} else if funcBlock { // 在 func 區塊內
// 			funcBlock = false
// 			splitStr := strings.Split(currentLine, "//")
// 			currentLine = splitStr[0]
// 			if currentLine[len(currentLine)-1] == '{' {
// 				funcBlockCount++
// 			} else if currentLine[len(currentLine)-1] == '}' {
// 				funcBlockCount--
// 			}
// 		} else { // 不再任何區塊內
// 			splitStr := strings.Split(currentLine, " ")
// 			if splitStr[0] == "package" { // 分析 package 命名
// 				fmt.Println("fin Package ", splitStr[1])
// 				packageInfo := NewPackageInfo(splitStr[1])
// 				allPackageMap[packageInfo.Name] = packageInfo
// 				currentPackage = packageInfo
// 			} else if splitStr[0] == "import" { // 進入 import 區塊
// 				fmt.Println("fin import ", splitStr[1])
// 				if splitStr[1] == "(" {
// 					importBlock = true
// 				} else { // 單行 import 直接分析
// 					// 建立關聯
// 					filePathSplit := strings.Split(splitStr[1], string(os.PathSeparator)) // os不一樣方向不一樣
// 					pacakageName := filePathSplit[len(filePathSplit)-1]
// 					pacakgeInfo := NewPackageInfo(pacakageName)
// 					packageLink := NewPackageLink("", pacakgeInfo)
// 					currentPackage.ImportLink[packageLink.Name()] = packageLink
// 				}
// 			} else if splitStr[0] == "func" { // 進入方法區塊
// 				funcBlock = true
// 				funcBlockCount++
// 				// SplitFunc(currentLine)
// 				//建立關聯
// 				currentFunc = NewFuncInfo(splitStr[1])
// 				fmt.Println("fin func ", currentFunc)
// 			}
// 		}
// 	}
// }

// const (
// 	charSkipFlag   = '\\'
// 	charFlag       = '\''
// 	stringFlag     = '"'
// 	noteFlag       = "//"
// 	noteFlag2Start = "/*"
// 	noteFlag2End   = "*/"
// 	nextLineFlag   = '\n'
// )

// type tokenInfo struct {
// 	TokenType token
// 	Idx       int
// }

// func Scanner(code string) []tokenInfo {
// 	var delimitersList []tokenInfo
// 	for idx, ch := range code {
// 		switch ch {
// 		case '{':
// 			delimitersList = append(delimitersList, tokenInfo{
// 				TokenType: _Lbrace,
// 				Idx:       idx,
// 			})
// 		case '}':
// 			delimitersList = append(delimitersList, tokenInfo{
// 				TokenType: _Rbrace,
// 				Idx:       idx,
// 			})
// 		case '(':
// 			delimitersList = append(delimitersList, tokenInfo{
// 				TokenType: _Lparen,
// 				Idx:       idx,
// 			})
// 		case ')':
// 			delimitersList = append(delimitersList, tokenInfo{
// 				TokenType: _Lparen,
// 				Idx:       idx,
// 			})
// 		}
// 	}
// 	return delimitersList
// }

// func BlockFormat(code string) {
// 	var left []int                                      // 左區塊便是字元索引
// 	var lineCount int                                   // 遞增行數
// 	var lineStartIndex int                              // 每行起始索引
// 	var blocks []BlockInfo                              // 區塊資料
// 	var blockLeft int                                   // 區塊起始索引
// 	var isInChar, isInString, isInNote1, isInNote2 bool // 當前索引區域旗標
// 	var currentBlock BlockInfo                          // 當前處理區塊
// 	// 當前區域指標範圍
// 	localcation := LocalCationPoint_Global
// 	// fmt.Println(charFlag, stringFlag, noteFlag, charSkipFlag)
// 	length := len(code) // 程式碼長度
// 	for idx := 0; idx < length; idx++ {
// 		char := code[idx]
// 		// 字元串區域檢查
// 		if isInChar {
// 			if code[idx-1] != charSkipFlag {
// 				isInChar = false
// 			}
// 			continue
// 		} else if char == charFlag {
// 			isInChar = true
// 		}
// 		// 字串區域檢查
// 		if isInString {
// 			if code[idx-1] != charSkipFlag {
// 				isInString = false
// 			}
// 			continue
// 		} else if char == stringFlag {
// 			isInString = true
// 		}
// 		// 註解區域檢查
// 		if isInNote1 {
// 			if char == nextLineFlag {
// 				isInNote1 = false
// 			}
// 			continue
// 		} else if char == '/' {
// 			label := string([]byte{code[idx], code[idx+1]})
// 			if label == noteFlag {
// 				isInNote1 = true
// 			}
// 			if label == noteFlag2Start {
// 				isInNote2 = true
// 			}
// 		}
// 		// 註解區域 Style2 檢查
// 		if isInNote2 {
// 			label := string([]byte{code[idx], code[idx-1]})
// 			if label == noteFlag2End {
// 				isInNote2 = false
// 			}
// 			continue /*  */
// 		}
// 		switch char {
// 		case '\n':
// 			// 紀錄當前行數
// 			lineCount++
// 			// 紀錄當前行數開頭在文章內的index
// 			lineStartIndex = idx + 1
// 		case '{':
// 			if code[idx+1] == '}' {
// 				idx++
// 				continue
// 			}
// 			// 全域宣告物件
// 			if localcation == LocalCationPoint_Global {
// 				localcation = LocalCationPoint_Local
// 				left = append(left, lineStartIndex)
// 				// DefineFormat(code[lineStartIndex:idx])
// 				// 生成區塊資料
// 				currentBlock = BlockInfo{}
// 				currentBlock.HeadIndexRange = [2]int{lineStartIndex, idx}
// 				currentBlock.BodyIndexRange = [2]int{idx}
// 			} else { // 區域宣告物件
// 				left = append(left, idx)
// 			}
// 		case '}':
// 			leftLen := len(left)
// 			// 當左區塊與右區塊當好抵銷時為全域宣告
// 			if leftLen > 1 {
// 				left = left[:leftLen-1]
// 				break
// 			}
// 			// func 結束區塊 紀錄區塊起始與結束索引
// 			blockLeft, left = left[leftLen-1], left[:leftLen-1]
// 			endIdx := idx + 1
// 			// 更新區塊資料
// 			currentBlock.BodyIndexRange[1] = endIdx
// 			currentBlock.DataIndexRange = [2]int{blockLeft, endIdx}
// 			blocks = append(blocks, currentBlock)
// 			// 將指標為改為全域
// 			localcation = LocalCationPoint_Global
// 		}
// 	}
// 	for _, blockInfo := range blocks {
// 		// fmt.Println("-------")
// 		// fmt.Printf("%s\n", code[blockInfo.HeadIndexRange[0]:blockInfo.HeadIndexRange[1]])
// 		// fmt.Printf("%s\n", code[blockInfo.BodyIndexRange[0]:blockInfo.BodyIndexRange[1]])
// 		// fmt.Printf("%s\n", code[blockInfo.DataIndexRange[0]:blockInfo.DataIndexRange[1]])
// 		DefineFormat(code[blockInfo.HeadIndexRange[0]:blockInfo.HeadIndexRange[1]])
// 	}
// }
// func funcBlockToFuncInfo(code string) FuncInfo {
// 	var blockStart, blockEnd byte = '(', ')'
// 	// var skipStr string = string([]byte{blockStart, blockEnd})
// 	var splite []string
// 	var spliteStart int
// 	var funcInfo FuncInfo
// 	var isMethodExpression bool
// 	var isParamsFinish bool
// 	lenght := len(code)
// 	for i := 0; i < lenght; i++ {
// 		if code[i] == blockStart { // 發現開始區塊
// 			if i == spliteStart { // 起始就進入區塊
// 				isMethodExpression = true
// 				continue
// 			}
// 			if code[spliteStart:i] == " " || code[spliteStart:i] == "" { // 空格或是無資料
// 				continue
// 			}
// 			if name := strings.TrimSpace(code[spliteStart:i]); name != "" { // 取得方法名稱
// 				splite = append(splite, name)
// 				funcInfo.Name = name
// 			}
// 			spliteStart = i
// 		} else if code[i] == blockEnd { // 處理區塊資料
// 			endindex := i + 1
// 			data := strings.TrimSpace(code[spliteStart+1 : endindex-1])
// 			spliteStart = endindex
// 			if isMethodExpression { // 選擇器
// 				isMethodExpression = false
// 				// fmt.Println(data)
// 				methodParams := strings.Split(data, " ")
// 				me := NewMethodExpression(methodParams[0], nil) // TODO:取得指定 struct 指標寫入
// 				funcInfo.Method = me
// 			} else if !isParamsFinish { // 輸入參數
// 				isParamsFinish = true
// 				if data == "" { // 無參數
// 					continue
// 				}
// 				params := strings.Split(data, ",")
// 				var currentStruct *StructInfo
// 				paramsLen := len(params)
// 				for i := paramsLen; i > 0; i-- { // params 為了對應簡易宣告倒著處理回來 ex: a, b, c string
// 					splites := strings.Split(strings.TrimSpace(params[i-1]), " ")
// 					if len(splites) == 2 {
// 						currentStruct = nil // TODO:取得指定 struct 指標寫入
// 					}
// 					funcInfo.ParamsInPoint = append(funcInfo.ParamsInPoint, NewFuncParams(splites[0], currentStruct))
// 				}
// 				if paramsLen > 1 {
// 					for i, count := 0, paramsLen/2; i < count; i++ {
// 						funcInfo.ParamsInPoint[i], funcInfo.ParamsInPoint[paramsLen-i-1] = funcInfo.ParamsInPoint[paramsLen-i-1], funcInfo.ParamsInPoint[i]
// 					}
// 				}
// 			} else { // 輸出參數
// 				// fmt.Println(data)
// 				params := strings.Split(data, ",")
// 				var currentStruct *StructInfo
// 				for i := len(params); i > 0; i-- { // params 為了對應簡易宣告倒著處理回來 ex: a, b, c string
// 					splites := strings.Split(strings.TrimSpace(params[i-1]), " ")
// 					if len(splites) == 2 {
// 						currentStruct = nil // TODO:取得指定 struct 指標寫入
// 					}
// 					funcInfo.ParamsOutPoint = append(funcInfo.ParamsInPoint, NewFuncParams(splites[0], currentStruct))
// 				}
// 			}
// 		}
// 	}
// 	// 無區塊的輸出
// 	if spliteStart != lenght && code[spliteStart:lenght-1] != " " && code[spliteStart:lenght-1] != "" {
// 		strings.TrimSpace(code[spliteStart : lenght-1])
// 		paramsInfo := NewFuncParams("undefine", nil) /// TODO:取得指定 struct 指標寫入
// 		funcInfo.ParamsOutPoint = append(funcInfo.ParamsOutPoint, paramsInfo)
// 	}
// 	return funcInfo
// }
// func funcBlockSplite(code string) []string {
// 	var blockStart, blockEnd byte = '(', ')'
// 	// var skipStr string = string([]byte{blockStart, blockEnd})
// 	var splite []string
// 	var spliteStart int
// 	lenght := len(code)
// 	for i := 0; i < lenght; i++ {
// 		if code[i] == blockStart { // 發現開始區塊
// 			// 起始就進入區塊
// 			if i == spliteStart {
// 				continue
// 			}
// 			if code[spliteStart:i] == " " || code[spliteStart:i] == "" {
// 				continue
// 			}
// 			// 取得方法名稱
// 			if name := strings.TrimSpace(code[spliteStart:i]); name != "" {
// 				splite = append(splite, name)
// 			}
// 			spliteStart = i
// 		} else if code[i] == blockEnd {
// 			endindex := i + 1
// 			splite = append(splite, strings.TrimSpace(code[spliteStart:endindex]))
// 			spliteStart = endindex
// 		}
// 	}
// 	if spliteStart != lenght && code[spliteStart:lenght-1] != " " && code[spliteStart:lenght-1] != "" {
// 		splite = append(splite, strings.TrimSpace(code[spliteStart:lenght-1]))
// 	}
// 	return splite
// }

// // 分析宣告結構
// func DefineFormat(line string) {
// 	// fmt.Println("DefindAnalysis", line)
// 	defineSplite := strings.SplitN(line, " ", 2)
// 	defineType := defineSplite[0]
// 	switch defineType {
// 	case DefineType_Func:
// 		var funcInfo FuncInfo
// 		funcInfo = funcBlockToFuncInfo(defineSplite[1])
// 		fmt.Print(funcInfo)
// 	case DefineType_Type:
// 		// fmt.Println("type:", defineSplite)
// 		// fmt.Println(line)
// 	case DefineType_Var:
// 		// fmt.Println("var:", defineSplite)
// 		// fmt.Println(line)
// 	}
// }

// func (s *source) funcBody() string {
// 	txt := ""
// 	// fmt.Println("FuncBody Start")
// 	for {
// 		s.nextToken()
// 		s.ch = rune(s.buf[s.r])
// 		txt += s.rangeStr()
// 		if s.ch == '/' {
// 			// 進入註解範圍
// 			if s.nextCh() == '/' {
// 				s.nextTargetToken('\n')
// 			} else if s.nextCh() == '*' {
// 				for {
// 					s.nextTargetToken('*')
// 					if s.nextCh() == '/' {
// 						break
// 					}
// 				}
// 			}
// 		}
// 		if s.ch == '"' { // 字串範圍
// 			s.nextTargetToken('"')
// 			// txt += s.rangeStr()
// 		}
// 		if s.ch == '\'' {
// 			s.nextTargetToken('"')
// 			// txt += s.rangeStr()
// 		}
// 		if s.ch == '{' { // 子區塊
// 			// s.nextTargetToken('}')
// 			s.funcBody()
// 			txt += s.rangeStr()
// 			// s.b = s.r + 1
// 		}
// 		if s.ch == '}' { // 方法區塊結束
// 			break
// 		}
// 	}
// 	return txt
// }

// func (s *source) InStructOnType() (str string) {
// 	// s.nextCh()
// 	switch s.buf[s.r+1] {
// 	case '*': // PointerType
// 		s.nextCh()
// 		str = string(s.ch)
// 		str += s.InStructOnType()
// 	case '[': // ArrayType, SliceType
// 		if s.buf[s.r+2] == ']' { // slice type
// 			s.nextCh()
// 			str += string(s.ch)
// 			s.nextCh()
// 			str += string(s.ch)
// 			str += s.InStructOnType()
// 			return str
// 		} else { // array type
// 			str = s.OnArrayType()
// 		}
// 	case '<': // OutPutChanelType
// 		if s.buf[s.r+1] == '<' { // 單出
// 			s.next()
// 			str = s.rangeStr()
// 		} else {
// 			s.next()
// 			str = s.rangeStr()
// 		}
// 		str = str + " " + s.InStructOnType()
// 	case '.': // ArranType
// 		if string(s.buf[s.r+1:s.r+4]) == "..." {
// 			s.nextCh()
// 			s.nextCh()
// 			s.nextCh()
// 			str = "..."
// 			str += s.InStructOnType()
// 		}
// 	default:
// 		nextEndIdx := s.nextIdx()
// 		nextTokenIdx := s.nextTokenIdx()
// 		if nextEndIdx < nextTokenIdx {
// 			s.next()
// 			str = strings.TrimSpace(s.rangeStr())
// 		}
// 		if nextTokenIdx == nextEndIdx {
// 			s.nextToken()
// 		} else if s.buf[nextTokenIdx] == '.' {
// 			str = s.OnQualifiedIdentType()
// 		} else if s.buf[nextTokenIdx] == '(' {
// 			tmpStr := strings.TrimSpace(string(s.buf[s.r+1 : nextTokenIdx]))
// 			if tmpStr == "func" {
// 				str = s.OnFuncType('\n')
// 			}
// 		}
// 	}
// 	return
// }
