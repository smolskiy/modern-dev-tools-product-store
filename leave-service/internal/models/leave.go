// internal/model/leave.go
package model

import "time"

type Leave struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	Type       string    `json:"type"`   // например: "основной", "учебный"
	Status     string    `json:"status"` // например: "заявка", "одобрен", "отклонён"
	DateFrom   time.Time `json:"date_from"`
	DateTo     time.Time `json:"date_to"`
	CreatedAt  time.Time `json:"created_at"`
}
