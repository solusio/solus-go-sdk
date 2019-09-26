package solus

type ComputeResource struct {
	Id                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	CanRetryInstallation bool                   `json:"can_retry_installation"`
	Host                 string                 `json:"host"`
	AgentPort            string                 `json:"agent_port"`
	Status               ComputerResourceStatus `json:"status"`
}

type ComputerResourceStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
