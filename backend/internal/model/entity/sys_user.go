package entity

import "time"

// SysUser 用户表实体
type SysUser struct {
	Id          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	PasswordHash string    `json:"-"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
