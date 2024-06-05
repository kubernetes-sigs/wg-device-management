package gen

import (
	"fmt"
	"strings"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ptr[T any](val T) *T {
	var v T = val
	return &v
}

func genPool(node, pool string, devicesPerNuma, numaNodes int, vendor, driver, model, firmwareVer, driverVer string) api.DevicePool {
	var devices []api.Device
	for nn := 0; nn < numaNodes; nn++ {
		for d := 0; d < devicesPerNuma; d++ {
			devices = append(devices, api.Device{
				Name: fmt.Sprintf("dev-%02d", nn*devicesPerNuma+d),
				Attributes: []api.Attribute{
					{Name: "numa", StringValue: ptr(fmt.Sprintf("%d", nn))},
				},
			})
		}
	}

	return api.DevicePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: api.DevMgmtAPIVersion,
			Kind:       "DevicePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: pool,
		},
		Spec: api.DevicePoolSpec{
			NodeName: &node,
			Driver:   driver,
			Attributes: []api.Attribute{
				{Name: "vendor", StringValue: ptr(vendor)},
				{Name: "model", StringValue: ptr(model)},
				{Name: "firmwareVersion", SemVerValue: ptr(api.SemVer(firmwareVer))},
				{Name: "driverVersion", SemVerValue: ptr(api.SemVer(driverVer))},
			},
			Devices: devices,
		},
	}
}

type generator func(num int) []api.DevicePool

func getGenerators() map[string]generator {
	generators := make(map[string]generator)

	models := []string{"foozer-1000", "foozer-4000"}
	sizes := map[string]struct {
		numaNodes, devicesPerNuma int
	}{
		"tiny": {
			numaNodes:      1,
			devicesPerNuma: 2,
		},
		"small": {
			numaNodes:      1,
			devicesPerNuma: 4,
		},
		"medium": {
			numaNodes:      2,
			devicesPerNuma: 4,
		},
		"large": {
			numaNodes:      4,
			devicesPerNuma: 4,
		},
	}

	vendor := "example.com"
	driver := "example.com-foozer"
	firmwareVersion := "1.8.2"
	driverVersion := "4.2.1-gen3"

	for _, model := range models {
		for size, sizeInfo := range sizes {
			nodeBase := model + "-" + size
			generators[nodeBase] = func(num int) []api.DevicePool {
				return gen(num, nodeBase, "foozer", sizeInfo.devicesPerNuma, sizeInfo.numaNodes, vendor, driver, model, firmwareVersion, driverVersion)
			}
		}
	}

	generators["dgxa100"] = func(num int) []api.DevicePool {
		var pools []api.DevicePool
		for i := 0; i < num; i++ {
			nodeName := fmt.Sprintf("nvidia-%02d", i)
			p, err := dgxa100Pool(nodeName, "dgxa100")
			if err != nil {
				fmt.Printf("Error generating dgxa100 pool for %q: %s\n", nodeName, err.Error())
				continue
			}
			pools = append(pools, *p)
		}
		return pools
	}

	return generators
}

func Gen(nodeType string, num int) ([]api.DevicePool, error) {
	generators := getGenerators()

	generate, ok := generators[nodeType]
	if !ok {
		var valid []string
		for k := range generators {
			valid = append(valid, k)
		}
		return nil, fmt.Errorf("generator %q not found, valid generators are: %s", nodeType, strings.Join(valid, ", "))
	}

	return generate(num), nil
}

func gen(num int, nodeBase, poolBase string, devicesPerNuma, numaNodes int, vendor, driver, model, firmwareVer, driverVer string) []api.DevicePool {
	var pools []api.DevicePool
	for i := 0; i < num; i++ {
		nodeName := fmt.Sprintf("%s-%02d", nodeBase, i)
		poolName := fmt.Sprintf("%s-pool-%s", nodeBase, poolBase)

		pools = append(pools, genPool(nodeName, poolName, devicesPerNuma, numaNodes, vendor, driver, model, firmwareVer, driverVer))
	}

	return pools
}

// 4 cpus/numa nodes
// Each CPU has two Foozers and two Barzers associated
/*
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
*/
