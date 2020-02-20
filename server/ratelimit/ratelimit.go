package ratelimit

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Source: https://github.com/s12i/gin-throttle/blob/master/throttle.go
// But as I don't want to add a dependency on such a small thing,
// I just copied it here.
func Throttle(maxEventsPerSec int, maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(context *gin.Context) {
		if limiter.Allow() {
			context.Next()
			return
		}

		context.Error(errors.New("Limit exceeded"))
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}
