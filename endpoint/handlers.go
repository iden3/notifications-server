package endpoint

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"github.com/iden3/go-iden3/core"
	auth "github.com/iden3/go-iden3/middleware/iden-assert-auth"
)

func handleGetInfo(c *gin.Context) {
	c.JSON(200, gin.H{})
}

type NotificationMsg struct {
	Jws string `json:"jws" binding:"required"`
}

type Notification struct {
	Id   uint64  `json:"id"`
	Date int64   `json:"date"`
	Jws  string  `json:"jws"`
	To   core.ID `json:"to"`
}

func handlePostNotification(c *gin.Context) {
	var notificationMsg NotificationMsg
	c.BindJSON(&notificationMsg)

	id, err := core.IDFromString(c.Param("id"))
	if err != nil {
		fail(c, "bad id", err)
		return
	}

	if err := counter.incCounter(id[:], func(n uint64) error {
		notification := Notification{
			Id:   n,
			Date: time.Now().Unix(),
			Jws:  notificationMsg.Jws,
			To:   id,
		}
		return mongodb.GetCollections()["notifications"].Insert(notification)
	}); err != nil {
		fail(c, "error on handlePostNotifications", err)
		return
	}
	c.JSON(200, gin.H{})
}

func handleGetNotifications(c *gin.Context) {
	user := auth.GetUser(c)

	afterid, err := strconv.ParseInt(c.DefaultQuery("afterid", "0"), 10, 64)
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}
	beforeid, err := strconv.ParseInt(c.DefaultQuery("beforeid", "9999999999"), 10, 64)
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}

	var notifications []Notification
	err = mongodb.GetCollections()["notifications"].Find(bson.M{
		"to": user.Id,
		"id": bson.M{
			"$gt": afterid,
			"$lt": beforeid,
		},
	}).Sort("-id").Limit(10).All(&notifications)
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}

	c.JSON(200, gin.H{
		"notifications": notifications,
	})
}

func handleDeleteNotifications(c *gin.Context) {
	user := auth.GetUser(c)

	afterid, err := strconv.Atoi(c.DefaultQuery("afterid", "0"))
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}
	beforeid, err := strconv.Atoi(c.DefaultQuery("beforeid", "9999999999"))
	if err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}

	info, err := mongodb.GetCollections()["notifications"].RemoveAll(bson.M{
		"to": user.Id,
		"id": bson.M{
			"$gte": afterid,
			"$lte": beforeid,
		},
	})
	if err != nil {
		fail(c, "error on handleDeleteNotifications", err)
		return
	}
	c.JSON(200, gin.H{
		"status":  "notifications deleted",
		"removed": info.Removed,
	})
}
