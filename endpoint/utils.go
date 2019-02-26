package endpoint

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type Counter struct {
	mutex      sync.Mutex
	collection *mgo.Collection
}

func NewCounter(collection *mgo.Collection) *Counter {
	return &Counter{collection: collection}
}

func (c *Counter) incCounter(key string, f func(n uint64) error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var n struct{ N uint64 }
	if err := c.collection.FindId(key).One(&n); err != nil {
		if err.Error() == "not found" {
			if err := c.collection.Insert(bson.M{"_id": key, "n": 0}); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	n.N++
	if err := c.collection.UpdateId(key, bson.M{"n": n.N}); err != nil {
		return err
	}

	return f(n.N)
}

func fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}
