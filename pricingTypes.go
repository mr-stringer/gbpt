package main

type VmPrice struct {
	VmSku      string
	Location   string
	Currency   string
	PaygHrRate float32
	OneYrRi    float32
	ThreeYrRi  float32
}

const (
	pssd = iota
	pssdv2
)

type DiskPrice struct {
	DiskType uint
	Pssd     string
	Location string
	Size     uint
	Iops     uint
	Mbps     uint
	Price    float32
}
