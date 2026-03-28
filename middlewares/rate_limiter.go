package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gblog/utils"
)

var globalStore = memory.NewStore()
var strictStore = memory.NewStore()

var globalLimiter = limiter.Rate{
	Period: time.Minute,
	Limit:  200,
}

var strictLimiter = limiter.Rate{
	Period: time.Minute,
	Limit:  5,
}

func GlobalRateLimiter() gin.HandlerFunc {
	instance := limiter.New(globalStore, globalLimiter)
	return func(c *gin.Context) {
		key := c.ClientIP()
		context, err := instance.Get(c.Request.Context(), key)
		if err != nil {
			utils.Logger.WithError(err).Error("Rate limiter error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if context.Reached {
			utils.Logger.WithField("ip", key).Warn("Global rate limit exceeded")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

func StrictRateLimiter() gin.HandlerFunc {
	instance := limiter.New(strictStore, strictLimiter)
	return func(c *gin.Context) {
		key := c.ClientIP()
		context, err := instance.Get(c.Request.Context(), key)
		if err != nil {
			utils.Logger.WithError(err).Error("Rate limiter error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if context.Reached {
			utils.Logger.WithField("ip", key).Warn("Strict rate limit exceeded")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}
