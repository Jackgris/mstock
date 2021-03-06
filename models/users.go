package models

import (
	"errors"
	"log"
	"time"

	"github.com/astaxie/beego/validation"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection name
const table string = "users"

// Data representing the application user, is the structure that is
// going to relate all application data, to perform effective monitoring
// of user actions.
type User struct {
	IdUser    string    `json:"id_user"`
	Name      string    `json:"name"`
	Pass      string    `json:"pass"`
	Email     string    `json:"email"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update"`
	Token     Token     `json:"token"`
}

// We will check if the user data are valid
func (u User) Valid() bool {
	v := validation.Validation{}
	v.Required(u.Name, "name")
	v.MaxSize(u.Name, 20, "nameMax")
	v.Required(u.Pass, "pass")
	v.MaxSize(u.Pass, 30, "passMax")
	v.Email(u.Email, "email")

	if v.HasErrors() {
		for _, e := range v.Errors {
			log.Println("Check valid user data:", e)
		}
		return false
	}

	return true
}

// Updated user data, and if it does not exist create a new one
func (u User) Save() error {
	session, err := mgo.Dial(DB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if u.Name == "" {
		return errors.New("User name can't be empty for save on database")
	}
	if u.Email == "" {
		return errors.New("User email can't be empty for save on database")
	}
	if u.Pass == "" {
		return errors.New("User pass can't be empty for save on database")
	}

	c := session.DB(DB_NAME).C(table)
	update := bson.M{"$set": u}
	_, err = c.UpsertId(u.IdUser, update)
	return err
}

// We will return the user data associated with ID
func (u User) Get() (User, error) {
	session, err := mgo.Dial(DB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(DB_NAME).C(table)

	result := User{}
	err = c.FindId(u.IdUser).One(&result)
	if err != nil {
		log.Println("Get user", err)
		return User{}, err
	}
	return result, err
}

// It will remove the user associated with that ID
func (u User) Delete() error {
	session, err := mgo.Dial(DB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(DB_NAME).C(table)

	return c.Remove(bson.M{"iduser": u.IdUser})
}

// We are going to search for users in the database,
// given through an e-mail
func FindUsers(email string) ([]User, error) {
	session, err := mgo.Dial(DB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(DB_NAME).C(table)

	query := bson.M{"email": email}
	result := []User{}
	err = c.Find(query).All(&result)
	if err != nil {
		log.Println("Find user", err)
		return []User{}, err
	}

	return result, err
}
