package commands

import (
	"fmt"

	"github.com/iden3/notifications-server/config"
	"github.com/iden3/notifications-server/db"
	log "github.com/sirupsen/logrus"
)

func LoadMongoService() *db.Mongodb {
	collectionsArray := []string{"notifications", "counters"}
	mongoservice, err := db.NewMongodb(config.C.Mongodb.Url, config.C.Mongodb.Database, collectionsArray)
	if err != nil {
		fmt.Println("Cannot open mongodb storage")
		panic(err)
	}
	log.WithField("path", config.C.Mongodb.Url).Info("Mongodb storage opened")
	return mongoservice
}
