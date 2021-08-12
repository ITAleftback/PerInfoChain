/**
 * @Author: Anpw
 * @Description:
 * @File:  Education
 * @Version: 1.0.0
 * @Date: 2021/7/22 11:53
 */

package models

import "time"

type EducationInfo struct {
	EducationID string `json:"education_id" gorm:"type:varchar(64);primary_key"`
	UserID      string `json:"user_id" gorm:"type:varchar(64);"`
	LoginUserID string `json:"login_user_id" gorm:"type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	SchoolName  string     `json:"school_name" gorm:"type:varchar(60)"`
	StartYear   uint       `json:"start_year" gorm:"type:int"`
	EndYear     uint       `json:"end_year" gorm:"type:int"`
	Major       string     `json:"major" gorm:"type:varchar(60)"`
	EducationInfoFilePath       string `json:"education_info_file_path" gorm:"type:varchar(255)"`
	EducationInfoFileHash       string `json:"education_info_file_hash" gorm:"type:varchar(255)"`
	EducationInfoFileUploadTime time.Time `json:"education_info_file_upload_time" gorm:"type:datetime"`
}

func (EducationInfo) TableName() string {
	return "education_info"
}
