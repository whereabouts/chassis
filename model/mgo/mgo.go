package mgo

import (
	"context"
	"github.com/globalsign/mgo"
	"github.com/whereabouts/utils/mapper"
	"github.com/whereabouts/utils/timer"
	"reflect"
	"time"
)

type MongoDB struct {
	database   string
	collection string
	value      interface{}
}

func (db *MongoDB) Database() string {
	return db.database
}

func (db *MongoDB) Collection() string {
	return db.collection
}

func (db *MongoDB) client() Client {
	return defaultClient()
}

func (db *MongoDB) Do(f func(c *mgo.Collection) error) error {
	return db.client().Do(db, f)
}

func (db *MongoDB) Remove(selector interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		return c.Remove(selector)
	})
}

func (db *MongoDB) RemoveID(id interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		return c.RemoveId(id)
	})
}

func (db *MongoDB) RemoveAll(selector interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	err = db.Do(func(c *mgo.Collection) error {
		changeInfo, err = c.RemoveAll(selector)
		return err
	})
	return changeInfo, err
}

func (db *MongoDB) Insert(ctx context.Context, doc ...interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		out := make([]interface{}, 0, len(doc))
		for _, in := range doc {
			v := reflect.ValueOf(in)
			for v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			now := time.Now()
			if v.Kind() == reflect.Struct {
				m, err := mapper.Struct2Map(in)
				if err != nil {
					return err
				}
				if !mapper.IsKeyTimeValid(m, "update_time") {
					m["update_time"] = timer.Format(now)
				}
				if !mapper.IsKeyTimeValid(m, "create_time") {
					m["create_time"] = timer.Format(now)
				}
				out = append(out, m)
			} else if v.Kind() == reflect.Map {
				if !v.MapIndex(reflect.ValueOf("update_time")).IsValid() {
					v.SetMapIndex(reflect.ValueOf("update_time"), reflect.ValueOf(timer.Format(now)))
				}
				if !v.MapIndex(reflect.ValueOf("create_time")).IsValid() {
					v.SetMapIndex(reflect.ValueOf("create_time"), reflect.ValueOf(timer.Format(now)))
				}
				out = append(out, v.Interface())
			}
		}
		return c.Insert(out...)
	})
}
