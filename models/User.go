/**
 * @Author: Anpw
 * @Description:
 * @File:  User
 * @Version: 1.0.0
 * @Date: 2021/7/12 16:57
 */

package models

import (
	"time"
)

// User Role一栏 0 代表未分配角色 1教育 2公安 3住房 4银行 5超级管理员
type User struct {
	LoginUserID string `json:"user_id" gorm:"primary_key;type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Username    string     `json:"username" gorm:"type:varchar(20)"`
	Password    string     `json:"password" gorm:"type:varchar(60)"`
	Phone       string     `json:"login_phone" gorm:"type:varchar(32);unique_index:phone"`
	LastLogin   time.Time  `json:"last_login" gorm:"type:datetime"`
	Role        uint       `json:"role" gorm:"type:char(1)"`
	Status      uint       `json:"status" gorm:"type:char(1)"`
}

func (User) TableName() string {
	return "user"
}
