package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/xianggong/m2svis/server/modules/database"
	"github.com/xianggong/m2svis/server/modules/router"
	"github.com/xianggong/m2svis/server/modules/trace"
)

func main() {
	// Commandline flags
	configPtr := flag.String("c", "./config.json", "Backend configuration file")
	tracePtr := flag.String("t", "", "Trace file to be processed")
	flag.Parse()

	// Database module is required
	err := database.Init(*configPtr)
	if err != nil {
		glog.Fatal("Initalize database failed")
		return
	}

	// 2 modes: processe trace or start the backend server
	if *tracePtr != "" {
		fmt.Println("Processing trace file:", *tracePtr)
		trace.Init()
		trace.Process(*tracePtr)
	} else {
		// Initialization
		router.Init()
	}
}
