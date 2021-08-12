/**
 * @Author: Anpw
 * @Description:
 * @File:  BankInfo
 * @Version: 1.0.0
 * @Date: 2021/7/23 21:37
 */

package models

import "time"

type BankInfo struct {
	BankCardNo  string     `json:"bank_card_no" gorm:"primary_key;type:varchar(64)"`
	UserID      string `json:"user_id" gorm:"type:varchar(64)"`
	LoginUserID string `json:"login_user_id" gorm:"type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	BankName    string     `json:"bank_name" gorm:"type:varchar(20)"`
	Deposit     float64    `json:"deposit" gorm:"type:decimal(10,2)"`
	BankInfoFilePath       string `json:"bank_info_file_path" gorm:"type:varchar(255)"`
	BankInfoFileHash       string `json:"bank_info_file_hash" gorm:"type:varchar(255)"`
	BankInfoFileUploadTime time.Time `json:"bank_info_file_upload_time" gorm:"type:datetime"`
}

func (BankInfo) TableName() string {
	return "bank_info"
}
