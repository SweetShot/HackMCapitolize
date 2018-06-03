package store

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "mongodb://test:test123@ds141720.mlab.com:41720/capitolize"

// DBNAME the name of the DB instance
const DBNAME = "capitolize"

// COLLECTION is the name of the collection in DB
const IdeaCollection = "store"

const LoginCollection = "users"

var IdeaId = 10

func (r Repository) Register() bool {
	return false
}

func (r Repository) Login(user User) User {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()

	c := session.DB(DBNAME).C(LoginCollection)
	result := User{}

	if err := c.Find(bson.M{"username": user.Username}).One(&result); err != nil {
		fmt.Println("Failed to write results:", err)
	}
	return result
}

// GetIdeas returns the list of Ideas
func (r Repository) GetIdeas() Ideas {

	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(IdeaCollection)
	results := Ideas{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to read results:", err)
	}
	return results
}

func (r Repository) GetIdeasByUsername(Username string) Ideas {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(IdeaCollection)
	results := Ideas{}

	if err := c.Find(bson.M{"username": Username}).All(&results); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return results
}

// GetIdeasByString takes a search string as input and returns Ideas
func (r Repository) GetIdeasByString(query string) Ideas {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(IdeaCollection)
	result := Ideas{}

	// Logic to create filter
	qs := strings.Split(query, " ")
	and := make([]bson.M, len(qs))
	for i, q := range qs {
		and[i] = bson.M{"title": bson.M{
			"$regex": bson.RegEx{Pattern: ".*" + q + ".*", Options: "i"},
		}}
	}
	filter := bson.M{"$and": and}

	if err := c.Find(&filter).Limit(10).All(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// GetIdea
func (r Repository) GetUniqueIdea(username string, date_posted int64) Idea {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	idea := Idea{}
	log.Println("Getting unique Idea : ", username, date_posted)
	err = session.DB(DBNAME).C(IdeaCollection).Find(
		bson.M{"$and": []bson.M{
			{"username": username},
			{"date_posted": date_posted},
		}}).One(&idea)

	if err != nil {
		log.Fatal(err)
		return Idea{Title : "Error"}
	}

	log.Println("Retrieved Idea - ", idea.Username, idea.DatePosted)

	return idea
}


// AddIdea adds an Idea in the DB
func (r Repository) AddIdea(idea Idea) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB(DBNAME).C(IdeaCollection).Insert(idea)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Idea ID- ", idea.Title)
	return true
}

// UpdateIdea updates a Idea in the DB
func (r Repository) UpdateIdea(idea Idea) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	log.Println("Updating Idea : ", idea.Username, idea.DatePosted)
	err = session.DB(DBNAME).C(IdeaCollection).Update(
		bson.M{"$and": []bson.M{
			{"username": idea.Username},
			{"date_posted": idea.DatePosted},
		}}, idea)

	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Println("Updated Idea - ", idea.Username, idea.DatePosted)

	return true
}

// DeleteIdea Deletes an Idea
func (r Repository) DeleteIdea(idea Idea) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// Remove Idea
	if err = session.DB(DBNAME).C(IdeaCollection).Remove(bson.M{"username": idea.Username, "date_posted": idea.DatePosted}); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	fmt.Println("Deleted Idea- ", idea.Title)
	// Write status
	return "OK"
}
