package user

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	rkredis "github.com/rookie-ninja/rk-db/redis"
	"net/http"
	"time"
)

type HandlerUser struct{}

func (*HandlerUser) getcaptcha(ctx *gin.Context) {
	ctx.JSON(200, "getCaptcha success")
}

// ==============================================================
var redisClient *redis.Client
var redisEntry *rkredis.RedisEntry

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func init() {
	redisEntry = rkredis.GetRedisEntry("redis")
	if redisEntry != nil {
		redisClient, _ = redisEntry.GetClient()
	}
}

func (*HandlerUser) GetRedis(ctx *gin.Context) {
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
func (*HandlerUser) SetRedis(ctx *gin.Context) {
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
