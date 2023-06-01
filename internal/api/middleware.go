package api

import (
	"errors"
	"rating-api/internal/util/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware
// Logs HTTP requests with a predefined structure.
func LoggingMiddleware(loggr logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		host := c.Request.Host
		route := c.FullPath()
		remoteAddr := c.Request.RemoteAddr
		clientIp := c.ClientIP()
		protocol := c.Request.Proto
		method := c.Request.Method
		uri := c.Request.RequestURI
		queryString := c.Request.URL.RawQuery
		elapsedMilliseconds := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		hasError := len(c.Errors.Errors()) > 0

		var errs []error
		if hasError {
			for i := 0; i < len(c.Errors.Errors()); i++ {
				errs = append(errs, errors.New(c.Errors.Errors()[i]))
			}
		}

		pingRoute := "/api/ping"
		if route != pingRoute {
			logMessage := protocol + " " + method + " " + uri + " responded " + strconv.Itoa(statusCode) + " in " + strconv.Itoa(int(elapsedMilliseconds)) + " ms"

			if hasError {
				loggr.Error(
					logMessage,
					zap.String("host", host),
					zap.String("route", route),
					zap.String("protocol", protocol),
					zap.String("uri", uri),
					zap.String("method", method),
					zap.String("remoteAddr", remoteAddr),
					zap.String("clientIp", clientIp),
					zap.String("queryString", queryString),
					zap.Int("statusCode", statusCode),
					zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
					zap.Errors("errors", errs),
				)
			} else {
				loggr.Info(
					logMessage,
					zap.String("host", host),
					zap.String("route", route),
					zap.String("protocol", protocol),
					zap.String("uri", uri),
					zap.String("method", method),
					zap.String("remoteAddr", remoteAddr),
					zap.String("clientIp", clientIp),
					zap.String("queryString", queryString),
					zap.Int("statusCode", statusCode),
					zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
					zap.Errors("errors", errs),
				)
			}
		}
	}
}
