/**
 * @Author: Anpw
 * @Description:
 * @File:  Menu
 * @Version: 1.0.0
 * @Date: 2021/7/25 22:15
 */

package models

type Menu struct {
	MenuID       string `json:"menu_id" gorm:"type:varchar(64)"`
	FatherMenuID string `json:"father_menu_id" gorm:"type:varchar(64)"`
	MenuName     string `json:"menu_name" gorm:"type:varchar(30)"`
	Grade        string `json:"grade" gorm:"type:varchar(8)"`
	Path         string `json:"path" gorm:"type:varchar(128)"`
}

func (Menu) TableName() string {
	return "menu"
}
