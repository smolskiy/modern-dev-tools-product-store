// internal/model/employee.go
package model

type Employee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Position   string `json:"position"`
}
