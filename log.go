package main

import (
	"fmt"
)

const (
	LogInfo    = 0
	LogWarning = 1
	LogErr     = 2
)

type LogClass struct {
}

var Log LogClass

func (this *LogClass) Write(msg string, _type int) {
	fmt.Println(msg)
}
