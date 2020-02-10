package types

type InstanceLimit struct {
	Total MachineSKU
	Individual map[string]DiskSKU
}

type DiskSKU struct {
	Name  string
	Size  int
	Limit IOLimit
}

type MachineSKU struct {
	Name     string
	Cached   IOLimit
	Uncached IOLimit
}

type IOLimit struct {
	IOPS       float64
	Throughput float64
}
