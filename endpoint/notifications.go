package endpoint

import (
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func handleGetInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func handlePostNotification(c *gin.Context) {
	var notificationMsg NotificationMsg
	c.BindJSON(&notificationMsg)

	idAddr := c.Param("idaddr")

	notification := Notification{
		Date:   uint64(time.Now().Unix()),
		Data:   notificationMsg.Data,
		ToAddr: idAddr,
	}

	err := mongodb.GetCollections()["notifications"].Insert(notification)
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}
	c.JSON(200, gin.H{
		"foo": "baz",
	})
}

func handleGetNotifications(c *gin.Context) {
	var m GetNotificationMsg
	c.BindJSON(&m)

	// TODO url parameters filter
	// ?beforeid={notificationid}
	// ?afterid={notificationid}

	var notifications []Notification
	err := mongodb.GetCollections()["notifications"].Find(bson.M{"toaddr": m.IdAddr}).Sort("-$natural").All(&notifications)
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}

	c.JSON(200, gin.H{
		"notifications": notifications,
	})
}

func handleDeleteNotifications(c *gin.Context) {
	var m GetNotificationMsg
	c.BindJSON(&m)

	info, err := mongodb.GetCollections()["notifications"].RemoveAll(bson.M{"toaddr": m.IdAddr})
	if err != nil {
		fail(c, "error on handleDeleteNotifications", err)
		return
	}
	c.JSON(200, gin.H{
		"status":  "notifications deleted",
		"removed": info.Removed,
	})
}
