/**
 * @Author: Anpw
 * @Description:
 * @File:  BankDetails
 * @Version: 1.0.0
 * @Date: 2021/7/24 22:43
 */

package models

import "time"

type BankDetails struct {
	BankMonthID               string `json:"bank_month_id" gorm:"primary_key;type:varchar(64)"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 *time.Time `sql:"index"`
	OccurTime                 string
	LoginUserID               string    `json:"login_user_id" gorm:"type:varchar(64)"`
	BankCardNo                string    `json:"bank_card_no" gorm:"type:varchar(64)"`
	BeforeDeposit             float64   `json:"before_deposit" gorm:"type:decimal(10,2)"`
	AfterDeposit              float64   `json:"after_deposit" gorm:"type:decimal(10,2)"`
	Difference                float64   `json:"difference" gorm:"type:decimal(10,2)"`
	BankDetailsFilePath       string    `json:"bank_details_file_path" gorm:"type:varchar(255)"`
	BankDetailsFileHash       string    `json:"bank_details_file_hash" gorm:"type:varchar(255)"`
	BankDetailsFileUploadTime time.Time `json:"bank_details_file_upload_time" gorm:"type:datetime"`
}
