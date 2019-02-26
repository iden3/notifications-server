package endpoint

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func handleGetInfo(c *gin.Context) {
	c.JSON(200, gin.H{})
}

type NotificationMsg struct {
	Data string `json:"data"`
}

type Notification struct {
	Id     uint64 `json:"id"`
	Date   int64  `json:"date"`
	Data   string `json:"data"`
	ToAddr string `json:"toAddr"`
}

func handlePostNotification(c *gin.Context) {
	var notificationMsg NotificationMsg
	c.BindJSON(&notificationMsg)

	idAddr := c.Param("idaddr")

	if err := counter.incCounter(idAddr, func(n uint64) error {
		notification := Notification{
			Id:     n,
			Date:   time.Now().Unix(),
			Data:   notificationMsg.Data,
			ToAddr: idAddr,
		}
		return mongodb.GetCollections()["notifications"].Insert(notification)
	}); err != nil {
		fail(c, "error on handleGetNotifications", err)
		return
	}
	c.JSON(200, gin.H{})
}

type GetNotificationMsg struct {
	IdAddr string `json:"idAddr"`
	// SignedPacket string `json:"signedPacket"`
	// ProofKSign ProofClaim `json:"proofKSign"`
}

func handleGetNotifications(c *gin.Context) {
	var m GetNotificationMsg
	c.BindJSON(&m)

	// TODO check signature of requester id

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
		"toaddr": m.IdAddr,
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
	var m GetNotificationMsg
	c.BindJSON(&m)

	// TODO check signature of requester id

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
		"toaddr": m.IdAddr,
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