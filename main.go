package main

import (
	"database/sql"
	"fmt"
	"os"

	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var tpl *template.Template
var udb map[string]user

type user struct {
	Email string
	Password []byte
}



func init(){
	// get templates
	tpl = template.Must(template.ParseGlob("templates/*"))
	// connect to db
    DBept := os.Getenv("DB_ENDPOINT")
	DBpass := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	dataSourceName := fmt.Sprint("admin:%s@tcp(%s)/%s?charset=utf8", DBept, DBpass, DBname)
	db, err := sql.Open("mysql", dataSourceName)
	CheckError(err, "Could not open DB")
	defer db.Close()
	err = db.Ping()
	CheckError(err, "DB does not respond")
}

func main (){
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, req *http.Request){
 	var nu user
	data := struct{
		User user
		Login bool
	 }{
		nu,
		false,
	 }
	if req.Method == http.MethodPost {
		data.User.Email = req.FormValue("email")
		data.Login = false
		password, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			log.Fatal(err)
		}

		data.User.Password = password
		u, ok := udb[data.User.Email];
		if ok {
			if err := bcrypt.CompareHashAndPassword(u.Password, data.User.Password); err == nil {
				tpl.ExecuteTemplate(w, "index.gohtml", data)
				return
			}
		}
		data.Login = true
	}
	
	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func index(w http.ResponseWriter, req *http.Request){
   tpl.ExecuteTemplate(w, "index.gohtml", nil)
}