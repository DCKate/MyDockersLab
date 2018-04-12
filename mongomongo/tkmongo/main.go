package main

import (
	"fmt"
	"tkmongo/mongo"
)

type Person struct {
	Name  string "bson:`name`"
	Phone string "bson:`phone`"
	Age   int    "bson:`age`"
}

func main() {
	se, db := mongo.ConnectDB("localhost:27017", "testlog", "kk", "pass")
	defer mongo.Disconnect(se, db)
	mm := make(map[string]interface{}, 2)
	db.C("people").Insert(&Person{"Ale", "+55 53 8116 9639", 15}, &Person{"Cla", "+55 53 8402 8510", 30})
	mm["phone"] = "Ale"
	mm["age"] = 15
	rr := mongo.FindDocuments(db.C("people"), mm)
	fmt.Printf("%v\n", rr)
}
