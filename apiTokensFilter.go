package solus

// FilterAPITokens represents available filters for fetching list of API tokens.
type FilterAPITokens struct {
	filter
}

// ByName filter API tokens by specified name.
func (f *FilterAPITokens) ByName(name string) *FilterAPITokens {
	f.add("filter[search]", name)
	return f
}

// ByUserID filter API tokens by specified user ID.
func (f *FilterAPITokens) ByUserID(userID int) *FilterAPITokens {
	f.addInt("filter[user_id]", userID)
	return f
}
