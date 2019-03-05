package endpoint

import (
	"context"
	"crypto/rand"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/iden3/notifications-server/config"
	"github.com/iden3/notifications-server/db"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/middleware/iden-assert-auth"
)

var mongodb db.Mongodb
var counter *Counter
var nonceDb *core.NonceDb

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server {
	// start serviceapi
	api := gin.Default()
	api.Use(cors.Default())

	serviceapi := api.Group("/api/unstable")
	serviceapi.GET("/", handleGetInfo)

	var key [256 / 8]byte
	if _, err := rand.Read(key[:]); err != nil {
		panic(err)
	}
	authapi, err := auth.AddAuthMiddleware(serviceapi, config.C.Server.Domain, nonceDb, key[:])
	if err != nil {
		panic(err)
	}

	serviceapi.POST("/notifications/:idaddr", handlePostNotification)
	authapi.GET("/notifications", handleGetNotifications)
	authapi.DELETE("/notifications", handleDeleteNotifications)

	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: api}
	go func() {
		log.Info("API server at ", config.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func Serve(mgodb *db.Mongodb) {

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	mongodb = *mgodb
	counter = NewCounter(mongodb.GetCollections()["counters"])
	nonceDb = core.NewNonceDb()
	// start servers
	serviceapisrv := serveServiceApi()

	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}
}
