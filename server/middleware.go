package server

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"pod-service-relations/dao"
	"pod-service-relations/database"
	"time"

	"github.com/gin-gonic/gin"

	"pod-service-relations/logging"
)

func logMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logging.GetLogger().Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func loginMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		basicAuth, err := extractBasicAuth(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "message": err.Error()})
			c.Abort()
			return
		}
		ud := dao.NewUserDao(database.GetDB())
		if _, err := ud.Get(0, basicAuth.User, basicAuth.Token); err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "message": http.StatusText(http.StatusUnauthorized)})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
			c.Abort()
			return
		}
		logging.GetLogger().WithFields(logrus.Fields{
			"user": basicAuth.User,
		}).Infof("login")
		c.Next()
	}
}
