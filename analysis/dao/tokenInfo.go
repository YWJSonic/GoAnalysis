package dao

type TokenInfo struct {
	PointBase
}

func NewTokenInfo() *TokenInfo {
	return &TokenInfo{}
}

var BaseTokenInfo map[string]*TokenInfo = map[string]*TokenInfo{
	"+": {PointBase: PointBase{TypeName: "+"}},
}
