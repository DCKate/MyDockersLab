package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ConnectDB(host string, dbname string, user string, pass string) (*mgo.Session, *mgo.Database) {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	// defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbname)
	if len(user) > 0 && len(pass) > 0 {
		err = db.Login(user, pass)
		if err != nil {
			panic(err)
		}
	}
	return session, db
}
func Disconnect(ss *mgo.Session, db *mgo.Database) {
	db.Logout()
	ss.Close()
}

func FindDocuments(c *mgo.Collection, mm map[string]interface{}) []interface{} {
	bb := bson.M{}
	for kk, vv := range mm {
		bb[kk] = vv
	}
	var rr []interface{}
	var tmp interface{}
	log.Printf("%v\n", bb)
	iter := c.Find(bb).Iter() //.One(&tmp)
	//c.Find(bson.M{"toys": bson.M{"$exists": true}})

	// Select enables selecting which fields should be retrieved for the results found. For example,
	// the following query would only retrieve the name field:
	// err := collection.Find(nil).Select(bson.M{"name": 1}).One(&result)

	//err := c.Find(nil).Sort("-age").Skip(2).Limit(2).All(&users)

	for iter.Next(&tmp) {
		rr = append(rr, tmp)
		// log.Printf("%v\n", tmp)
	}
	return rr
}
