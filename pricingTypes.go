package main

// VmPrice is a struct that stores all possible prices for a VM at a given
// location.
type VmPrice struct {
	VmSku      string
	Location   string
	Currency   string
	PaygHrRate float32
	OneYrRi    float32
	ThreeYrRi  float32
}

// constant that store values for pssd types
const (
	pssd = iota
	pssdv2
)

// DiskPrice is a struct that stores disk prices for a specific disk in a
// specific location. It can be used with either pssd pssdv2 types.
type DiskPrice struct {
	DiskType uint
	Pssd     string
	Location string
	Size     uint
	Iops     float32
	MBps     float32
	GBs      float32
	Price    float32
}
