package solus

// FilterUsage represent available filters for fetching usage statistics.
type FilterUsage struct {
	filter
}

// ByPeriod filter usage by specified period (e.g., "2020-07").
func (f *FilterUsage) ByPeriod(period string) *FilterUsage {
	f.add("filter[period]", period)
	return f
}

// ByUserID filter usage by specified user ID.
func (f *FilterUsage) ByUserID(userID int) *FilterUsage {
	f.addInt("filter[user_id]", userID)
	return f
}

// ByBillingUserID filter usage by specified billing user ID.
func (f *FilterUsage) ByBillingUserID(billingUserID int) *FilterUsage {
	f.addInt("filter[billing_user_id]", billingUserID)
	return f
}

// WithTimezone sets timezone for usage statistics.
func (f *FilterUsage) WithTimezone(timezone string) *FilterUsage {
	f.add("timezone", timezone)
	return f
}
