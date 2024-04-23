package solus

type AdditionalDiskCreateRequest struct {
	Name    string `json:"name"`
	Size    int    `json:"size"`
	OfferID int    `json:"offer_id"`
}

type Disk struct {
	ID         int     `json:"id"`
	IsPrimary  bool    `json:"is_primary"`
	Name       string  `json:"name"`
	Size       int     `json:"size"`
	ActualSize int     `json:"actual_size"`
	Path       string  `json:"path"`
	FullPath   string  `json:"full_path"`
	Storage    Storage `json:"storage"`
	Offer      Offer   `json:"offer,omitempty"`
}

type disksResponse struct {
	Data []Disk `json:"data"`
}
