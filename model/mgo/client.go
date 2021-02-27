package mgo

import (
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"time"
)

type Client interface {
	GetSession() *mgo.Session
	Close()
	Do(model Model, exec func(s *mgo.Collection) error) error
}

var globalClient Client

func Init(option Options) (Client, error) {
	c, err := NewClient(option)
	if err != nil {
		return nil, err
	}
	globalClient = c
	return c, err
}

func InitFast(url string) (Client, error) {
	c, err := Dial(url)
	if err != nil {
		return nil, err
	}
	globalClient = c
	return c, err
}

func getGlobalClient() Client {
	return globalClient
}

// example: mongodb://myuser:mypass@localhost:27017,otherhost:27017/db
func Dial(url string) (Client, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect mongodb")
	}
	return &client{session: session}, nil
}

func NewClient(option Options) (Client, error) {
	poolLimit := defaultPoolLimit
	if option.PoolLimit != 0 {
		poolLimit = option.PoolLimit
	}

	maxIdleTime := defaultMaxIdleTime
	if option.MaxIdleTime != time.Duration(0) {
		maxIdleTime = option.MaxIdleTime * time.Second
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          option.Addrs,
		Database:       option.Database,
		Username:       option.Username,
		Password:       option.Password,
		Source:         option.Source,
		ReplicaSetName: option.ReplicaSetName,
		Timeout:        option.Timeout * time.Second,
		PoolLimit:      poolLimit,
		MaxIdleTimeMS:  int(maxIdleTime / time.Millisecond),
		AppName:        option.AppName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect mongodb")
	}

	if option.Mode < mgo.Primary {
		option.Mode = mgo.PrimaryPreferred
	}
	session.SetMode(option.Mode, true)

	return &client{session: session}, nil
}

type client struct {
	session *mgo.Session
}

func (c *client) GetSession() *mgo.Session {
	return c.session.Copy()
}

func (c *client) Close() {
	c.session.Close()
}

func (c *client) Do(model Model, exec func(s *mgo.Collection) error) error {
	s := c.GetSession()
	defer s.Close()
	return exec(s.DB(model.Database()).C(model.Collection()))
}
