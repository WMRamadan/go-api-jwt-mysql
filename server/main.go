package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var mySigningKey = []byte("mysecret")

type User struct {
	Name string `json:"Name"`
}

type Users []User

func getUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/go-api-jwt-mysql")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT name FROM users")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var user User

		err = results.Scan(&user.Name)
		if err != nil {
			panic(err.Error())
		}
		users := Users{
			User{Name: user.Name},
		}
		json.NewEncoder(w).Encode(users)
	}
	fmt.Println("Endpoint Hit: All Users")
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	fid := 0
	fname := r.Form.Get("name")

	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/go-api-jwt-mysql")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insert, err := db.Prepare("INSERT INTO users (id, name) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	insert.Exec(fid, fname)
	defer insert.Close()
	fmt.Fprintf(w, "POST Endpoint")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/", isAuthorized(homePage))
	myRouter.Handle("/users", isAuthorized(getUsers)).Methods("GET")
	myRouter.HandleFunc("/users", postUsers).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main() {
	handleRequests()
}
