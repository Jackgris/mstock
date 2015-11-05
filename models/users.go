package models

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection name
const table string = "users"

// Data representing the application user, is the structure that is
// going to relate all application data, to perform effective monitoring
// of user actions.
type User struct {
	IdUser    string
	Name      string
	Pass      string
	LastLogin time.Time
	CreatedAt time.Time
	UpdateAt  time.Time
}

// Updated user data, and if it does not exist create a new one
func (u User) Save() error {
	session, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(dbname).C(table)
	update := bson.M{"$set": u}
	info, err := c.UpsertId(u.IdUser, update)
	log.Println("Save user", info)
	return err
}

// We will return the user data associated with ID
func (u User) Get() (User, error) {
	session, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(dbname).C(table)

	result := User{}
	err = c.FindId(u.IdUser).One(&result)
	if err != nil {
		log.Fatalln("Get user", err)
		return User{}, err
	}
	return result, err
}

// It will remove the user associated with that ID
func (u User) Delete() error {
	session, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(dbname).C(table)

	return c.Remove(bson.M{"iduser": u.IdUser})
}
