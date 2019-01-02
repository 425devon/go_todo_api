package models

import "gopkg.in/mgo.v2/bson"

type Task struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Completed bool          `bson:"completed" json:"completed"`
}
