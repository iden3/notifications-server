package db

import "github.com/globalsign/mgo"

type Mongodb struct {
	collections map[string]*mgo.Collection
}

func NewMongodb(url string, databaseName string, collectionsArray []string) (*Mongodb, error) {
	session, err := mgo.Dial("mongodb://" + url)
	if err != nil {
		return &Mongodb{}, err
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// get all the collections specified in the collectionsArray
	collections := make(map[string]*mgo.Collection)
	for _, collection := range collectionsArray {
		collections[collection] = session.DB(databaseName).C(collection)
	}
	return &Mongodb{collections}, nil
}

// GetCollections returns the current mongodb collections of the mongosrv
func (mongodb *Mongodb) GetCollections() map[string]*mgo.Collection {
	return mongodb.collections
}
