package v1

// RegisterReq 用户注册请求
type RegisterReq struct {
	Username string `json:"username" v:"required|length:3,64#用户名不能为空|用户名长度3-64位"`
	Email    string `json:"email"    v:"email#邮箱格式不正确"`
	Password string `json:"password" v:"required|length:8,128#密码不能为空|密码长度8-128位"`
}

// RegisterRes 用户注册响应
type RegisterRes struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

// LoginReq 用户登录请求
type LoginReq struct {
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

// LoginRes 用户登录响应
type LoginRes struct {
	Token string       `json:"token"`
	User  UserInfoRes  `json:"user"`
}

// UserInfoRes 用户信息响应
type UserInfoRes struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
