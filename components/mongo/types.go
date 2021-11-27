package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type IClient interface {
	// DB clones the existing mongo session and returns the requested DB instance
	DB(name ...string) (mongo.Session, *mongo.Database, error)
	// Client returns the underlying mongo client contained by this wrapper
	Client() *mongo.Client
}
