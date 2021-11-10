package main

import (
	"codeanalysis/load/project/goloader"
	"testing"
)

func TestGolandProjectLoad(t *testing.T) {
	goloader.LoadRoot("./")
}
