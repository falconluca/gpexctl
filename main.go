package main

import (
	"flag"
	"github.com/golang/glog"
	"gpex/cmd"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	cmd.Execute()
}
