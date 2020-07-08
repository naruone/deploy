package api

import (
    "deploy/config"
    "deploy/model"
    "deploy/model/request"
    "deploy/model/response"
    "deploy/router/middleware"
    "deploy/service"
    "deploy/utils"
    "errors"
    "github.com/dchest/captcha"
    "github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var (
        L          request.LoginStruct
        userVerify map[string][]string
        userModel  *model.User
        genToken   string
        expiresAt  int64
        err        error
    )
    _ = c.ShouldBindJSON(&L)
    userVerify = utils.Rules{
        "CaptchaId": {utils.NotEmpty()},
        "Captcha":   {utils.NotEmpty()},
        "Username":  {utils.NotEmpty()},
        "Password":  {utils.NotEmpty(), utils.Le("18"), utils.Ge("6")},
    }
    if err = utils.Verify(L, userVerify); err != nil {
        goto ERR
    }
    if captcha.VerifyString(L.CaptchaId, L.Captcha) == false {
        err = errors.New("验证码错误")
        goto ERR
    }
    if userModel, err = service.Login(&model.User{UserName: L.Username, Password: L.Password}); err != nil {
        goto ERR
    }
    if expiresAt, genToken, err = middleware.SignToken(userModel); err != nil {
        goto ERR
    }
    userModel.Password = "-"
    utils.OkWithData(response.LoginResponse{
        User:      *userModel,
        Token:     genToken,
        ExpiresAt: expiresAt,
    }, c)
    return

ERR:
    utils.FailWithMessage(err.Error(), c)
    return
}

//验证码生成
func Captcha(c *gin.Context) {
    captchaId := captcha.NewLen(config.GConfig.Captcha.Long)
    utils.OkDetailed(response.SysCaptchaResponse{
        CaptchaId: captchaId,
        PicPath:   "/auth/captcha/" + captchaId + ".png",
    }, "获取验证码成功", c)
}

//验证码图片读取
func CaptchaImg(c *gin.Context) {
    utils.GinCaptchaServeHTTP(c.Writer, c.Request)
}

func GetUserList(c *gin.Context) {
    var (
        pageInfo request.ComPageInfo
        list       []model.User
        total      int
        err        error
    )
    _ = c.ShouldBindJSON(&pageInfo)
    if list, total, err = service.UserList(&pageInfo); err != nil {
        utils.FailWithMessage("获取失败, Message: "+err.Error(), c)
        return
    }
    utils.OkDetailed(response.PageResult{
        List:        list,
        Total:       total,
        PageSize:    pageInfo.PageSize,
        CurrentPage: pageInfo.CurrentPage,
    }, "获取成功", c)
}

func SaveUser(c *gin.Context) {
    var (
        u          model.User
        userVerify map[string][]string
        err        error
    )
    _ = c.ShouldBindJSON(&u)
    userVerify = utils.Rules{
        "UserName": {utils.NotEmpty(), utils.Le("20")},
        "Password": {utils.Le("18")},
        "NickName": {utils.NotEmpty(), utils.Le("20")},
    }
    if err = utils.Verify(u, userVerify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.SaveUser(&u); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}

func ModifyPwd(c *gin.Context) {
    var (
        p      request.ModifyPwd
        verify map[string][]string
        err    error
    )
    u, _ := c.Get("claims")
    us := u.(*request.CustomClaims)
    _ = c.ShouldBindJSON(&p)
    verify = utils.Rules{
        "OldPwd":  {utils.NotEmpty(), utils.Le("18"), utils.Ge("6")},
        "NewPwd":  {utils.NotEmpty(), utils.Le("18"), utils.Ge("6")},
        "CNewPwd": {utils.NotEmpty(), utils.Le("18"), utils.Ge("6")},
    }
    if err = utils.Verify(p, verify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if p.NewPwd != p.CNewPwd {
        utils.FailWithMessage("密码确认不正确", c)
        return
    }
    p.UserId = int(us.ID)
    if err = service.UpdatePwd(&p); err != nil {
        utils.FailWithMessage("修改密码失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("修改成功!", c)
}

func SetUserStatus(c *gin.Context) {
    var (
        u   model.User
        err error
    )
    _ = c.ShouldBindJSON(&u)
    if err = service.SaveUser(&u); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}
