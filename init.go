package main

import (
	"flag"
	"sync"
)

var (
	flgPort   int
	waitGroup sync.WaitGroup
)

func init() {
	parseFlags()
	waitGroup = sync.WaitGroup{}
	waitGroup.Add(1)
}

func parseFlags() {
	flag.IntVar(&flgPort, "port", 8282, "Port for flux connection")
	flag.Parse()
}
