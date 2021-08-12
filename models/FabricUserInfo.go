/**
 * @Author: Anpw
 * @Description:
 * @File:  FabricUserInfo
 * @Version: 1.0.0
 * @Date: 2021/7/19 0:23
 */

package models

import "time"

// FabricUserInfo 添加了一些唯一索引 性别一栏 1为男 2为女
type FabricUserInfo struct {
	UserID        string `json:"user_id" gorm:"type:varchar(64);primary_key"`
	LoginUserID   string `json:"login_user_id" gorm:"type:varchar(64)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	IDCard        string     `json:"id_card" gorm:"type:char(18);unique_index:id_card"`
	Name          string     `json:"name" gorm:"type:varchar(20)"`
	Sex           uint       `json:"sex"  gorm:"type:tinyint(1)"`
	Nation        string     `json:"nation"  gorm:"type:varchar(20)"`
	Native        string     `json:"native"  gorm:"type:varchar(20)"`
	Birthday      string     `json:"birth_day" gorm:"type:date"`
	Phone         string     `json:"phone" gorm:"type:char(11);unique_index:phone"`
	Email         string     `json:"email" gorm:"type:varchar(32);"`
	PoliticalLook string     `json:"political_look" gorm:"type:varchar(16)"`
	HomeAddress   string     `json:"home_address" gorm:"type:varchar(100)"`
	UserInfoFilePath       string `json:"user_info_file_path" gorm:"type:varchar(255)"`
	UserInfoFileHash       string `json:"user_info_file_hash" gorm:"type:varchar(255)"`
	UserInfoFileUploadTime time.Time `json:"user_info_file_upload_time" gorm:"type:datetime"`
}

func (FabricUserInfo) TableName() string {
	return "fabric_user_info"
}
