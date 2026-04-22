package request

import "github.com/golang-jwt/jwt/v4"

type LoginStruct struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    Captcha   string `json:"captcha"`
    CaptchaId string `json:"captchaId"`
}

type CustomClaims struct {
    ID       uint
    NickName string
    jwt.RegisteredClaims
}

type ModifyPwd struct {
    UserId  int
    OldPwd  string `json:"old_pwd"`
    NewPwd  string `json:"new_pwd"`
    CNewPwd string `json:"c_new_pwd"`
}
