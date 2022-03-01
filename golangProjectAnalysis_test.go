package main

import (
	"testing"

	_ "net/http/pprof"
)

func TestGolandAnalysis(t *testing.T) {
	RunAnalysis()
}
