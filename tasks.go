package solus

type Task struct {
	Id                int    `json:"id"`
	ComputeResourceId int    `json:"compute_resource_id"`
	Queue             string `json:"queue"`
	Action            string `json:"action"`
	Status            string `json:"status"`
	Output            string `json:"output"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
	StartedAt         Date   `json:"started_at"`
	CreatedAt         Date   `json:"created_at"`
	FinishedAt        Date   `json:"finished_at"`
}

type Date struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}
