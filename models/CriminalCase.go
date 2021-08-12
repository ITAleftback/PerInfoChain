/**
 * @Author: Anpw
 * @Description:
 * @File:  CriminalCase
 * @Version: 1.0.0
 * @Date: 2021/7/23 20:18
 */

package models

import "time"

type CriminalCase struct {
	CaseID      string `json:"case_id" gorm:"primary_key;type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	LoginUserID string     `json:"login_user_id" gorm:"type:varchar(64)"`
	UserID      string     `json:"user_id" gorm:"type:varchar(64)"`
	Time        string     `json:"time" gorm:"datetime"`
	Place       string     `json:"place" gorm:"type:varchar(255)"`
	Case        string     `json:"case" gorm:"type:varchar(255)"`
	Degree      uint       `json:"degree" gorm:"type:int(1)"`
	Punish      string     `json:"punish" gorm:"type:varchar(255)"`
	CaseInfoFilePath       string `json:"case_info_file_path" gorm:"type:varchar(255)"`
	CaseInfoFileHash       string `json:"case_info_file_hash" gorm:"type:varchar(255)"`
	CaseInfoFileUploadTime time.Time `json:"case_info_file_upload_time" gorm:"type:datetime"`
}

func (CriminalCase) TableName() string {
	return "criminal_case"
}
