package main

import (
	"flag"
	"fmt"
	"os"
)

var flagPodName, flagKubeconfig string
var flagVerbose bool

func init() {
	flag.StringVar(&flagKubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.BoolVar(&flagVerbose, "v", false, "verbose output")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [ -v ] -kubeconfig <file> <pod-name>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	podName := args[0]
	fmt.Fprintf(flag.CommandLine.Output(), "scheduling pod %q not implemented yet\n", podName)
}
