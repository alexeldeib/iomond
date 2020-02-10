package data

import "github.com/alexeldeib/iomond/types"

var DiskSKUs = append([]types.DiskSKU(nil), PremiumSSDs...)
var MachineSKUs = append([]types.MachineSKU(nil), dsv2...)

var dsv2 = []types.MachineSKU{
	{
		Name: "Standard_DS1_v2",
		Cached: types.IOLimit{
			IOPS:       4000,
			Throughput: 32,
		},
		Uncached: types.IOLimit{
			IOPS:       3200,
			Throughput: 48,
		},
	},
	{
		Name: "Standard_DS2_v2",
		Cached: types.IOLimit{
			IOPS:       8000,
			Throughput: 64,
		},
		Uncached: types.IOLimit{
			IOPS:       6400,
			Throughput: 96,
		},
	},
	{
		Name: "Standard_DS3_v2",
		Cached: types.IOLimit{
			IOPS:       16000,
			Throughput: 128,
		},
		Uncached: types.IOLimit{
			IOPS:       12800,
			Throughput: 192,
		},
	},
	{
		Name: "Standard_DS4_v2",
		Cached: types.IOLimit{
			IOPS:       32000,
			Throughput: 256,
		},
		Uncached: types.IOLimit{
			IOPS:       25600,
			Throughput: 384,
		},
	},
	{
		Name: "Standard_DS5_v2",
		Cached: types.IOLimit{
			IOPS:       64000,
			Throughput: 512,
		},
		Uncached: types.IOLimit{
			IOPS:       51200,
			Throughput: 768,
		},
	},
}

var PremiumSSDs = []types.DiskSKU{
	{
		Name:       "P4",
		Size:       32,
		Limit: types.IOLimit{
			IOPS:       120,
			Throughput: 25,
		},
	},
	{
		Name:       "P6",
		Size:       64,
		Limit: types.IOLimit{
			IOPS:       240,
			Throughput: 50,
		},
	},
	{
		Name:       "P10",
		Size:       128,
		Limit: types.IOLimit{
			IOPS:       500,
			Throughput: 100,
		},
	},
	{
		Name:       "P15",
		Size:       256,
		Limit: types.IOLimit{
			IOPS:       1100,
			Throughput: 125,
		},
	},
	{
		Name:       "P20",
		Size:       512,
		Limit: types.IOLimit{
			IOPS:       2300,
			Throughput: 150,
		},
	},
	{
		Name:       "P30",
		Size:       1024,
		Limit: types.IOLimit{
			IOPS:       5000,
			Throughput: 200,
		},
	},
	{
		Name:       "P40",
		Size:       2048,
		Limit: types.IOLimit{
			IOPS:       7500,
			Throughput: 250,
		},
	},
	{
		Name:       "P50",
		Size:       4096,
		Limit: types.IOLimit{
			IOPS:       7500,
			Throughput: 250,
		},
	},
	{
		Name:       "P60",
		Size:       8192,
		Limit: types.IOLimit{
			IOPS:       16000,
			Throughput: 500,
		},
	},
	{
		Name:       "P70",
		Size:       16384,
		Limit: types.IOLimit{
			IOPS:       18000,
			Throughput: 750,
		},
	},
	{
		Name:       "P80",
		Size:       32767,
		Limit: types.IOLimit{
			IOPS:       20000,
			Throughput: 900,
		},
	},
}