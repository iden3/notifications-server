module github.com/iden3/notifications-server

require (
	github.com/ethereum/go-ethereum v1.8.23
	github.com/gin-contrib/cors v0.0.0-20190226021855-50921afdc5c1
	github.com/gin-gonic/gin v1.3.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/iden3/go-iden3 v0.0.3-0.20190325144113-c072f0edf3fd
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/viper v1.3.1
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
)

// replace github.com/iden3/go-iden3 => ../go-iden3
