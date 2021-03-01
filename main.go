package main

import (
	"fmt"
	"github.com/whereabouts/chassis/logger"
	"github.com/whereabouts/chassis/model/mongo"
	"github.com/whereabouts/chassis/model/mongo/test_user"
	"time"
)

func main() {
	client, err := mongo.Init(mongo.Config{
		Addrs:          []string{"127.0.0.1:27017"},
		Database:       "test",
		Username:       "root",
		Password:       "root",
		Source:         "admin",
		ReplicaSetName: "",
		Timeout:        3,
		InsertTimeAuto: true,
		UpdateTimeAuto: true,
	})
	defer client.Close()
	if err != nil {
		logger.Fatalln(err)
	}
	user1 := test_user.User{Name: "hezebin1", Age: 21}
	user2 := test_user.User{Name: "hezebin2", Age: 22}
	user3 := test_user.User{Name: "hezebin3", Age: 23}
	user4 := test_user.User{Name: "hezebin4", Age: 24}
	err = test_user.GetUserDB().Add(user1, user2, user3, user4)
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = test_user.GetUserDB().DeleteByAge(24)
	if err != nil {
		logger.Errorln(err)
		return
	}
	time.Sleep(3 * time.Second)
	err = test_user.GetUserDB().ModifyAgeByName("hezebin1", 221)
	if err != nil {
		logger.Errorln(err)
		return
	}
	u, err0 := test_user.GetUserDB().GetByName("hezebin2")
	if err0 != nil {
		logger.Errorln(err0)
		return
	}
	fmt.Printf("%+v\n", u)
	err = test_user.GetUserDB().ReplaceByName("hezebin2", user4)
	if err != nil {
		logger.Errorln(err)
		return
	}
	users, err1 := test_user.GetUserDB().GetAll()
	if err1 != nil {
		logger.Errorln(err1)
		return
	}
	for _, user := range users {
		fmt.Printf("%+v\n", *user)
	}
}
