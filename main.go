package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rookie-ninja/rk-boot/v2"
	rkmysql "github.com/rookie-ninja/rk-db/mysql"
	rkredis "github.com/rookie-ninja/rk-db/redis"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	rkginctx "github.com/rookie-ninja/rk-gin/v2/middleware/context"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var MySecret = []byte("my-secret")

// JWT claims contains UID
type CustomClaims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

var userDb *gorm.DB
var redisClient *redis.Client

type Base struct {
	CreatedAt time.Time      `yaml:"-" json:"-"`
	UpdatedAt time.Time      `yaml:"-" json:"-"`
	DeletedAt gorm.DeletedAt `yaml:"-" json:"-" gorm:"index"`
}

type User struct {
	Base
	Id   int    `yaml:"id" json:"id" gorm:"primaryKey"`
	Name string `yaml:"name" json:"name"`
}

func main() {
	//worden_test.Test_HookFunc("aa")
	//worden_test.Test_b("bb")
	//return
	//
	boot := rkboot.NewBoot()
	boot.Bootstrap(context.TODO())
	entry := rkgin.GetGinEntry("greeter")
	entry.Router.GET("/v1/get", Get)
	entry.Router.POST("/v1/set", Set)

	redisEntry := rkredis.GetRedisEntry("redis")
	redisClient, _ = redisEntry.GetClient()

	//JWT
	entry.Router.GET("/v1/greeter", Greeter)
	entry.Router.GET("/v1/login", Login)
	//
	//_ = os.Setenv("DOMAIN", "dev")
	mysqlEntry := rkmysql.GetMySqlEntry("user-db")
	userDb = mysqlEntry.GetDB("rk-boot")
	if !userDb.DryRun {
		userDb.AutoMigrate(&User{})
	}
	entry.Router.GET("/v1/user/:id", GetUser)
	entry.Router.PUT("/v1/user", CreateUser)

	boot.WaitForShutdownSig(context.TODO())
}
func GetUser(ctx *gin.Context) {
	uid := ctx.Param("id")
	user := &User{}
	res := userDb.Where("id = ?", uid).Find(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}

	if res.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context) {
	user := &User{
		Name: ctx.Query("name"),
	}

	res := userDb.Create(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Greeter service
// @Id 1
// @version 1.0
// @produce application/json
// @Param name query string true "Input name"
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	jwtToken := rkginctx.GetJwtToken(ctx)
	ctx.JSON(http.StatusOK, map[string]string{
		"Message": fmt.Sprintf("Hello %s!", GetPhoneFromJwtToken(jwtToken)),
	})
}

// Login API
func Login(ctx *gin.Context) {
	token, _ := GenerateAccessToken()
	ctx.JSON(http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

func GenerateAccessToken() (tokenString string, err error) {
	claim := CustomClaims{
		Phone: "18561122236",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

func GetPhoneFromJwtToken(jwtToken *jwt.Token) string {
	claims := &CustomClaims{}
	bytes, _ := json.Marshal(jwtToken.Claims)
	_ = json.Unmarshal(bytes, claims)

	return claims.Phone
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Set(ctx *gin.Context) {
	payload := &KV{}

	if err := ctx.BindJSON(payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	cmd := redisClient.Set(ctx.Request.Context(), payload.Key, payload.Value, time.Minute)

	if cmd.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, cmd.Err())
		return
	}

	ctx.Status(http.StatusOK)
}

func Get(ctx *gin.Context) {
	key := ctx.Query("key")

	cmd := redisClient.Get(ctx.Request.Context(), key)

	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			ctx.JSON(http.StatusNotFound, "Key not found!")
		} else {
			ctx.JSON(http.StatusInternalServerError, cmd.Err())
		}
		return
	}

	payload := &KV{
		Key:   key,
		Value: cmd.Val(),
	}

	ctx.JSON(http.StatusOK, payload)
}
