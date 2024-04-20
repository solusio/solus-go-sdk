package solus

const (
	// OfferTypeAdditionalDisk is a constant for additional disk offer type.
	OfferTypeAdditionalDisk OfferType = "additional_disk"
)

type OfferType string

type Offer struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description,omitempty"`
	Type               OfferType  `json:"type"`
	IsVisible          bool       `json:"is_visible"`
	AvailablePlans     []Plan     `json:"available_plans"`
	AvailableLocations []Location `json:"available_locations"`
}
