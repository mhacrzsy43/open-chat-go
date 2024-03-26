package models

import (
	"database/sql"
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime `gorm:"index" json:"-"`
	ID            uint         `json:"userId" gorm:"primaryKey;autoIncrement:true"`
	Name          string       `json:"name"`
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string `json:"token"`
	ClientIP      string `json:"clientIp"`
	ClientPort    string
	Salt          int32
	LoginTime     time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	HeartbeatTime time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	LogoutTime    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:Logout_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string { // 加入小括号表示这是一个函数
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)
	token, _ := utils.GenerateToken(user.ID)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", token)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) UserBasic {
	err := utils.DB.Create(&user)
	if err != nil {
		fmt.Println("插入用户失败", err)
		return UserBasic{}
	}
	return user
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		PassWord: user.PassWord,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}
