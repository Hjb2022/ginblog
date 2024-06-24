package model

import (
	"Bluebell/utils/errmsg"
	"encoding/base64"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null " json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// 查询用户是否存在
func CheckUser(name string) (code int) {
	var users User
	db.Select("id").Where("username = ?", name).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.SUCCESS
}

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) (code int) {
	var user User
	db.Select("id, username").Where("username = ?", name).First(&user)
	if user.ID == uint(id) {
		return errmsg.SUCCESS
	}
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCESS
}

// 新增用户
func CreateUser(data *User) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询单个用户
func GetUser(id int) (User, int) {
	var user User
	err = db.Where("ID=?", id).First(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCESS
}

// 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int) {
	var total int
	var users []User

	if username != "" {
		db.Select("id,username,role,created_at").Where(
			"username LIKE ?", username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		db.Model(&users).Where(
			"username LIKE ?", username+"%",
		).Count(&total)
		return users, total
	}
	db.Select("id,username,role,created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	db.Model(&users).Count(&total)

	if err == gorm.ErrRecordNotFound {
		return nil, 0
	}
	return users, total

}

// 编辑用户,限定密码以外
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err = db.Model(&user).Where("id=?", id).Update(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除用户
func DeleteUser(id int) int {
	var user User
	err = db.Where("id=?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS

}

// 密码加密
func (u *User) BeforeSave() {
	u.Password = ScryptPw(u.Password)
}
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 6, 66, 22, 222, 11}
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw

}

//登录验证

func CheckLogin(username string, password string) int {
	var user User

	db.Where("username=?", username).First(&user)

	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCESS
}

// CheckLogin 后台登录验证
//func CheckLogin(username string, password string) (User, int) {
//	var user User
//	var PasswordErr error
//
//	db.Where("username = ?", username).First(&user)
//
//	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//
//	if user.ID == 0 {
//		return user, errmsg.ERROR_USER_NOT_EXIST
//	}
//	if PasswordErr != nil {
//		return user, errmsg.ERROR_PASSWORD_WRONG
//	}
//	if user.Role != 1 {
//		return user, errmsg.ERROR_USER_NO_RIGHT
//	}
//	return user, errmsg.SUCCESS
//}
//
//// CheckLoginFront 前台登录
//func CheckLoginFront(username string, password string) (User, int) {
//	var user User
//	var PasswordErr error
//
//	db.Where("username = ?", username).First(&user)
//
//	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//	if user.ID == 0 {
//		return user, errmsg.ERROR_USER_NOT_EXIST
//	}
//	if PasswordErr != nil {
//		return user, errmsg.ERROR_PASSWORD_WRONG
//	}
//	return user, errmsg.SUCCESS
//}
