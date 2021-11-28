package mongo

import (
	"github.com/jucardi/go-strings/stringx"
	"go.mongodb.org/mongo-driver/mongo"
)

type mgoClient struct {
	client *mongo.Client
	dbName string
}

func (c *mgoClient) Client() *mongo.Client {
	return c.client
}

func (c *mgoClient) DB(name ...string) (mongo.Session, *mongo.Database, error) {
	session, err := c.client.StartSession()
	if err != nil {
		return nil, nil, err
	}
	dbName := stringx.GetOrDefault(c.dbName, name...)
	return session, session.Client().Database(dbName), nil
}
