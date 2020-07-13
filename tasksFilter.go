package solus

import "strconv"

func NewFilterTasks() *FilterTasks {
	return &FilterTasks{
		filter: map[string]string{},
	}
}

type FilterTasks struct {
	filter map[string]string
}

func (f *FilterTasks) Get() map[string]string {
	return f.filter
}

func (f *FilterTasks) FilterByAction(action string) *FilterTasks {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["filter[action]"] = action

	return f
}

func (f *FilterTasks) FilterByStatus(status string) *FilterTasks {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["filter[status]"] = status

	return f
}

func (f *FilterTasks) FilterByComputeResourceID(id int) *FilterTasks {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["filter[compute_resource_id]"] = strconv.Itoa(id)

	return f
}

func (f *FilterTasks) FilterByComputeResourceVmID(id int) *FilterTasks {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["filter[compute_resource_vm_id]"] = strconv.Itoa(id)

	return f
}
