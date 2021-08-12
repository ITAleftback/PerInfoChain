/**
 * @Author: Anpw
 * @Description:
 * @File:  HousingInfo
 * @Version: 1.0.0
 * @Date: 2021/7/23 21:08
 */

package models

import "time"

type HousingInfo struct {
	HousingID   string `json:"housing_id" gorm:"primary_key;type:varchar(64)"`
	UserID      string `json:"user_id" gorm:"type:varchar(64)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	LoginUserID string     `json:"login_user_id" gorm:"type:varchar(64)"`
	Address     string     `json:"address" gorm:"type:varchar(255)"`
	HouseSize   uint       `json:"house_size" gorm:"type:int"`
	Price       float32    `json:"price" gorm:"type:decimal(10,2)"`
	HousingNo   string     `json:"housing_no" gorm:"type:varchar(255)"`
	PayTime     string     `json:"pay_time" gorm:"type:datetime"`
	UsingYear   uint       `json:"using_year" gorm:"type:int"`
	HouseInfoFilePath       string `json:"house_info_file_path" gorm:"type:varchar(255)"`
	HouseInfoFileHash       string `json:"house_info_file_hash" gorm:"type:varchar(255)"`
	HouseInfoFileUploadTime time.Time `json:"house_info_file_upload_time" gorm:"type:datetime"`
}

func (HousingInfo) TableName() string {
	return "housing_info"
}
