package solus

import "strconv"

func NewFilterServers() *FilterServers {
	return &FilterServers{
		filter: map[string]string{},
	}
}

type FilterServers struct {
	filter map[string]string
}

func (f *FilterServers) Get() map[string]string {
	return f.filter
}

func (f *FilterServers) FilterByUserID(id int) *FilterServers {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["user_id"] = strconv.Itoa(id)

	return f
}

func (f *FilterServers) FilterByComputeResourceID(id int) *FilterServers {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["compute_resource_id"] = strconv.Itoa(id)

	return f
}

func (f *FilterServers) FilterByStatus(status string) *FilterServers {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["status"] = status

	return f
}
