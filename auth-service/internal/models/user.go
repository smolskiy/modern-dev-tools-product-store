// internal/model/user.go
package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // хранить только хеш!
	Role     string `json:"role"`               // например: admin, hr, manager, employee
}
