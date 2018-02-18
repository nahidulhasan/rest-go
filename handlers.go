package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"strconv"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/go-sql-driver/mysql"
)


//var db *sql.DB
var mdb *sql.DB


//dbConn is create a connection for local database
func dbConn() (db *sql.DB) {
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_DATABASE")
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	//dbDriver := "mysql"


	//db, err := sql.Open(dbDriver, user+":"+pass+"@"+"("+host+":"+port+")/"+database)


	cfg := mysql.Config{
		User:   user,
		Passwd: pass,
		DBName: database,
		Addr:   host + ":" + port,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	//db, err := sql.Open("mysql", user+":"+pass+"@"+scheme+"("+host+":"+port+")/"+database)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}


func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request)  {

	fmt.Fprintln(w, "Welcome in GO todo app!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	qs := fmt.Sprintf("SELECT id, name, completed, due FROM task")

	db := dbConn()
	defer db.Close()

	selDB, err := db.Query(qs)
	if err != nil {
		panic(err)
	}

	todos := []Todo{}

	for selDB.Next() {
		var todo Todo
		err = selDB.Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.Due)
		if err != nil {
			panic(err)
		}

		todos = append(todos, todo)
	}


	if err := json.NewEncoder(w).Encode(todos); err != nil {

		panic(err)
	}

}

func TodoShow( w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)

}

func TodoCreate(w http.ResponseWriter, r *http.Request)  {

	var todo Todo

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {

		panic(err)
	}


	if err := r.Body.Close(); err != nil {

		panic(err)
	}


	if err := json.Unmarshal(body, &todo); err != nil {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(422)


		if err := json.NewEncoder(w).Encode(todos); err != nil {

			panic(err)
		}

	}

	//t := RepoCreateTodo(todo)

	RepoCreateTodo(todo)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusCreated)


	if err := json.NewEncoder(w).Encode(todos); err != nil {

		panic(err)
	}

}


func TodoDelete( w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	todoId := vars["todoId"]

	id, err := strconv.Atoi(todoId)

	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	RepoDestroyTodo(id)


}



