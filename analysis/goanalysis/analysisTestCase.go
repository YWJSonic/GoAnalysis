package goanalysis

import (
	wwe "codeanalysis/load/project/goloader"
	"fmt"
	"sync"
	"time"
)

const ca1, cb1, cc1 = 0b1001, "\\foo" + "fqw", 3_5 + .52/4
const ca2, cb3, cc4 int = 0b1001, 443 + 333/2, 2 / 4 // a = 3, b = 4, c = "foo", untyped integer and string constants
const Pi float64 = 3.14159265358979323846
const zero = 0.0 // untyped floating-point constant
const (
	size int64 = 1024
	eof        = -1 // untyped integer constant
)
const u, v float32 = 0, 3 // u = 0.0, v = 3.0
const ()
const ca = 2 + 3.0 // a == 5.0   (untyped floating-point constant)
const (
	c1e = iota
	c2e
	c3e
)
const cb = 15 / 4        // b == 3     (untyped integer constant)
const c = 15 / 4.0       // c == 3.75  (untyped floating-point constant)
const Θ float64 = 3 / 2  // Θ == 1.0   (type float64, 3/2 is integer division)
const Π float64 = 3 / 2. // Π == 1.5   (type float64, 3/2. is float division)
const d = 1 << 3.0       // d == 8     (untyped integer constant)
const e = 1.0 << 3       // e == 8     (untyped integer constant)
// const f = int32(1) << 33    // illegal    (constant 8589934592 overflows int32)
// const g = float64(2) >> 1 // illegal    (float64(2) is a typed floating-point constant)
const h = "foo" > "bar"  // h == true  (untyped boolean constant)
const cj = true          // j == true  (untyped boolean constant)
const ck = 'w' + 1       // k == 'x'   (untyped rune constant)
const l = "hi"           // l == "hi"  (untyped string constant)
const cm = string(ck)    // m == "x"   (type string)
const Σ = 1 - 0.707i     //            (untyped complex constant)
const Δ = Σ + 2.0e-4     //            (untyped complex constant)
const Φ = iota*1i - 1/1i //            (untyped complex constant)

var GX [size]int = [size]int{2}

type Feature2 struct {
	SetMockWebsocketServer func(websocketResponseMap map[string][]byte)
	Run                    func(wwe.ANAT1)
}

type Feature3 struct {
	SetMockWebsocketServer func(websocketResponseMap map[string][]byte)
	Run                    func(wwe.ANAT1)
}

type GSA int

func (self GSA) M32(
	struct{},
	struct {
		t int
		s struct{}
		i wwe.ANAT1
	},
) (a, b *wwe.GoFileNode, c, d int) {
	self += 1
	return &wwe.GoFileNode{}, &wwe.GoFileNode{}, 0, 0
}
func M3(chan chan<- int) (*wwe.GoFileNode, *wwe.GoFileNode) {
	return &wwe.GoFileNode{}, &wwe.GoFileNode{}
}

type Gs = Ws
type Ws = []As

type (
	t1 [1024]byte
	t2 int
	// t3 []int
	t4 []*int
	t5 = t1
	t6 struct {
		T6_1 t1
		T6_2 t2
	}
	t7 interface{}
	t8 interface {
		T8_1(int) int
	}
	// t9 t1
)

func (As) m5(Bs, [1238]byte, float64, interface{}) (success bool) {
	return
}

var tf1 func(a, b int, z float64, opt ...interface{}) (success bool)
var tf2 func(fd, dss int, f wwe.ANAT1, w float64) (float64, *[]int)

var HandlerIdxgfwer, HandlerIdx12332, HandlerIdx5tt2 = 1_5.0, 32 + 2, 5 / 11 // 本次取用的 arangodb client index
var once sync.Once
var initialized bool

var a [1024]byte
var s uint = 33
var i = 1 << s         // 1 has type int
var j int32 = 1 << s   // 1 has type int32; j == 0
var k = uint64(1 << s) // 1 has type uint64; k == 1<<33
var m int = 1.0 << s   // 1.0 has type int; m == 1<<33
var n = 1.0<<s == j    // 1.0 has type int; n == true
var o = 1<<s == 2<<s   // 1 and 2 have type int; o == false
var p = 1<<s == 1<<33  // 1 has type int; p == true
// var u = 1.0 << s             // illegal: 1.0 has type float64, cannot shift
// var u1 = 1.0<<s != 0         // illegal: 1.0 has type float64, cannot shift
// var u2 = 1<<s != 1.0         // illegal: 1 has type float64, cannot shift
// var v float32 = 1 << s       // illegal: 1 has type float32, cannot shift
var w int64 = 1.0 << 33 // 1.0<<33 is a constant shift expression; w == 1<<33
// var x = a[1.0<<s]             // panics: 1.0 has type int, but 1<<33 overflows array bounds
var b = make([]byte, 1.0<<s) // 1.0 has type int; len(b) == 1<<33
var mm int = 1.0 << s        // 1.0 has type int; mm == 0
var oo = 1<<s == 2<<s        // 1 and 2 have type int; oo == true
var pp = 1<<s == 1<<33       // illegal: 1 has type int, but 1<<33 overflows int
// var xx = a[1.0<<s]            // 1.0 has type int; xx == a[0]
var bb = make([]byte, 1.0<<s) // 1.0 has type int; len(bb) == 0
var (
	WSX   = 3
	fgh   = 1.0
	r334f = 0x12
	hgqe  = 1 << 2
)

type G16Options struct {
	writeFlow chan *wwe.ANAT1 // 提供使用者發布訊息的通道
	ww        chan<- string

	// f
	// w2w             chan string

	s, b, a, d      string            // 123
	writeFlowBackup <-chan *wwe.ANAT1 // 暫存publish失敗的訊息，若重新連線成功再次發送
	writeEnd        chan struct{}     // 用於停止 publish() 生命週期
	shutdownChan    chan struct{}     // graceful shutdown 信號

	subscriberMap map[string]func() error // map[RoomGuid]func() error 記住所有註冊方法
	mu            sync.Mutex

	// read contains the portion of the map's contents that are safe for
	// concurrent access (with or without mu held).
	wwe.GoFileNode   `protobuf:"bytes,1,rep,name=ecSiteID,proto3" json:"ecSiteID,omitempty"` // 底注
	MinCoins         []wwe.GoFileNode                                                        `json:"minCoins"`         // 最低金額
	Rounds           []wwe.GoFileNode                                                        `json:"rounds"`           // 局數
	DisabledPatterns []wwe.GoFileNode                                                        `json:"disabledPatterns"` // 缺少的獎項
	BirdCatchCounts  []wwe.GoFileNode                                                        `json:"birdCatchCounts"`  // 抓鳥牌型
}

type Feature struct {
	SetMockWebsocketServer func(websocketResponseMap map[string][]byte)
	Run                    func(wwe.ANAT1)
}

func Test18(n, f func() int, w, fg func() int) func(p *int) {
	return func(p *int) {}
}
func Test8(n int, f func() int) func(p *int) {
	return func(p *int) {}
}

type Ds struct {
	WWe int
	asd string
}

func (Ds) m91(a, b, w, z float64, opt interface{}) (success bool) {
	return
}
func (As) m4(a, b int, w, z float64, opt ...interface{}) (success bool) {
	return
}
func (As) m9(prefix string, node ...*wwe.GoFileNode) (ww *wwe.GoFileNode, e *wwe.GoFileNode) {
	return &wwe.GoFileNode{}, &wwe.GoFileNode{}
}

type f1 func() int

type ()

type m1 map[string]interface{}
type m2 map[struct {
	p1, p2 [1024]byte
	p3     int
}]interface{}

type m3 []map[wwe.ANAT1]float32
type m4 []map[[1024]byte]map[string]float32

type g7 interface{}
type g8 interface {
	Test1([1024]byte) int
	Test2(int64) float32
}
type As struct {
	WWe int
	asd string
}

type gg1 = struct {
	Gw1 int
}

func GGB1() {
	var gga gg1
	fmt.Println(gga.Gw1)
}
func GGA1() {
	var gga gg1
	fmt.Println(gga.Gw1)
}

type g61 struct{}
type g9 g1

type g1 string
type g2 int
type g3 []int
type g4 []*int
type g5 = t1

func m8(node *wwe.GoFileNode) *wwe.GoFileNode {
	return &wwe.GoFileNode{}
}
func (self *As) M3_1() (*wwe.GoFileNode, *wwe.GoFileNode) {
	self.WWe += 1
	return &wwe.GoFileNode{}, &wwe.GoFileNode{}
}

func (As) M5(*Bs, int, float64) (ww chan string, m *[]int) {
	return make(chan string), nil
}

func (As) m1(a, _ int, z float32) bool {
	return false
}
func (As) m6(n int) func(p *As) int {
	return func(p *As) int { return 0 }
}

func (As) m7(n int) interface{} {
	return ""
}

func (As) m(x *int) {

}

func (As) m2(a, b int, z float32) bool {
	return false
}

type Bs struct{}

func Test6(a, b int, z float64, opt ...interface{}) (success bool) {
	return
}
func Test5(prefix string, values ...int) {

}
func Test7(int, int, float64) (float64, *[]int) {
	return 0, &[]int{}
}
func Test1() {

}
func Test2(x int) int {
	return 0
}
func Test3(a, _ int, z float32) bool {
	return false
}
func Test4(a, b int, z float32) bool {
	return true
}

// 實驗性質
// code 為程式碼片段 相當於一個類型
// 例如: map[string]interface{}
// 可以輸入
// "map[string]interface{}"
// "string"
// "interface{}"
// 但只輸入 map 會錯誤
func analysisType(code []byte) (str string) {
	tokenList := findToken(code)
	fmt.Println(tokenList)
	return
}

type GRedisAdapater interface {
	WW3()
	HGet(key, field string, dw int) (string, error)
	HGetAll(key string) (map[string]string, error)
	HMSet(key string, values ...interface{}) (bool, error)
	HDel(key string, fields ...string) (int64, error)

	Del(keys string) (int64, error)

	UpdateTotalBet(key string, playerBet map[string]string) error
	CreateGameRecordId(now time.Time) string
	Nw1(
		string,
		string,
		int64,
	) string
	Nw2(
		string,
		string,
		int64) string
	UpdateGame29Paradise(
		// fe
		ecSiteId string, gameId string, themeId string,
		memberIncomeOfWheel,
		memberOutcomeOfWheel,
		memberIncomeOfBanker,
		memberOutcomeOfBanker float64,
	) error
}

func findToken(code []byte) (tokenIndexList []int) {
	var r int
	for {
		switch code[r] {
		case '"':
			tokenIndexList = append(tokenIndexList, r)
		case '`':
			tokenIndexList = append(tokenIndexList, r)
		case '\'':
			tokenIndexList = append(tokenIndexList, r)
		case '(':
			tokenIndexList = append(tokenIndexList, r)
		case '{':
			tokenIndexList = append(tokenIndexList, r)
		case '[':
			tokenIndexList = append(tokenIndexList, r)
		case ';':
			tokenIndexList = append(tokenIndexList, r)
		case ',':
			tokenIndexList = append(tokenIndexList, r)
		case ']':
			tokenIndexList = append(tokenIndexList, r)
		case ')':
			tokenIndexList = append(tokenIndexList, r)
		case ':':
			tokenIndexList = append(tokenIndexList, r)
		case '}':
			tokenIndexList = append(tokenIndexList, r)
		case '+':
			tokenIndexList = append(tokenIndexList, r)
		case '.':
			tokenIndexList = append(tokenIndexList, r)
		case '*':
			tokenIndexList = append(tokenIndexList, r)
		case '-':
			tokenIndexList = append(tokenIndexList, r)
		case '%':
			tokenIndexList = append(tokenIndexList, r)
		case '/':
			tokenIndexList = append(tokenIndexList, r)
		case '|':
			tokenIndexList = append(tokenIndexList, r)
		case '&':
			tokenIndexList = append(tokenIndexList, r)
		case '<':
			tokenIndexList = append(tokenIndexList, r)
		case '^':
			tokenIndexList = append(tokenIndexList, r)
		case '=':
			tokenIndexList = append(tokenIndexList, r)
		case '>':
			tokenIndexList = append(tokenIndexList, r)
		case '~':
			tokenIndexList = append(tokenIndexList, r)
		case '!':
			tokenIndexList = append(tokenIndexList, r)
		}

		r++
		if r >= len(code) {
			break
		}
	}
	return
}

//
