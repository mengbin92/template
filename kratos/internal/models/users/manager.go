package users

import (
	"context"
	"explorer/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username    string
	Password    string
	Email       string
	Phone       string
	Permissions []string `gorm:"serializer:json"`
}

type UserManager struct {
	AuthConfig *conf.AuthConfig
	Logger     *log.Helper
}

func NewUserManager(authConfig *conf.AuthConfig, logger log.Logger) *UserManager {
	return &UserManager{
		AuthConfig: authConfig,
		Logger:     log.NewHelper(log.With(logger, "models", "user")),
	}
}

func (m *UserManager) CreateUser(ctx context.Context, user *User, db *gorm.DB) (*User, error) {
	m.Logger.Infof("create user: %s", user.Username)

	if err := db.Save(user).Error; err != nil {
		m.Logger.Errorf("create user: %s, error: %v", user.Username, err)
		return nil, errors.Wrap(err, "create user failed")
	}

	return user, nil
}

func (m *UserManager) GetUserByName(ctx context.Context, username string, db *gorm.DB) (*User, error) {
	m.Logger.Infof("get user: %s", username)

	var user User
	if err := db.First(&user, username).Error; err != nil {
		m.Logger.Errorf("get user: %s, error: %v", username, err)
		return nil, errors.Wrap(err, "get user failed")
	}

	return &user, nil
}

func (m *UserManager) GetUserByID(ctx context.Context, id uint, db *gorm.DB) (*User, error) {
	m.Logger.Infof("get user: %d", id)
	var user User
	if err := db.First(&user, id).Error; err != nil {
		m.Logger.Errorf("get user: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "get user failed")
	}

	return &user, nil
}

func (m *UserManager) UpdateUser(ctx context.Context, user *User, db *gorm.DB) error {
	m.Logger.Infof("update user: %d", user.ID)

	var oldUser User
	if err := db.First(&oldUser, user.ID).Error; err != nil {
		m.Logger.Errorf("update user: %s, error: %v", user.Username, err)
		return errors.Wrap(err, "update user failed")
	}

	// update user
	if user.Username != "" {
		oldUser.Username = user.Username
	}
	if user.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(user.Password)); err != nil {
			oldUser.Password = user.Password
		}
	}

	if  user.Email != "" {
		oldUser.Email = user.Email
	}
	if user.Phone != "" {
		oldUser.Phone = user.Phone
	}
	oldUser.UpdatedAt = time.Now()
	if err := db.Save(&oldUser).Error; err != nil {
		m.Logger.Errorf("update user: %s, error: %v", user.Username, err)
		return errors.Wrap(err, "update user failed")
	}

	return m.Logout(ctx, oldUser.ID)
}

func (m *UserManager) DeleteUser(ctx context.Context, id uint, db *gorm.DB) error {
	m.Logger.Infof("delete user: %d", id)
	return db.Delete(&User{}, id).Error
}

func (m *UserManager) ListUsers(ctx context.Context, db *gorm.DB) ([]*User, error) {
	m.Logger.Infof("list users")

	var users []*User
	if err := db.Find(&users).Error; err != nil {
		m.Logger.Errorf("list users, error: %v", err)
		return nil, errors.Wrap(err, "list users failed")
	}
	return users, nil
}
