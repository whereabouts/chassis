package test_user

import (
	"github.com/globalsign/mgo/bson"
	"github.com/whereabouts/chassis/model/mgo"
)

var (
	userDB = &UserDB{mgo.New("test", "user", User{})}
)

func GetUserDB() *UserDB {
	return userDB
}

type UserDB struct {
	*mgo.MongoDB
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (user *UserDB) Add(users ...interface{}) error {
	return user.Insert(users...)
}

func (user *UserDB) AddOne(u User) error {
	return user.Insert(u)
}

func (user *UserDB) GetByName(name string) (*User, error) {
	u := &User{}
	selector := bson.M{"name": name}
	err := user.FindOne(selector, nil, u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (user *UserDB) GetAll() ([]*User, error) {
	ret := make([]*User, 0)
	err := user.FindAll(nil, nil, nil, 0, 0, &ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (user *UserDB) ModifyAgeByName(name string, age int) error {
	selector := bson.M{"name": name}
	return user.Update(selector, User{Name: name, Age: age})
}

func (user *UserDB) DeleteByAge(age int) error {
	selector := bson.M{"age": age}
	return user.Remove(selector)
}
