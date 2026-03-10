package solus

// FilterVPCNetworks represents available filters for fetching list of VPC networks.
type FilterVPCNetworks struct {
	filter
}

// ByName filter VPC networks by specified name.
func (f *FilterVPCNetworks) ByName(name string) *FilterVPCNetworks {
	f.add("filter[search]", name)
	return f
}
