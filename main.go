package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"

	"github.com/figassis/wpoffload/util"
)

func main() {
	fmt.Printf("Starting WP Offloader %s\n", util.Version)

	var wg sync.WaitGroup
	wg.Add(1)
	go util.Start(&wg)
	wg.Wait()

	util.Log("Done. Exiting")
	os.Exit(0)
}
