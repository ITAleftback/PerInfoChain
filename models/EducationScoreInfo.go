/**
 * @Author: Anpw
 * @Description:
 * @File:  EducationScoreInfo
 * @Version: 1.0.0
 * @Date: 2021/7/23 19:33
 */

package models

import "time"

type EducationScoreInfo struct {
	ScoreID     string `json:"score_id" gorm:"primary_key;type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	EducationID string     `json:"education_id" gorm:"type:varchar(64)"`
	LoginUserID string     `json:"login_user_id" gorm:"type:varchar(64)"`
	TestTime    uint       `json:"school_year" gorm:"type:varchar(20)"`
	ClassName   string     `json:"class_name" gorm:"type:varchar(20)"`
	Score       float32    `json:"score" gorm:"type:decimal(5,2)"`
	UploadTime  time.Time     `json:"upload_time" gorm:"type:datetime"`
	ScoreInfoFilePath       string `json:"score_info_file_path" gorm:"type:varchar(255)"`
	ScoreInfoFileHash       string `json:"score_info_file_hash" gorm:"type:varchar(255)"`
	ScoreInfoFileUploadTime time.Time `json:"score_info_file_upload_time" gorm:"type:datetime"`
}

func (EducationScoreInfo) TableName() string {
	return "education_score_info"
}
