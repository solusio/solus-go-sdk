package solus

// FilterLanguages represents available filters for fetching list of languages.
type FilterLanguages struct {
	filter
}

// ByName filter languages by specified name.
func (f *FilterLanguages) ByName(name string) *FilterLanguages {
	f.add("filter[search]", name)
	return f
}
