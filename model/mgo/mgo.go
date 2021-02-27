package mgo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func New(database string, collection string, value interface{}) *MongoDB {
	return &MongoDB{database: database, collection: collection, value: value}
}

func (db *MongoDB) Database() string {
	return db.database
}

func (db *MongoDB) Collection() string {
	return db.collection
}

func (db *MongoDB) client() Client {
	return getGlobalClient()
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

func (db *MongoDB) Insert(doc ...interface{}) error {
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

func (db *MongoDB) Upsert(selector, update interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	err = db.Do(func(c *mgo.Collection) error {
		changeInfo, err = c.Upsert(selector, update)
		return err
	})
	return changeInfo, err
}

func (db *MongoDB) upsertId(id, update interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	err = db.Do(func(c *mgo.Collection) error {
		changeInfo, err = c.UpsertId(id, update)
		return err
	})
	return changeInfo, err
}

func (db *MongoDB) UpdateAll(selector, update interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	err = db.Do(func(c *mgo.Collection) error {
		v := reflect.ValueOf(update)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		now := time.Now()
		if v.Kind() == reflect.Struct {
			m := make(map[string]interface{})
			m, err = mapper.Struct2Map(update)
			if err != nil {
				return err
			}
			m["update_time"] = timer.Format(now)
			changeInfo, err = c.UpdateAll(selector, m)
		} else if v.Kind() == reflect.Map {
			v.SetMapIndex(reflect.ValueOf("update_time"), reflect.ValueOf(timer.Format(now)))
			changeInfo, err = c.UpdateAll(selector, v.Interface())
		}
		changeInfo, err = c.UpdateAll(selector, update)
		return err
	})
	return changeInfo, err
}

func (db *MongoDB) Update(selector, update interface{}) error {
	return db.Do(func(c *mgo.Collection) (err error) {
		v := reflect.ValueOf(update)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		now := time.Now()
		if v.Kind() == reflect.Struct {
			m := make(map[string]interface{})
			m, err = mapper.Struct2Map(update)
			if err != nil {
				return err
			}
			if !mapper.IsKeyTimeValid(m, "update_time") {
				m["update_time"] = timer.Format(now)
			}
			err = c.Update(selector, m)
		} else if v.Kind() == reflect.Map {
			if !v.MapIndex(reflect.ValueOf("update_time")).IsValid() {
				v.SetMapIndex(reflect.ValueOf("update_time"), reflect.ValueOf(timer.Format(now)))
			}
			err = c.Update(selector, v.Interface())
		}
		return err
	})
}

func (db *MongoDB) UpdateId(id, update interface{}) (err error) {
	err = db.Do(func(c *mgo.Collection) error {
		err = c.UpdateId(id, update)
		return err
	})
	return err
}

func (db *MongoDB) FindId(id interface{}, picker interface{}, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		query := c.FindId(id)
		if picker != nil {
			query = query.Select(picker)
		}
		err := query.One(ret)
		return err
	})
}

func (db *MongoDB) FindObjectId(id string, picker interface{}, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		_id := bson.ObjectIdHex(id)
		query := c.FindId(_id)
		if picker != nil {
			query = query.Select(picker)
		}
		err := query.One(ret)
		return err
	})
}

func (db *MongoDB) FindOne(selector interface{}, picker interface{}, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		query := c.Find(selector)
		if selector != nil {
			query = query.Select(picker)
		}
		err := query.One(ret)
		return err
	})
}

func (db *MongoDB) FindOneWithSort(selector interface{}, sort []string, picker interface{}, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		query := c.Find(selector)
		if sort != nil {
			query = query.Sort(sort...)
		}
		if selector != nil {
			query = query.Select(picker)
		}
		err := query.One(ret)
		return err
	})
}

func (db *MongoDB) Count(selector interface{}) (count int, err error) {
	err = db.Do(func(c *mgo.Collection) error {
		count, err = c.Find(selector).Count()
		return err
	})
	return count, err
}

func (db *MongoDB) FindAll(selector interface{}, sort []string, picker interface{}, skip int, limit int, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		query := c.Find(selector)
		if selector != nil {
			query = query.Select(picker)
		}
		if sort != nil {
			query = query.Sort(sort...)
		}
		if skip > 0 {
			query.Skip(skip)
		}
		if limit > 0 {
			query.Limit(limit)
		}
		return query.All(ret)
	})
}

func (db *MongoDB) PipeAll(selector interface{}, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		return c.Pipe(selector).All(ret)
	})
}

func (db *MongoDB) Distinct(selector interface{}, key string, ret interface{}) error {
	return db.Do(func(c *mgo.Collection) error {
		return c.Find(selector).Distinct(key, ret)
	})
}
