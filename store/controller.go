package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

//Controller ...
type Controller struct {
	Repository Repository
}

const SECRET = "1234aswedf1234"

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

/* Middleware handler to handle all requests for authentication */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		req.Header.Set("Username", "") // Clean if any
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(SECRET), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					claims := token.Claims.(*MyClaims)
					log.Println("TOKEN WAS VALID ", claims.Username)
					req.Header.Set("Username", claims.Username)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)

	// Validate user info
	userInfo := c.Repository.Login(user)
	log.Println(userInfo.Username + " " + userInfo.Password)
	if userInfo.Password != user.Password {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid password"})
		return
	}
	// Else if proper user
	claims := MyClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//log.Println("Username: " + user.Username)
	log.Println("Log In: " + user.Username)

	tokenString, error := token.SignedString([]byte(SECRET))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// GET /Ideas
func (c *Controller) GetIdeas(w http.ResponseWriter, r *http.Request) {
	ideas := c.Repository.GetIdeas() // list of all ideas
	log.Println("Ideas Request")
	data, _ := json.Marshal(ideas)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// SearchIdea GET /Ideas/Search/{query}
func (c *Controller) GetIdeasByString(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)

	query := vars["query"] // param query
	log.Println("Search Query - " + query)

	ideas := c.Repository.GetIdeasByString(query)
	data, _ := json.Marshal(ideas)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// /Idea/User/{name}
func (c *Controller) GetIdeasByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)

	username := vars["name"] // param id
	log.Println(username)

	ideas := c.Repository.GetIdeasByUsername(username)
	data, _ := json.Marshal(ideas)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddIdea POST /
func (c *Controller) AddIdea(w http.ResponseWriter, r *http.Request) {
	var idea Idea
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	//log.Println(body)

	if err != nil {
		log.Fatalln("Error AddIdea", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddIdea", err)
	}

	if err := json.Unmarshal(body, &idea); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddIdea unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Update Idea fields
	idea.Username = r.Header.Get("Username")
	idea.DatePosted = time.Now().Unix()
	idea.TotalFundsRaised = 0
	idea.NumContributors = 0

	// Add it to DB
	log.Println("Adding Idea")
	success := c.Repository.AddIdea(idea) // adds the idea to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// UpdateIdea PUT /
func (c *Controller) UpdateIdea(w http.ResponseWriter, r *http.Request) {
	var idea Idea
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error UpdateIdea", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdateIdea", err)
	}

	if err := json.Unmarshal(body, &idea); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateIdea unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Update username
	idea.Username = r.Header.Get("Username")

	log.Println("Updating Idea", idea.Title)
	success := c.Repository.UpdateIdea(idea) // updates the Idea in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// DeleteIdea DELETE /
func (c *Controller) DeleteIdea(w http.ResponseWriter, r *http.Request) {
	var idea Idea
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error DeleteIdea", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error DeleteIdea", err)
	}

	if err := json.Unmarshal(body, &idea); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateIdea unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Update username
	idea.Username = r.Header.Get("Username")

	if err := c.Repository.DeleteIdea(idea); err != "" { // delete
		log.Println(err)
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}
