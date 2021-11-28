package mongo

import (
	"context"

	"github.com/jucardi/go-beans/beans"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/logx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultClient = "default"
)

// The reference type to use in the `beans` dependency manager.
var ref = (*IClient)(nil)

func Dial(cfg ...*Config) (IClient, error) {
	var c *Config
	if len(cfg) > 0 && cfg[0] != nil {
		c = cfg[0]
	} else {
		c = getConfig()
	}

	client, err := mongo.Connect(context.Background(), c.opts())
	if err != nil {
		return nil, err
	}

	ret := &mgoClient{
		client: client,
		dbName: c.dbName(),
	}

	name := defaultClient
	if c.Name != "" {
		name = c.Name
	}

	current := Get(name)

	if err := beans.Register(ref, name, ret); err != nil {
		logx.WithObj(
			ret.client.Disconnect(context.Background()),
		).Error("failed to disconnect mongo session")
		return nil, err
	}

	// If overrides are allowed and another connection with the same name existed, closes that connection.
	if current != nil {
		logx.WithObj(
			current.Client().Disconnect(context.Background()),
		).Error("failed to disconnect existing mongo session")
	}

	return ret, nil
}

// Get attempts to retrieve an open connection by the given name. Returns nil if no connection by the given
// name is found
//
//   {name} - (Optional) The unique name for the connection to be retrieved. If not provided, returns the
//            primary connection.
//
func Get(name ...string) IClient {
	n := stringx.GetOrDefault(defaultClient, name...)
	if beans.Exists(ref, n) {
		return beans.Resolve(ref, n).(IClient)
	}
	return nil
}
