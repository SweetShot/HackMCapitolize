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
	"os/exec"
	"strconv"
	

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

func ExecuteBlockchain(command string) string{
	log.Println("exec blockchain ")
	cmd := exec.Command("truffle", "console")
	cmd.Dir = "/home/om/truffletest"
	//
	command += ".then(r => console.log(r.toString()))"
	
	//log.Println(command)
	stdin,err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, command + "\n")
	}()
	out, err := cmd.Output()
	if err != nil {
		//log.Println("part1")
		log.Fatal(err)
		
	}
	//log.Println("chain out ", string(out))
	strout := strings.Replace(string(out), "truffle(development)> ", "", -1)
	strout = strings.Replace(strout, "undefined", "", -1)
	//log.Println(strout)
	strout = string(([]rune(strout))[0:len(strout)-2])
	//log.Println(strout)
	return strout
}

/* Middleware handler to handle all requests for authentication */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		req.Header.Set("Username", "") // Clean if any
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {

				token, err := jwt.ParseWithClaims(bearerToken[1], &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(SECRET), nil
				})

				if err != nil {
					json.NewEncoder(w).Encode(Exception{Message: err.Error()})
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
			ExpiresAt: time.Now().Add(time.Hour * 100).Unix(),
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

// GET /Ideas
func (c *Controller) GetIdeas(w http.ResponseWriter, r *http.Request) {
	log.Println(ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.getIdeaCount.call())"))

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
	//log.Println(vars)

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
	//log.Println(vars)

	username := vars["name"] // param id
	//log.Println(username)

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
		return
	}

	// Update Idea fields
	idea.Username = r.Header.Get("Username")
	idea.DatePosted = time.Now().Unix()
	idea.TotalFundsRaised = 0
	idea.NumContributors = 0

	// Add it to DB
	log.Println("Adding Idea")
	ideaId := idea.Username + strconv.FormatInt(idea.DatePosted, 10)
	data,_ := json.Marshal(idea)
	data2 := string(data)
	data2 = strings.Replace(data2, "\"", "\\\"", -1)
	out := ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.addIdea.sendTransaction(\"" + 
		ideaId + "\",\"" + data2 + "\"))")

	log.Println(out)

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

	ideaId := idea.Username + strconv.FormatInt(idea.DatePosted, 10)
	data,_ := json.Marshal(idea)
	data2 := string(data)
	data2 = strings.Replace(data2, "\"", "\\\"", -1)
	out := ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.updateIdea.sendTransaction(\"" +
		 ideaId + "\",\"" + data2 + "\"))");

	log.Println(out)

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

	ideaId := idea.Username + strconv.FormatInt(idea.DatePosted, 10)
	out := ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.removeIdea.sendTransaction(\"" +
		ideaId + "\"))");

    log.Println(out)

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


// AddTransation POST /
func (c *Controller) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error AddTransaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddTransaction", err)
	}

	if err := json.Unmarshal(body, &transaction); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddTransaction unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Update username
	transaction.Username = r.Header.Get("Username")

	// Retrive idea with this id
	//transaction.DatePosted = time.Now().Unix()//hax
	idea := c.Repository.GetUniqueIdea(transaction.Username, transaction.DatePosted)
	transaction.DatePosted = time.Now().Unix()

	log.Println("Adding Transaction for Idea ", idea.Title, " price ", transaction.Price)

	out := ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.addTransaction.sendTransaction(" + 
		strconv.FormatInt(transaction.Price, 10) + ",\"" + transaction.IdeaId + "\",\"" + transaction.Username +  
		"\",\"" + transaction.OtherInfo + "\"," + strconv.FormatInt(transaction.DatePosted, 10) + "))");

	log.Println(out)
	idea.TotalFundsRaised += transaction.Price;
	
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

// AddTransation GET / DONT use next
func (c *Controller) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions Transactions

	out := ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.getTransactionCount.call())");
	log.Println(out);
	num, _ := strconv.Atoi(out)
	log.Println(out);

	for i := 0; i < num; i++ {
		out = ExecuteBlockchain("IdeaContract.deployed().then(instance => instance.getTransaction.call("+ 
			strconv.Itoa(i) +"))")
		log.Println(out)
		split := strings.Split(out, ",")
		prc, _ := strconv.Atoi(split[0])
		dpt, _ := strconv.Atoi(split[4])
		transactions = append(transactions, Transaction{
			IdeaId : split[1],
			DatePosted : int64(dpt),
			Username : split[2],
			Price : int64(prc),
			OtherInfo : split[3],
		});
	}

	

	/*if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} */
	data, _ := json.Marshal(transactions)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
