// internal/model/leave_type.go
package model

type LeaveType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
