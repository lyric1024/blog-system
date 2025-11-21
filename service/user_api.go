package service

import (
	"github.com/lyric1024/blog-system/model/system"
	"github.com/lyric1024/blog-system/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// 通过昵称和邮箱，校验用户是否存在
func (userService *UserService) ValidUser(userName, email string) error {

	var count int64
	err := userService.db.Model(&system.User{}).Where("user_name = ? OR email = ?", userName, email).Count(&count).Error
	if err != nil {
		return errors.Internal("检查用户昵称or邮箱是否存在失败", err)
	} else {
		if count > 0 {
			return errors.Exsit("用户已存在！", err)
		}
	}

	return nil
}

// 创建用户
func (userService *UserService) CreateUser(user system.User) (*system.User, error) {
	// 加密密码
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Internal("账号密码加密失败！", err)
	}

	createUser := user
	createUser.Password = string(hashPwd)
	if err := userService.db.Create(&createUser).Error; err != nil {
		return nil, errors.Internal("创建账号失败！", err)
	}

	returnUser := &createUser
	returnUser.Password = ""

	return returnUser, nil
}

// 创建用户
func (userService *UserService) Login(userName, password string) (*system.User, error) {

	var user system.User
	if err := userService.db.First(&user).Where("user_name = ?", userName).Error; err != nil {
		return nil, errors.NotFound("用户不存在", err)
	}

	// 对比psw
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.Unauthorized("密码错误", err)
	}

	returnUser := &user
	returnUser.Password = ""

	return returnUser, nil
}
