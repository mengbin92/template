package users

import (
	"context"
	"explorer/provider/cache"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAuthentication struct {
	UserID      uint
	Username    string
	Permissions []string
}

func (m *UserManager) generateToken(user *User) (string, error) {
	now := time.Now()

	ua := &UserAuthentication{
		UserID:      user.ID,
		Username:    user.Username,
		Permissions: user.Permissions,
	}
	uaBytes, err := sonic.Marshal(ua)
	if err != nil {
		return "", errors.Wrap(err, "marshal user authentication error")
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(m.AuthConfig.TokenExpireInSeconds.AsDuration().Abs())),
		Issuer:    m.AuthConfig.Issuer,
		Subject:   string(uaBytes),
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(m.AuthConfig.Algorithm), claims)
	tokenString, err := token.SignedString([]byte(m.AuthConfig.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, "create token error")
	}

	return tokenString, nil
}

func (m *UserManager) generateTokenCacheKey(id uint) string {
	return fmt.Sprintf("login_token::%d", id)
}

func (m *UserManager) Login(ctx context.Context, username, password string, db *gorm.DB) (string, error) {
	m.Logger.Infof("login user: %s", username)

	user := &User{}
	if err := db.Where("username = ?", username).First(user).Error; err != nil {
		m.Logger.Errorf("login user: %s, error: %v", username, err)
		return "", errors.Wrap(err, "login user failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		m.Logger.Errorf("login user: %s, password: %s, error: %v", username, password, err)
		return "", errors.Wrap(err, "login user failed")
	}

	token, err := m.generateToken(user)
	if err != nil {
		m.Logger.Errorf("login user: %s token error: %v", username, err)
		return "", errors.Wrap(err, "create login user token failed")
	}

	err = cache.GetRedisClient().Set(ctx, m.generateTokenCacheKey(user.ID), token, m.AuthConfig.TokenExpireInSeconds.AsDuration().Abs()).Err()
	if err != nil {
		m.Logger.Errorf("login user: %s token cache error: %v", username, err)
		return "", errors.Wrap(err, "cache login user token failed")
	}

	return token, nil
}

func (m *UserManager) RefreshToken(ctx context.Context, ua *UserAuthentication, db *gorm.DB) (string, error) {
	user := &User{}
	if err := db.Where("id = ?", ua.UserID).First(user).Error; err != nil {
		m.Logger.Errorf("refresh user: %s token error: %v", ua.Username, err)
		return "", errors.Wrap(err, "refresh user token failed")
	}
	
	token, err := m.generateToken(user)
	if err != nil {
		m.Logger.Errorf("refresh user: %s token error: %v", user.Username, err)
		return "", errors.Wrap(err, "refresh user token failed")
	}

	err = cache.GetRedisClient().Set(ctx, m.generateTokenCacheKey(user.ID), token, m.AuthConfig.TokenExpireInSeconds.AsDuration().Abs()).Err()
	if err != nil {
		m.Logger.Errorf("refresh user: %s token cache error: %v", user.Username, err)
		return "", errors.Wrap(err, "cache login user token failed")
	}

	return token, nil
}

func (m *UserManager) loadUserAuthenticationFromToken(ctx context.Context, tokenStr string) (*UserAuthentication, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.AuthConfig.SecretKey), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "new token parse function error")
	}

	// 验证 JWT 的有效性
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	uaStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid claims sub")
	}
	ua := &UserAuthentication{}
	err = sonic.Unmarshal([]byte(uaStr), ua)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal user authentication error")
	}
	return ua, nil
}

func (m *UserManager) Logout(ctx context.Context, id uint) error {
	m.Logger.Infof("logout user with id: %d", id)
	return cache.GetRedisClient().Del(ctx, m.generateTokenCacheKey(id)).Err()
}

func (m *UserManager) ValidToken(ctx context.Context, tokenStr string) (*UserAuthentication, error) {
	ua, err := m.loadUserAuthenticationFromToken(ctx, tokenStr)
	if err != nil {
		return nil, errors.Wrap(err, "load user authentication from token failed")
	}

	// // token 续期
	// rdb := cache.GetRedisClient()
	// tokenCacheKey := m.generateTokenCacheKey(ua.UserID)
	// found, err := cache.GetRedisClient().Get(ctx, tokenCacheKey).Result()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "check user token from cache failed")
	// }
	// if found != tokenStr {
	// 	return nil, errors.New("token cache not match")
	// }
	// if rdb.TTL(ctx, tokenCacheKey).Val().Abs() < m.AuthConfig.TokenExpireInSeconds.AsDuration().Abs() {
	// 	rdb.Expire(ctx, tokenCacheKey, m.AuthConfig.TokenExpireInSeconds.AsDuration().Abs())
	// }

	return ua, nil
}
