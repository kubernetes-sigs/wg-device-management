package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/gen"

	"sigs.k8s.io/yaml"
)

var flagPodName, flagKubeconfig string
var flagVerbose bool

func init() {
	flag.BoolVar(&flagVerbose, "v", false, "verbose output")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "       %s <shape>\n", os.Args[0])
	flag.PrintDefaults()
}

func genCapacityExample(shape string) {
	pools := gen.Gen(shape, 2)
	if pools == nil {
		fmt.Printf("could not generate shape %q\n", shape)
	}

	if pools != nil {
		b, _ := yaml.Marshal(pools)
		fmt.Println(string(b))
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	genCapacityExample(args[0])
}
