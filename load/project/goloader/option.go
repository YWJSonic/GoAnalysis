package goloader

func IsGoFile(fileName string) bool {
	nameCount := len(fileName)

	if nameCount < 4 {
		return false
	}

	if ".go" != fileName[nameCount-3:nameCount] {
		return false
	}

	return true
}
