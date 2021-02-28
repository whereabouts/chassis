package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/whereabouts/chassis/logger"
	"github.com/whereabouts/chassis/model/mongo"
	"github.com/whereabouts/chassis/model/mongo/test_user"
)

func main() {
	//session, err := mgo.Dial("mongodb://root:root@127.0.0.1:27017")
	//defer session.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c := session.DB("test").C("student")
	////data := test_user.User{
	////	Name:   "a11",
	////	//Age:    2111,
	////}
	////dataM := map[string]interface{} {
	////	"name": "b1",
	////	"age": 12,
	////}
	//err = c.Update(bson.M{"name":"aaa"}, map[string]interface{}{"age":222})
	//if err != nil {
	//	log.Println(err)
	//}
	client, err := mongo.Init(mongo.Options{
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
	//u := &test_user.User{
	//	Id: bson.NewObjectId(),
	//	Name: "hezebin2",
	//	Age: 100,
	//}
	err = test_user.GetUserDB().Do(func(c *mgo.Collection) error {
		return c.Update(bson.M{"name": "hezebin2"}, bson.M{"$set": bson.M{"age": 100}})
	})
	if err != nil {
		fmt.Println("err: ", err)
	}
	//user1 := test_user.User{Name: "hezebin1", Age: 21}
	//user2 := test_user.User{Name: "hezebin2", Age: 22}
	//user3 := test_user.User{Name: "hezebin3", Age: 23}
	//user4 := test_user.User{Name: "hezebin4", Age: 24}
	//err = test_user.GetUserDB().Add(user1, user2, user3, user4)
	//if err != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//err = test_user.GetUserDB().DeleteByAge(21)
	//if err != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//time.Sleep(3 * time.Second)
	//err = test_user.GetUserDB().ModifyAgeByName("hezebin2", 32)
	//if err != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//users, err1 := test_user.GetUserDB().GetAll()
	//if err1 != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//fmt.Printf("%+v", users)
}
