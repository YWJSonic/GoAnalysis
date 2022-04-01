package goanalysis

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/analysis/dao"
	"codeanalysis/util"
	"fmt"
)

// 嘗試一次就將資料解析出來
// 失敗:
// 多個相同結構無法判斷
// ex:
// identifier "." : 後續可接續的東西過多無法辨識, 在其他資料未建立完成前無法已名稱查詢
// func (s *source) onVarExpressionList(infos []*dao.VarInfo) []string {
// 	var expressions []string
// 	for _, info := range infos {
// 		s.nextCh()
// 		exp := s.onVarExpression(info.ContentTypeInfo)
// 		info.Expression = exp
// 		if info.ContentTypeInfo == nil {
// 			// 從表達式解析型態
// 			if exp.ExpressionType != nil {
// 				info.ContentTypeInfo = exp.ExpressionType
// 			} else if exp.PrimaryExpr != nil {
// 				if exp.PrimaryExpr.Operand != nil {
// 					if exp.PrimaryExpr.Operand.Literal != nil {
// 						if exp.PrimaryExpr.Operand.Literal.CompositeLit != nil {
// 							info.ContentTypeInfo = exp.PrimaryExpr.Operand.Literal.CompositeLit.LiteralType
// 						}
// 					}
// 				}
// 			}
// 		}
// 		fmt.Println("-----------------------------")
// 		fmt.Printf("Name: %s = %s\n", info.GetName(), info.Expression.ContentStr)
// 		if s.ch == ',' {
// 			s.nextCh()
// 		}
// 	}
// 	return expressions
// }

func (s *source) onVarExpression(iInfo dao.ITypeInfo) *dao.Expression {

	offset := s.r
	exp := &dao.Expression{}

	for {
		switch s.ch {
		case '{':
			if s.buf[s.r+1] == '}' {
				// 结束符号在下一行
				s.nextCh()
				return nil
			}
			return exp

		case '[':
			if s.buf[s.r+1] == ']' {
				// CompositeLit
				liteType := s.OnTypeSwitch(string(s.ch))
				liteValue := s.onCompositeLit(liteType)

				prim := &dao.PrimaryExpr{
					Operand: &dao.Operand{
						Literal: &dao.Literal{
							CompositeLit: liteValue,
						},
					},
				}
				exp.PrimaryExpr = prim
				exp.ContentStr = string(s.buf[offset:s.r])
				return exp

			} else if string(s.buf[s.r+1:s.r+3]) == "..." {
				// dineamic array
				panic("")

			} else {
				// Expression array: [5+6], KeyType map: [string], Expression range: [1-1:8*6],
				if arrayInfo, ok := iInfo.(*dao.TypeInfoArray); ok {
					s.nextCh()
					subexp := s.onVarExpression(nil)
					exp.SubExpression = append(exp.SubExpression, *subexp)
					s.OnDeclarationsType()

					// CompositeLit
					liteValue := s.onCompositeLit(arrayInfo)

					prim := &dao.PrimaryExpr{
						Operand: &dao.Operand{
							Literal: &dao.Literal{
								CompositeLit: liteValue,
							},
						},
					}

					exp.PrimaryExpr = prim
					exp.ContentStr = string(s.buf[offset:s.r])
					return exp
				} else {
					panic("")
				}

			}
		case '\'':
			s.scanRuneLit()
			s.nextCh()
			exp.ExpressionType = iInfo
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp

		case '`':
			// 解析字串
			s.scanStringLit(s.ch)
			exp.ExpressionType = iInfo
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp

		case '"':
			// 解析字串
			s.scanStringLit(s.ch)
			if iInfo == nil {
				exp.ExpressionType = dao.BaseTypeInfo["string"]
			} else {
				exp.ExpressionType = iInfo
			}
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp

		case '+':
			panic("")
		case '-':
			panic("")
		case '!':
			panic("")
		case '^':
			panic("")
		case '*':
			// unary_op UnaryExpr, Expression binary_op Expression,
			typeInfo := s.OnTypeSwitch(string(s.ch))
			exp.ExpressionType = typeInfo
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp
		case '(':
			panic("")
		case '&':
			// 解析指標
			liteType := dao.NewTypeInfoPointer()
			liteType.SetTypeName("point")
			liteType.ContentTypeInfo = s.OnDeclarationsType()
			liteValue := s.onCompositeLit(liteType.ContentTypeInfo)

			prim := &dao.PrimaryExpr{
				Operand: &dao.Operand{
					Literal: &dao.Literal{
						CompositeLit: liteValue,
					},
				},
			}

			exp.PrimaryExpr = prim
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp

		case '<':
			panic("")
		case '0':
			ident := s.scanIdentifier()
			var basicLit dao.ITypeInfo

			if s.ch == '.' {
				// is float
				s.nextCh()
				ident = fmt.Sprintf("%s.%s", ident, s.scanIdentifier())
				basicLit = dao.BaseTypeInfo["float64"]

			} else {
				basicLit = dao.BaseTypeInfo["int"]
			}

			if s.ch != ',' && s.ch != '\n' && s.ch != '}' {
				// Expression binary_op Expression
				fstExp := dao.Expression{
					ExpressionType: basicLit,
					ContentStr:     string(s.buf[offset:s.r]),
				}
				exp.SubExpression = append(exp.SubExpression, fstExp)

				//排除 common, / 符號重疊
				if !s.checkCommon() {

					secExps := s.onSubBinaryOperator(&fstExp)
					exp.SubExpression = append(exp.SubExpression, secExps...)
				}
			} else {
				exp.ExpressionType = basicLit
			}
			exp.ContentStr = string(s.buf[offset:s.r])
			return exp

		default:
			if util.IsLetter(s.ch) {
				exp = s.onIdentifier(iInfo)
				exp.ContentStr = string(s.buf[offset:s.r])
				return exp
			} else if util.IsDecimal(s.ch) {
				exp = s.onExpressionNumber()
				exp.ContentStr = string(s.buf[offset:s.r])
				return exp
			} else {
				// 解析特殊符號
				panic("")
			}
		}
	}
}

func (s *source) onIdentifier(iInfo dao.ITypeInfo) *dao.Expression {
	exp := &dao.Expression{}
	// 解析 宣告名稱 或 型別名稱
	identifierOrType := s.scanIdentifier()

	if s.ch == '.' {
		// selector, TypeAssertion, QualifiedIdent, MethodExpr, CompositeLit
		if importPackage := s.PackageInfo.GetPackage(identifierOrType); importPackage != nil {

			s.nextCh()
			typeName := s.scanIdentifier()
			_, expType := s.PackageInfo.GetIdentifier(identifierOrType, typeName)

			if expType != nil {
				// 可判別類型
				switch info := expType.(type) {
				case *dao.TypeInfo:
					// TypeInfo: QualifiedIdent
					fmt.Println(info)

					switch s.ch {
					case '{':
						// CompositeLit
						liteType := info
						liteValue := s.onCompositeLit(liteType)
						if liteValue.LiteralValue == nil {
							exp.ExpressionType = liteType

						} else {
							prim := &dao.PrimaryExpr{
								Operand: &dao.Operand{
									Literal: &dao.Literal{
										CompositeLit: liteValue,
									},
								},
							}

							exp.PrimaryExpr = prim
						}
					case '(':
						// Conversion
						s.nextCh()
						subExp := s.onVarExpression(nil)
						pri := &dao.PrimaryExpr{
							Conversion: &dao.Conversion{
								ConversionType: info,
								SubExpression:  subExp,
							},
						}
						exp.ExpressionType = info
						exp.PrimaryExpr = pri
					}

					return exp

				case *dao.FuncInfo:
					// MethodExpr
					fmt.Println(info)

				}
			} else {
				// 無法判別, 存在於未解析檔案
				if s.ch == '{' {
					// struct
					_, liteType := s.PackageInfo.GetPackageType(identifierOrType, typeName)
					s.onCompositeLit(liteType)
				}
				fmt.Println("")
			}
			fmt.Println("")

			// primary := &dao.PrimaryExpr{
			// 	Operand: &dao.Operand{
			// 		OperandName: info,
			// 	},
			// }
			// s.onSubPrimaryExpr(primary)
			// exp.PrimaryExpr = primary
		}
		return exp
	} else if _, ok := constant.KeyWorkFunc[identifierOrType]; ok {
		// 字串為系統保留 方法名稱
		info := dao.BaseFuncInfo[identifierOrType]

		primary := &dao.PrimaryExpr{
			Operand: &dao.Operand{
				OperandName: info,
			},
		}

		// Arguments
		subExpList := dao.NewExpression()
		for {
			// 跨過 ',' 和 '('
			s.nextCh()
			s.toNextCh()

			subexp := s.onVarExpression(nil)
			subExpList.SubExpression = append(subExpList.SubExpression, *subexp)
			if s.ch != ',' {
				s.nextCh()
				break
			}
		}

		// Arguments
		s.onSubPrimaryExpr(primary)
		exp.PrimaryExpr = primary
		return exp

	} else if _, ok := constant.KeyWordType[identifierOrType]; ok {

		if identifierOrType == "func" {

			funcLit := s.onFuncLit()
			exp.PrimaryExpr = &dao.PrimaryExpr{
				Operand: &dao.Operand{
					Literal: &dao.Literal{
						FunctionLit: funcLit,
					},
				},
			}
		} else {

			// 字串為系統保留 類型名稱
			liteType := s.OnTypeSwitch(identifierOrType)
			liteValue := s.onCompositeLit(liteType)

			if liteValue.LiteralValue == nil {
				exp.ExpressionType = liteType

			} else {
				prim := &dao.PrimaryExpr{
					Operand: &dao.Operand{
						Literal: &dao.Literal{
							CompositeLit: liteValue,
						},
					},
				}

				exp.PrimaryExpr = prim
			}
		}
		return exp

	} else if baseConstInfo, ok := dao.BaseConstInfo[identifierOrType]; ok {
		// 字串為基礎型別
		exp.ExpressionType = baseConstInfo
		return exp

	} else if baseTypeInfo, ok := dao.BaseTypeInfo[identifierOrType]; ok {
		// 字串為基礎型別
		exp.ExpressionType = baseTypeInfo

		// Conversion
		if s.ch == '(' {
			s.nextCh()
			subexp := s.onVarExpression(nil)
			exp.SubExpression = append(exp.SubExpression, *subexp)
		}
		s.nextCh()
		return exp

	} else if constInfo, ok := s.PackageInfo.AllConstInfos[identifierOrType]; ok {
		// 字串為 const
		exp.ExpressionType = constInfo
		return exp

	} else if varInfo, ok := s.PackageInfo.AllVarInfos[identifierOrType]; ok {
		// 字串為 var
		exp.ExpressionType = varInfo
		return exp

	} else {
		if info, ok := iInfo.(*dao.TypeInfoStruct); ok {
			// 未查詢到結構名稱
			if subInfo, ok := info.VarInfos[identifierOrType]; ok {
				exp.ExpressionType = subInfo.ContentTypeInfo
			} else {
				for _, subInfo := range info.ImplicitlyVarInfos {
					if subInfo.ContentTypeInfo.GetName() == identifierOrType {
						exp.ExpressionType = subInfo.ContentTypeInfo
					}
				}
			}

		} else {
			// 字串可能是 未定義 const or var
			s.PackageInfo.UndefVarOrConst[identifierOrType] = info
			exp.ContentStr = identifierOrType
		}
		return exp
	}
}

func (s *source) onSubPrimaryExpr(primary *dao.PrimaryExpr) {
	subprimary := &dao.PrimaryExpr{}

	switch s.ch {
	case '.':
		s.nextCh()

	case '(':
		subprimary = s.onArguments()
	case '[':
	default:
		return
	}
	s.onSubPrimaryExpr(subprimary)
	primary.SubPrimaryExpr = subprimary
}

func (s *source) onCompositeLit(literalType dao.ITypeInfo) *dao.CompositeLit {
	var exp = &dao.CompositeLit{}

	// 初始化 LiteralType
	exp.LiteralType = literalType

	// 解析 LiteralValue
	if string(s.buf[s.r:s.r+2]) == "{}" {
		s.nextCh()
		s.nextCh()
		return exp
	}

	if s.ch == '{' {
		s.toNextCh()

		switch info := literalType.(type) {
		case *dao.TypeInfoSlice:
			exp.LiteralValue = s.onScanSliceElementList(info)
		case *dao.TypeInfoArray:
			exp.LiteralValue = s.onScanArrayElementList(info)
		case *dao.TypeInfoMap:
			exp.LiteralValue = s.onScanMapElementList(info)
		case *dao.TypeInfo:
			if info.ContentTypeInfo == nil {
				exp.LiteralValue = s.onScanStructUnknowElementList()
			} else if structInfo, ok := info.ContentTypeInfo.(*dao.TypeInfoStruct); ok {
				exp.LiteralValue = s.onScanStructElementList(structInfo)
			}
		case *dao.TypeInfoQualifiedIdent:
			if info.ContentTypeInfo == nil {
				exp.LiteralValue = s.onScanStructUnknowElementList()
			} else if structInfo, ok := info.ContentTypeInfo.(*dao.TypeInfoStruct); ok {
				exp.LiteralValue = s.onScanStructElementList(structInfo)
			}
		}

		if s.ch != '}' {
			s.nextCh()
		}
	}
	return exp
}

func (s *source) onScanSliceElementList(info *dao.TypeInfoSlice) *dao.ElementList {
	if s.buf[s.r+1] == '}' {
		// 提前結束
		s.nextCh()
		return nil
	}

	var elementList = &dao.ElementList{}
	for {
		if s.buf[s.r+1] == '}' {
			// 结束符号在下一行
			s.nextCh()
			break
		}
		s.toNextCh()
		s.nextCh()

		keyedElement := dao.KeyedElement{}
		keyedElement.SubElement = &dao.Element{
			SubExpression: s.onVarExpression(info.ContentTypeInfo),
		}

		elementList.KeyedElements = append(elementList.KeyedElements, keyedElement)
		if s.ch == ',' {
			s.toNextCh()
		} else if s.ch == '}' {
			// 结束符号在同一行
			break
		}
	}

	return elementList
}

func (s *source) onScanArrayElementList(info *dao.TypeInfoArray) *dao.ElementList {
	if s.buf[s.r+1] == '}' {
		// 提前結束
		s.nextCh()
		return nil
	}

	var elementList = &dao.ElementList{}
	for {
		if s.buf[s.r+1] == '}' {
			// 结束符号在下一行
			s.nextCh()
			break
		}
		s.toNextCh()
		s.nextCh()

		keyedElement := dao.KeyedElement{}
		keyedElement.SubElement = &dao.Element{
			SubExpression: s.onVarExpression(info),
		}

		elementList.KeyedElements = append(elementList.KeyedElements, keyedElement)
		if s.ch == ',' {
			s.toNextCh()
		} else if s.ch == '}' {
			// 结束符号在同一行
			break
		}
	}

	return elementList
}

func (s *source) onScanMapElementList(info *dao.TypeInfoMap) *dao.ElementList {
	if s.buf[s.r+1] == '}' {
		// 提前結束
		s.nextCh()
		return nil
	}

	var elementList = &dao.ElementList{}
	for {
		if s.buf[s.r+1] == '}' {
			// 结束符号在下一行
			s.nextCh()
			break
		}
		s.toNextCh()
		s.nextCh()

		// key
		keyedElement := dao.KeyedElement{}
		key := s.onVarExpression(info.KeyType)
		keyedElement.SubExpressionKey = &dao.Key{
			SubExpression: key,
		}

		if s.ch == ':' {
			// element Value
			s.toNextCh()
			s.nextCh()

			if string(s.buf[s.r:s.r+2]) == "{}" {
				s.nextCh()
				s.nextCh()
			} else {
				exp := s.onVarExpression(info.ValueType)
				if exp.PrimaryExpr != nil {
					keyedElement.SubElement = &dao.Element{
						SubExpression: exp,
					}
				} else if exp.SubExpression != nil {
					keyedElement.SubElement = &dao.Element{
						SubLiteralValue: exp.PrimaryExpr.Operand.Literal.CompositeLit.LiteralValue,
					}
				}
			}
		}

		elementList.KeyedElements = append(elementList.KeyedElements, keyedElement)
		if s.ch == ',' {
			s.toNextCh()

			if s.checkCommon() {
				s.OnComments(string(s.buf[s.r+1 : s.r+3]))
			}
		} else if s.ch == '}' {
			// 结束符号在同一行
			break
		}
	}

	return elementList
}

// 解析未知結構, 只紀錄名稱對應的關聯
func (s *source) onScanStructUnknowElementList() *dao.ElementList {
	if s.buf[s.r+1] == '}' {
		// 提前結束
		s.nextCh()
		return nil
	}

	var elementList = &dao.ElementList{}
	for {
		if s.buf[s.r+1] == '}' {
			// 结束符号在下一行
			s.nextCh()
			break
		}
		s.toNextCh()
		s.nextCh()

		// Key
		keyedElement := dao.KeyedElement{}
		key := s.scanIdentifier()
		keyedElement.SubExpressionKey = &dao.Key{
			SubFieldName: key,
		}

		if s.ch == ':' {
			// element Value
			s.toNextCh()
			s.nextCh()

			keyedElement.SubElement = &dao.Element{
				SubExpression: s.onVarExpression(nil),
			}
		}

		elementList.KeyedElements = append(elementList.KeyedElements, keyedElement)
		if s.ch == ',' {
			s.toNextCh()
		} else if s.ch == '}' {
			// 结束符号在同一行
			break
		}
	}

	return elementList
}

func (s *source) onScanStructElementList(info *dao.TypeInfoStruct) *dao.ElementList {
	if s.buf[s.r+1] == '}' {
		// 提前結束
		s.nextCh()
		return nil
	}

	var elementList = &dao.ElementList{}
	for {
		if s.buf[s.r+1] == '}' {
			// 结束符号在下一行
			s.nextCh()
			break
		}
		s.toNextCh()
		s.nextCh()

		// Key
		keyedElement := dao.KeyedElement{}
		elementOrKey := s.onVarExpression(info)
		keyedElement.SubExpressionKey = &dao.Key{
			SubExpression: elementOrKey,
		}

		if s.ch == ':' {
			// element Value
			s.toNextCh()
			s.nextCh()

			keyedElement.SubElement = &dao.Element{
				SubExpression: s.onVarExpression(info),
			}
		}

		elementList.KeyedElements = append(elementList.KeyedElements, keyedElement)
		if s.ch == ',' {
			s.toNextCh()
		} else if s.ch == '}' {
			// 结束符号在同一行
			break
		}
	}

	return elementList
}

func (s *source) onExpressionNumber() *dao.Expression {
	// 解析數字

	ident := s.scanIdentifier()
	var exp = &dao.Expression{}
	var basicLit dao.ITypeInfo

	if s.ch == '.' {
		// is float
		s.nextCh()
		ident = fmt.Sprintf("%s.%s", ident, s.scanIdentifier())
		basicLit = dao.BaseTypeInfo["float64"]

	} else {
		basicLit = dao.BaseTypeInfo["int64"]
	}

	if s.ch != ',' && s.ch != '\n' && s.ch != '}' {
		// Expression binary_op Expression
		fstExp := dao.Expression{
			ExpressionType: basicLit,
			ContentStr:     ident,
		}
		exp.SubExpression = append(exp.SubExpression, fstExp)
		secExps := s.onSubBinaryOperator(&fstExp)
		exp.SubExpression = append(exp.SubExpression, secExps...)

		if exp.ExpressionType == nil {
			var expType dao.ITypeInfo
			for _, subExp := range exp.SubExpression {
				if subExp.ExpressionType == nil {
					continue
				}

				if expType == nil {
					expType = subExp.ExpressionType
				} else if subExp.ExpressionType.GetName() == "float64" {
					expType = subExp.ExpressionType
				}

			}

			exp.ExpressionType = expType
		}

	} else {
		exp.ExpressionType = basicLit
	}
	return exp
}

// 表達式 解析二元運算
//
// @params *dao.Expression	第一部份表達式
//
// @params []dao.Expression 包含運算子以及第二部份表達式
// *如果無後續表達式回傳 nil
func (s *source) onSubBinaryOperator(fstExp *dao.Expression) []dao.Expression {
	var exps []dao.Expression

	//排除 common, / 符號重疊
	if s.checkCommon() {
		return exps
	}

	// 調整格式
	if s.ch == ' ' && s.buf[s.r+1] != ' ' {
		s.nextCh()
	}

	// 無運算子 or 結束符號
	if opType, len, ok := s.checkBinary_op(); ok {
		oper := string(s.buf[s.r : s.r+len])
		for i := 0; i < len; i++ {
			s.nextCh()
		}

		// 解析運算子
		opExp := dao.Expression{
			Operators:  oper,
			ContentStr: oper,
		}
		if _, ok := constant.Bool_op[opType]; ok {
			opExp.ExpressionType = dao.BaseTypeInfo["bool"]
		}

		exps = append(exps, opExp)

		// 調整格式
		if s.ch == ' ' && s.buf[s.r+1] != ' ' {
			s.nextCh()
		}

		// 解析第二部份表達式
		secExp := s.onVarExpression(nil)
		exps = append(exps, *secExp)

		//排除 common, / 符號重疊
		if !s.checkCommon() {

			// 第二部份後子表達式
			subExps := s.onSubBinaryOperator(secExp)
			exps = append(exps, subExps...)
		}

	}

	return exps
}

func (s *source) onArguments() *dao.PrimaryExpr {

	s.nextCh()
	subexp := &dao.Expression{}
	if s.ch != ')' {
		tmpexp := s.onVarExpression(nil)
		if s.ch == ',' {
			subexp.SubExpression = append(subexp.SubExpression, *tmpexp)

			for s.ch == ',' {
				s.toNextCh()
				s.nextCh()
				tmpexp = s.onVarExpression(nil)
				subexp.SubExpression = append(subexp.SubExpression, *tmpexp)
			}
		} else {
			subexp = tmpexp
		}
	}

	primary := &dao.PrimaryExpr{
		Arguments: &dao.Arguments{
			SubExpression: subexp,
		},
	}
	s.nextCh()
	return primary
}

// finc() "_"{}
func (s *source) onFuncLit() *dao.FunctionLit {
	info := &dao.FunctionLit{}
	info.ParamsInPoint, info.ParamsOutPoint = s.onSignature()
	info.Body = s.funcBodyBlock()
	return info
}
