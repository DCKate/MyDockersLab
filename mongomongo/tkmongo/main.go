package main

import (
	"fmt"
	"tkmongo/mongo"
)

type Person struct {
	Name    string "bson:`name`"
	Phone   string "bson:`phone`"
	Created string "bson:`created`"
	Age     int    "bson:`age`"
}

func main() {
	se, db := mongo.ConnectDB("192.168.33.10:27017", "test", "kk", "kk_passwd")
	defer mongo.Disconnect(se, db)

	for ii := 0; ii < 10000; ii++ {
		name := fmt.Sprintf("Ale[%v]", ii)
		pp := Person{name, "+55 53 8116 9639", fmt.Sprintf("%v", ii), ii % 100}
		db.C("people").Insert(&pp) //, &Person{"Cla", "+55 53 8402 8510", 30})
		fmt.Printf("%v\n", pp)
	}

	mm := make(map[string]interface{}, 2)
	mm["age"] = 15
	rr := mongo.FindDocuments(db.C("people"), mm)
	fmt.Printf("%v\n", rr)
}
