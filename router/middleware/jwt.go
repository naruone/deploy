package middleware

import (
    "deploy/config"
    "deploy/model"
    "deploy/model/request"
    "deploy/utils"
    "errors"
    "github.com/golang-jwt/jwt/v4"
    "github.com/gin-gonic/gin"
    "time"
)

// 处理跨域请求,支持options访问
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        var (
            token  string
            myJwt  *JWT
            claims *request.CustomClaims
            err    error
        )
        token = c.Request.Header.Get("x-token")
        if token == "" {
            utils.Result(utils.ERROR, gin.H{
                "reload": true,
            }, "非法访问,请先登录!", c)
            c.Abort()
            return
        }

        myJwt = NewJWT()
        // parseToken 解析token包含的信息
        if claims, err = myJwt.ParseToken(token); err != nil {
            if err == TokenExpired {
                utils.Result(utils.ERROR, gin.H{
                    "reload": true,
                }, "授权已过期", c)
                c.Abort()
                return
            }
            utils.Result(utils.ERROR, gin.H{
                "reload": true,
            }, err.Error(), c)
            c.Abort()
            return
        }
        c.Set("claims", claims)
        c.Next()
    }
}

type JWT struct {
    SigningKey []byte
}

func NewJWT() *JWT {
    return &JWT{
        []byte(config.GConfig.JwtSigningKey),
    }
}

var (
    TokenExpired     = errors.New("token 已经过期")
    TokenNotValidYet = errors.New("token 不能验证失败")
    TokenMalformed   = errors.New("token 格式不对")
    TokenInvalid     = errors.New("token 无法解析")
)

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
    var (
        token *jwt.Token
        err   error
    )
    token, err = jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return j.SigningKey, nil
    })
    if err != nil {
        if errors.Is(err, jwt.ErrTokenMalformed) {
            return nil, TokenMalformed
        }
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, TokenExpired
        }
        if errors.Is(err, jwt.ErrTokenNotValidYet) {
            return nil, TokenNotValidYet
        }
        return nil, TokenInvalid
    }
    if token != nil {
        if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
            return claims, nil
        }
        return nil, TokenInvalid
    } else {
        return nil, TokenInvalid
    }
}

func SignToken(userModel *model.User) (int64, string, error) {
    var (
        myJwt     *JWT
        tokenSign *jwt.Token
    )
    myJwt = NewJWT()
    now := time.Now()
    expiresAt := now.Add(24 * time.Hour * 7)
    claims := request.CustomClaims{
        ID:       userModel.UserId,
        NickName: userModel.NickName,
        RegisteredClaims: jwt.RegisteredClaims{
            NotBefore: jwt.NewNumericDate(now),                          // 签名生效时间
            ExpiresAt: jwt.NewNumericDate(expiresAt),                    // 过期时间一周
            Issuer:    "chao-da-ye",                  // 签名的发行者
        },
    }

    tokenSign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err := tokenSign.SignedString(myJwt.SigningKey)
    return expiresAt.Unix() * 1000, token, err
}
