package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/api"
	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/gen"

	"sigs.k8s.io/yaml"
)

var flagPodName, flagKubeconfig string
var flagVerbose bool

func init() {
	flag.StringVar(&flagKubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&flagPodName, "pod", "", "name of the pod to try to schedule")
	flag.BoolVar(&flagVerbose, "v", false, "verbose output")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: %s -kubeconfig <file> -pod <name> pod\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "       %s gen-example <shape>\n", os.Args[0])
	flag.PrintDefaults()
}

func genCapacityExample(shape string) {
	var pools []api.DevicePool

	switch shape {
	case "0":
		pools = gen.GenShapeZero(4)
	case "1":
		pools = gen.GenShapeOne(4)
	case "2":
		pools = gen.GenShapeTwo(4)
	case "3":
		pools = gen.GenShapeThree(4)
	default:
		fmt.Printf("unknown shape %q\n", shape)
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

	switch args[0] {
	case "gen-example":
		shape := "0"
		if len(args) > 1 {
			shape = args[1]
		}
		genCapacityExample(shape)
		break
	case "pod":
		fmt.Fprintf(flag.CommandLine.Output(), "not implemented yet\n")
		os.Exit(1)
		break
	default:
		usage()
		os.Exit(1)
	}
}
