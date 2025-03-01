package models

type Permission struct {
	Role string `json:"role"`
	Obj  string `json:"obj"`
	Act  string `json:"act"`
}
