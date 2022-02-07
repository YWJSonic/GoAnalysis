package goanalysis

import (
	"codeanalysis/analysis/dao"
	"strings"
)

func ConvertStringToType(data string) dao.ITypeInfo {

	if len(data) == 0 {
		return dao.BaseTypeInfo[_string]
	}

	switch data[0] {
	case '*':
	case '[':
	case '-': // 負號, int, float
		idx := strings.IndexAny(data, ".")
		if idx == -1 {
			return dao.BaseTypeInfo[_int]
		} else {
			return dao.BaseTypeInfo[_float64]
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // 數值, uint,int,flaot
		idx := strings.IndexAny(data, ".")
		if idx == -1 {
			return dao.BaseTypeInfo[_int]
		} else {
			return dao.BaseTypeInfo[_float64]
		}

	}
	panic("")
}
