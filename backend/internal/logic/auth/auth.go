package auth

import (
	"context"
	"errors"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
	"novel2script-backend/utility/jwt"
	"novel2script-backend/utility/password"
)

type authImpl struct{}

func init() {
	service.Auth = &authImpl{}
}

func (a *authImpl) Register(ctx context.Context, username, email, pass string) (*entity.SysUser, error) {
	// 检查用户名是否已存在
	existing, err := dao.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 哈希密码
	hash, err := password.HashPassword(pass)
	if err != nil {
		return nil, err
	}

	user := &entity.SysUser{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
		Status:       "active",
	}

	id, err := dao.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = id
	return user, nil
}

func (a *authImpl) Login(ctx context.Context, username, pass string) (*entity.SysUser, string, error) {
	user, err := dao.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	if !password.CheckPassword(pass, user.PasswordHash) {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, "", err
	}

	// 更新最后登录时间
	_ = dao.UpdateLastLogin(ctx, user.Id)

	return user, token, nil
}

func (a *authImpl) GetUserById(ctx context.Context, id int64) (*entity.SysUser, error) {
	return dao.GetUserById(ctx, id)
}
