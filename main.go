package main

type Student struct {
	Name string
	Age  int
}

func main() {
	//session, err := mgo.Dial("mongodb://root:root@127.0.0.1:27017")
	//defer session.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c := session.DB("test").C("student")
	//data := Student{
	//	Name:   "a2",
	//	Age:    22,
	//}
	////dataM := map[string]interface{} {
	////	"name": "b1",
	////	"age": 12,
	////}
	//err = c.Insert(&data)
	//if err != nil {
	//	log.Println(err)
	//}
}
