// internal/model/department.go
package model

type Department struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ParentID    *int   `json:"parent_id,omitempty"` // Для иерархии отделов (например, подотдел)
	ChiefID     *int   `json:"chief_id,omitempty"`  // ID руководителя отдела (связь с Employee)
	Description string `json:"description,omitempty"`
}
