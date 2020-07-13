package solus

func NewFilterLocations() *FilterLocations {
	return &filter{}
}

type FilterLocations = filter

func (f *FilterLocations) FilterByName(name string) *FilterLocations {
	f.add("filter[search]", name)

	return f
}
