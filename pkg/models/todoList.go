package models

import "gopkg.in/mgo.v2/bson"

type TodoList struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Tasks       []Task        `bson:"tasks" json:"tasks"`
}
