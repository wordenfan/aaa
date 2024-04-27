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
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkmid "github.com/rookie-ninja/rk-entry/v2/middleware"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	rkginctx "github.com/rookie-ninja/rk-gin/v2/middleware/context"
	_ "github.com/rookie-ninja/rk-grpc/v2/boot"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	rkzero "github.com/rookie-ninja/rk-zero/boot"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"tk-boot-worden/api/gen/v1"
	"tk-boot-worden/router"
	"tk-boot-worden/tools"
)

var MySecret = []byte("my-secret")

// JWT claims contains UID
type CustomClaims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

var userDb *gorm.DB
var redisClient *redis.Client
var logger *rkentry.LoggerEntry

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
	_ = os.Setenv("DEV_REGION", "qingdao")
	boot := rkboot.NewBoot()

	// Logger
	logger = rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	logger.Info("This is my-logger")

	// Grpc register  (bootstrap前)
	entryRpc := rkgrpc.GetGrpcEntry("go-grpc")
	entryRpc.AddRegFuncGrpc(registerGreeter)
	entryRpc.AddRegFuncGw(grt.RegisterGreeterHandlerFromEndpoint)

	// Zero register  (bootstrap前)
	zeroEntry := rkzero.GetZeroEntry("go-zero")
	zeroEntry.Server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/v1/zero_demo",
		Handler: ZeroGreeter,
	})

	// Router Group
	entryGin := rkgin.GetGinEntry("go-gin")
	//entryGin.AddMiddleware(RouterMiddle())
	//redisGroup := entryGin.Router.Group("v3")
	//{
	//	redisGroup.GET("/demo_api", demoRequest)
	//}

	// Error
	rkmid.SetErrorBuilder(&tools.MyErrorBuilder{})

	// 路由
	router.InitRouter(entryGin.Router)

	// JWT
	entryGin.Router.GET("/v1/jwt_token", JwtToken)
	entryGin.Router.GET("/v1/login", Login)

	// Mysql
	mysqlEntry := rkmysql.GetMySqlEntry("user-db")
	if mysqlEntry != nil {
		userDb = mysqlEntry.GetDB("rk-boot")
		if !userDb.DryRun {
			_ = userDb.AutoMigrate(&User{})
		}
	}
	entryGin.Router.GET("/v1/user/:id", GetUser)
	entryGin.Router.PUT("/v1/user", CreateUser)

	// Config
	fmt.Println(rkentry.GlobalAppCtx.GetConfigEntry("my-config").GetString("region"))

	// Run
	boot.Bootstrap(context.TODO())
	boot.WaitForShutdownSig(context.TODO())
}

// ================================================
// Greeter handler
// @Summary Greeter
// @Id 1
// @Tags Hello
// @version 1.0
// @Param name query string true "name"
// @produce application/json
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func ZeroGreeter(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	resp := &ZeroGreeterResponse{
		Message: fmt.Sprintf("Hello %s!", request.URL.Query().Get("name")),
	}
	bytes, _ := json.Marshal(resp)
	writer.Write(bytes)
}

type ZeroGreeterResponse struct {
	Message string
}

// ================================================
func demoRequest(ctx *gin.Context) {
	fmt.Println(" 我的测试 API! ")
	logger.Info("我的测试 API!")
}

// ================================================
func RouterMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("路由分组中间件-before")
		logger.Info("路由分组中间件-before")
		// 可以在这里添加任何预处理逻辑，比如验证token、记录日志等
		// ...
		// 然后一定要调用c.Next()来传递给下一个处理器
		c.Next()
		fmt.Println("路由分组中间件-after")
	}
}

// ================================================
type GreeterServer struct{}

func registerGreeter(server *grpc.Server) {
	grt.RegisterGreeterServer(server, &GreeterServer{})
}

func (server *GreeterServer) Hello(_ context.Context, _ *grt.HelloRequest) (*grt.HelloResponse, error) {
	return &grt.HelloResponse{
		MyMessage: "hello!",
	}, nil
}

func (server *GreeterServer) Person(_ context.Context, req *grt.PersonRequest) (*grt.PersonResponse, error) {
	p := &grt.PersonResponse{
		Id:    req.GetId(),
		Name:  "worden",
		Email: "rs@example.com",
		Phones: []*grt.PersonResponse_PhoneNumber{
			{Number: "555-4321", Type: grt.PersonResponse_HOME},
		},
	}
	return p, nil
}

// ================================================
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
// @Router /v1/jwt_token [get]
func JwtToken(ctx *gin.Context) {
	jwtToken := rkginctx.GetJwtToken(ctx)
	ctx.JSON(http.StatusOK, map[string]string{
		"Message": fmt.Sprintf("Hello %s!", GetPhoneFromJwtToken(jwtToken)),
	})
}

// ================================================
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
