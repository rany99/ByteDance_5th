package middle

import (
	"ByteDance_5th/models"
	"crypto/sha1"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

var JwtKey = []byte("imwave")

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// 生成Token
func GenerateToken(user models.User) (string, error) {
	//token有效期设置为一周
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.UserInfoId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "ByteDance_5th",
			Subject:   "imwave",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstr, err := token.SignedString(JwtKey)
	if err != nil {
		log.Println("token生成失败", tokenstr)
		return "", err
	}
	return tokenstr, nil
}

// 解析token
func DecodeToken(tokenstr string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenstr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}

// 中间件
func Permission() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenstr := context.Query("token")
		if tokenstr == "" {
			tokenstr = context.PostForm("token")
		}
		//缺少用户信息
		if tokenstr == "" {
			context.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 401,
				StatusMsg:  "用户未认证",
			})
			context.Abort()
			return
		}
		//验证token
		token_decoded, ok := DecodeToken(tokenstr)
		if !ok {
			context.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 403,
				StatusMsg:  "token解析错误",
			})
			context.Abort()
			return
		}
		if time.Now().Unix() > token_decoded.ExpiresAt {
			context.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 402,
				StatusMsg:  "token已过期",
			})
			context.Abort()
			return
		}
		context.Set("user_id", token_decoded.UserId)
		context.Next()
	}
}

func ShaMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		password := context.Query("password")
		if password == "" {
			password = context.PostForm("password")
		}
		context.Set("password", Sha1(password))
		context.Next()
	}
}

func NoAuthToGetUserId() gin.HandlerFunc {
	return func(context *gin.Context) {
		rawId := context.Query("user_id")
		if rawId == "" {
			rawId = context.PostForm("user_id")
		}
		//用户未认证
		if rawId == "" {
			context.JSON(http.StatusOK, models.CommonResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			context.Abort() //阻止执行
			return
		}
		userid, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			context.JSON(http.StatusOK, models.CommonResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			context.Abort()
		}
		context.Set("user_id", userid)
		context.Next()
	}
}

func Sha1(str string) string {
	o := sha1.New()
	o.Write([]byte(str))
	return hex.EncodeToString(o.Sum(nil))
}
