package main

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "net/http/pprof"
)

func TestGolandAnalysis(t *testing.T) {
	// RunAnalysis()
	a1 := &A1{make(map[string]*B1)}
	b1 := &B1{A1: a1}
	a1.M1["b1"] = b1

	dat, err := json.Marshal(a1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(dat))
}

type A1 struct {
	M1 map[string]*B1
}

type B1 struct {
	A1 *A1 `json:"-"`
}
