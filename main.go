package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template
var err error

// ExpireTime sets the lifespan of a session.
const expireTime int = 600 // session lifespan in sec

type user struct {
	Id int
	Email string
	Password []byte
	Username string
}

type session struct {
	Sid string
	Email string 
	LastUpdated time.Time
}

func init(){
	// get templates
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	// connect to db
	DBept := os.Getenv("DB_ENDPOINT") 
	DBpass := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	dataSourceName := fmt.Sprintf("admin:%s@tcp(%s)/%s?parseTime=true&charset=utf8", DBpass, DBept, DBname)
	db, err = sql.Open("mysql", dataSourceName)
	CheckError(err)
	//defer db.Close()
	err = db.Ping()
	CheckError(err)
}

func main (){
	http.HandleFunc("/", guestGuard(index))
	http.HandleFunc("/login", authGuard(login))
	http.HandleFunc("/logout", guestGuard(logout))
	http.HandleFunc("/signup", authGuard(signup))
	http.ListenAndServe(":8080", nil)
}