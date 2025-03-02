package models

type Permission struct {
	Role string `json:"role"`
	Obj  string `json:"obj"`
	Act  string `json:"act"`
}

// Struct untuk request body
type PermissionContainerRequest struct {
	UserID        string `json:"user_id"`
	ContainerName string `json:"container_name"`
	Action        string `json:"action"` // Misal: "read", "start", "stop", dll.
}
