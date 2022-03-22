package goanalysis

import "fmt"

func (s *source) onStatementList() {

}

func (s *source) onStatement() string {
	var statement string
	// var scanDone = false
	var offset = s.r
	var zoneTagMap = map[rune]rune{'{': '}', '(': ')', '[': ']'}
	var strZoneTagMap = map[rune]struct{}{'\'': {}, '`': {}, '"': {}}
	var zoneQueue = []rune{}
	var strZone = false

	for {

		if s.checkEnd() {
			break
		}

		s.nextCh()

		if !strZone {
			if zoneEndTag, ok := zoneTagMap[s.ch]; ok {
				zoneQueue = append(zoneQueue, zoneEndTag)
			} else if len(zoneQueue) > 0 && s.ch == zoneQueue[len(zoneQueue)-1] {
				zoneQueue = zoneQueue[:len(zoneQueue)-1]
			} else if _, ok := strZoneTagMap[s.ch]; ok {
				strZone = true
			}
		} else {
			if s.ch == '\\' {
				s.nextCh()
			} else if _, ok := strZoneTagMap[s.ch]; ok {
				strZone = false
			}
		}

		if s.isOnNewlineSymbol() {
			if len(zoneQueue) == 0 {
				break
			}
		} else if (s.ch == ' ' || s.ch == '\n' || s.ch == '\t') && s.checkCommon() {
			if len(zoneQueue) == 0 {
				break
			} else {
				s.nextLine()
			}
		}
	}

	statement = string(s.buf[offset:s.r])
	fmt.Println(statement)
	return statement
}
