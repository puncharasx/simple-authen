package main

import (
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB

const jwtSecret = "secretkeyjwt"

type User struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Reponse struct {
	Token string `json:"token"`
}

func main() {
	var err error
	db, err = sqlx.Open("mysql", "user:pass@/dbname")
	checkError(err)
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(200)
		rw.Write([]byte(string("Hello World")))

	})

	router.HandleFunc("/login", login).Methods(http.MethodPost)
	router.HandleFunc("/register", register).Methods(http.MethodPost)

	http.ListenAndServe(":8080", router)

}

func login(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("request body incorrect format"))
		return
	}
	request := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("request body not json"))
	}
	user := User{}
	query := "SELECT id,username,password FROM users WHERE username=?"
	err = db.Get(&user, query, request.Username)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Username or Password Incorrect"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Username or Password Incorrect"))
		return
	}

	cliams := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Issuer:    user.Username,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	token, err := jwtToken.SignedString([]byte(jwtSecret))

	if err != nil {
		rw.Write([]byte("Error JWT"))
	}
	reponse := Reponse{
		Token: token,
	}

	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(&reponse)

}

func register(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("request body incorrect format"))
		return
	}
	request := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("request body not json"))
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	checkError(err)
	query := "INSERT INTO users (username,password) VALUES (?,?)"
	result, err := db.Exec(
		query,
		request.Username,
		string(password),
	)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cannot register"))
	}
	id, err := result.LastInsertId()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cannot register "))
	}

	user := User{
		Id:       int(id),
		Username: request.Username,
		Password: string(password),
	}

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	json.NewEncoder(rw).Encode(user)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
