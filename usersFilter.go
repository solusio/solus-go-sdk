package solus

func NewFilterUsers() *FilterUsers {
	return &FilterUsers{
		filter: map[string]string{},
	}
}

type FilterUsers struct {
	filter map[string]string
}

func (f *FilterUsers) Get() map[string]string {
	return f.filter
}

func (f *FilterUsers) FilterByStatus(status string) *FilterUsers {
	if f.filter == nil {
		f.filter = map[string]string{}
	}

	f.filter["status"] = status

	return f
}
