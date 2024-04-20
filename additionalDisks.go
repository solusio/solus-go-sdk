package solus

type AdditionalDiskCreateRequest struct {
	Name    string `json:"name"`
	Size    int    `json:"size"`
	OfferID int    `json:"offer_id"`
}

type AdditionalDisk struct {
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
