package gen

import (
	"fmt"

	"github.com/kubernetes-sigs/wg-device-management/k8srm-prototype/pkg/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ptr[T any](val T) *T {
	var v T = val
	return &v
}

func genPoolForNumaNode(node, poolBase string, numaNode, count int, vendor, driver, model, firmwareVer, driverVer string) api.DevicePool {
	return api.DevicePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.DevMgmtAPIVersion,
			Kind:       "DevicePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s-%02d", node, poolBase, numaNode),
		},
		Spec: api.DevicePoolSpec{
			NodeName: &node,
			Driver:   driver,
			Attributes: []api.Attribute{
				{Name: "vendor", StringValue: ptr(vendor)},
				{Name: "model", StringValue: ptr(model)},
				{Name: "firmwareVersion", SemVerValue: ptr(api.SemVer(firmwareVer))},
				{Name: "driverVersion", SemVerValue: ptr(api.SemVer(driverVer))},
				{Name: "numa", StringValue: ptr(fmt.Sprintf("%d", numaNode))},
			},
			DeviceCount: count,
		},
	}
}

func genSimplePools(num, numa, count int, nodeBase, poolBase, vendor, driver, model, firmwareVer, driverVer string) []api.DevicePool {
	var pools []api.DevicePool
	for i := 0; i < num; i++ {
		node := fmt.Sprintf("%s-%02d", nodeBase, i)
		for nn := 0; nn < numa; nn++ {
			pools = append(pools, genPoolForNumaNode(node, poolBase, nn, count, vendor, driver, model, firmwareVer, driverVer))
		}
	}

	return pools
}

func GenShapeZero(num int) []api.DevicePool {
	return genSimplePools(num, 1, 2, "shape-zero", "foozer", "example.com", "example.com-foozer", "foozer-1000", "4.2.1-gen3", "1.8.2")
}

func GenShapeOne(num int) []api.DevicePool {
	return genSimplePools(num, 2, 2, "shape-one", "foozer", "example.com", "example.com-foozer", "foozer-1000", "4.2.1-gen3", "1.8.2")
}

func GenShapeTwo(num int) []api.DevicePool {
	return genSimplePools(num, 4, 4, "shape-two", "foozer", "example.com", "example.com-foozer", "foozer-4000", "4.2.1-gen7", "1.8.2")
}

func GenShapeThree(num int) []api.DevicePool {
	return genSimplePools(num, 4, 4, "shape-three", "barzer", "example.com", "example.com-barzer", "barzer-1000", "1.1.1", "1.8.2")
}

// 4 cpus/numa nodes
// Each CPU has two Foozers and two Barzers associated
func GenFoozerBarzerNodes(num int) []api.DevicePool {
	var pools []api.DevicePool
	for i := 0; i < num; i++ {
		node := fmt.Sprintf("%s-%02d", "shape-foozer-barzer", i)
		for nn := 0; nn < 4; nn++ {
			pools = append(pools, genPoolForNumaNode(node, "foozer", nn, 2, "example.com", "example.com-foozer", "foozer-1000", "4.2.1-gen3", "1.8.2"))
			pools = append(pools, genPoolForNumaNode(node, "barzer", nn, 2, "example.com", "example.com-barzer", "barzer-1000", "1.1.1", "1.8.2"))
		}
	}

	return pools
}
