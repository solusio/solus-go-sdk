package solus

type filter struct {
	data map[string]string
}

func (f *filter) get() map[string]string {
	return f.data
}

func (f *filter) add(k, v string) {
	if f.data == nil {
		f.data = map[string]string{}
	}

	f.data[k] = v
}
