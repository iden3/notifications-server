package endpoint

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/iden3/notifications-server/db"
	"github.com/stretchr/testify/assert"
)

const dbName = "iden3-test-notifications"

func dumpDb() {
	cmd := exec.Command("mongo", dbName, "--eval", "c=db.getCollection('counters'); c.find();")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(stdout))
}

func TestCounter(t *testing.T) {
	cmd := exec.Command("mongo", dbName, "--eval", "db.dropDatabase();")
	stdout, err := cmd.Output()
	fmt.Println(string(stdout))
	//dumpDb()
	assert.Nil(t, err)
	collectionsArray := []string{"notifications", "counters"}
	mongodb, err := db.NewMongodb("127.0.0.1:27017", dbName, collectionsArray)
	assert.Nil(t, err)

	counter := NewCounter(mongodb.GetCollections()["counters"])

	for i := 0; i < 3; i++ {
		err = counter.incCounter("id5", func(n uint64) error {
			fmt.Println(n)
			return nil
		})
		assert.Nil(t, err)
	}
}
