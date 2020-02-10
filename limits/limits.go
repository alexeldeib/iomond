package limits

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/alexeldeib/imds"
	"github.com/alexeldeib/iomond/data"
	"github.com/alexeldeib/iomond/types"
)

func New() (*types.InstanceLimit, error) {
	deviceNames, err := getBlockDevices()
	if err != nil {
		return nil, err
	}

	metadata, err := imds.New()
	if err != nil {
		return nil, err
	}

	vmSKU, err := getMachineSKU(metadata.Compute.VMSize)
	if err != nil {
		return nil, err
	}

	fmt.Printf(
		"found vm sku %s with %.2f cached iops, %.2f MB/s  cached throughput and %.2f uncached iops, %.2f MB/s uncached throughput\n",
		vmSKU.Name,
		vmSKU.Cached.IOPS,
		vmSKU.Cached.Throughput,
		vmSKU.Uncached.IOPS,
		vmSKU.Uncached.Throughput,
	)

	diskSKUs := map[string]types.DiskSKU{}
	isEphemeral := metadata.Compute.StorageProfile.OsDisk.DiffDiskSettings.Option == "local"
	if !isEphemeral {
		deviceNames = append(deviceNames, "sda")
	}

	for _, device := range deviceNames {
		size, err := getDeviceSizeGB(device)
		if err != nil {
			return nil, err
		}

		sku, err := getDiskSKUFromSize(size, data.DiskSKUs)
		if err != nil {
			return nil, err
		}
		fmt.Printf("found disk /dev/%s with size %d, iops limit %.2f and throughput limit %.2f  MB/s\n", device, size, sku.Limit.IOPS, sku.Limit.Throughput)
		diskSKUs[device] = sku
	}
	return &types.InstanceLimit{
		Total:      vmSKU,
		Individual: diskSKUs,
	}, nil
}

func getBlockDevices() ([]string, error) {
	files, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	names := []string{}
	for i := range files {
		if strings.Contains(files[i].Name(), "loop") {
			continue
		}
		if strings.EqualFold(files[i].Name(), "sr0") {
			continue
		}
		if strings.EqualFold(files[i].Name(), "sda") {
			continue
		}
		if strings.EqualFold(files[i].Name(), "sdb") {
			continue
		}
		if strings.Contains(files[i].Name(), "nvme") {
			continue
		}
		names = append(names, files[i].Name())
	}
	return names, nil
}

func getDeviceSizeGB(device string) (int, error) {
	sectorStr, err := ioutil.ReadFile(fmt.Sprintf("/sys/block/%s/size", device))
	if err != nil {
		return -1, err
	}

	sectorCount, err := strconv.Atoi(strings.TrimRight(string(sectorStr), "\n"))
	if err != nil {
		return -1, err
	}
	// 1 << 30 == 1 GB
	// 512 = size of 1 sector in bytes, apparently statically defined
	return sectorCount * 512 / (1 << 30), nil
}

func getDiskSKUFromSize(size int, choices []types.DiskSKU) (types.DiskSKU, error) {
	for i := range choices {
		if size <= choices[i].Size {
			return choices[i], nil
		}
	}
	return types.DiskSKU{}, fmt.Errorf("no sku found for disk size: %d", size)
}

func getMachineSKU(skuName string) (types.MachineSKU, error) {
	for i := range data.MachineSKUs {
		candidate := data.MachineSKUs[i]
		if strings.EqualFold(candidate.Name, skuName) {
			return candidate, nil
		}
	}
	return types.MachineSKU{}, fmt.Errorf("failed to match imds sku name to vm size")
}
