package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (table *Contact) TableName() string { // 加入小括号表示这是一个函数
	return "contact"
}

func SearchFriend(token string) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	user := UserBasic{}
	utils.DB.Where("identity = ?", token).First(&user)
	utils.DB.Where("owner_id = ? and type = 1", user.ID).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>>>", v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}
