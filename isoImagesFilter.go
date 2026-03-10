package solus

// FilterISOImages represents available filters for fetching list of ISO images.
type FilterISOImages struct {
	filter
}

// ByName filter ISO images by specified name.
func (f *FilterISOImages) ByName(name string) *FilterISOImages {
	f.add("filter[search]", name)
	return f
}
