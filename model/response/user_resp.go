package response

import "deploy/model"

type LoginResponse struct {
    User      model.User `json:"user"`
    Token     string     `json:"token"`
    ExpiresAt int64      `json:"expiresAt"`
}

type SysCaptchaResponse struct {
    CaptchaId string `json:"captchaId"`
    PicPath   string `json:"picPath"`
}
