package solus

import (
	"context"
)

type TasksService service

const (
	// status
	TaskStatusPending  = "pending"
	TaskStatusQueued   = "queued"
	TaskStatusRunning  = "running"
	TaskStatusDone     = "done"
	TaskStatusFailed   = "failed"
	TaskStatusCanceled = "canceled"

	// actions
	ServerActionCreate         = "vm-create"
	ServerActionReinstall      = "vm-reinstall"
	ServerActionDelete         = "vm-delete"
	ServerActionUpdate         = "vm-update"
	ServerActionPasswordChange = "vm-password-change"
	ServerActionStart          = "vm-start"
	ServerActionStop           = "vm-stop"
	ServerActionRestart        = "vm-restart"
	ServerActionSuspend        = "vm-suspend"
	ServerActionResume         = "vm-resume"
)

type Task struct {
	Id                int    `json:"id"`
	ComputeResourceId int    `json:"compute_resource_id"`
	Queue             string `json:"queue"`
	Action            string `json:"action"`
	Status            string `json:"status"`
	Output            string `json:"output"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
}

type TasksResponse struct {
	paginatedResponse

	Data []Task `json:"data"`
}
type Date struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

// Tasks return list of Task, filter can be nil
func (s *TasksService) List(ctx context.Context, filter *FilterTasks) (TasksResponse, error) {
	resp := TasksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "tasks", &resp, withFilter(filter.data))
}
