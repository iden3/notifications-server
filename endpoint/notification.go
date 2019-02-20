package endpoint

import (
	"github.com/gin-gonic/gin"
)

func handleGetInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "baz",
	})
}

func handlePostNotification(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "baz",
	})
}

func handleGetNotifications(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "baz",
	})
}

func handleDeleteNotifications(c *gin.Context) {
	c.JSON(200, gin.H{
		"foo": "baz",
	})
}
