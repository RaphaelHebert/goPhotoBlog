package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template
var err error

type user struct {
	Id int
	Email string
	Password string
	Username string
}

func init(){
	// get templates
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	// connect to db
	DBept := os.Getenv("DB_ENDPOINT") 
	DBpass := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	dataSourceName := fmt.Sprintf("admin:%s@tcp(%s)/%s?charset=utf8", DBpass, DBept, DBname)
	db, err = sql.Open("mysql", dataSourceName)
	CheckError(err)
	//defer db.Close()
	err = db.Ping()
	CheckError(err)
}

func main (){
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, req *http.Request){
	data := struct{
		u user
		WrongLogin bool
	}{}
	nu := user{}
	if req.Method == http.MethodPost {
		nu.Email = req.FormValue("email")
		if req.FormValue("email") == "" || req.FormValue("password") == "" {
			http.Error(w, "Email and Password do not match", http.StatusUnauthorized)
			return
		}
		err = GetUser(&nu)
		if err != nil {
			fmt.Println(err)
		}
		itTrue := req.FormValue("password") == nu.Password // bcrypt.CompareHashAndPassword(nu.Password,[]byte(req.FormValue("password")))
		nu.Password = ""
		nu.Id = 0
		if itTrue {
			// TODO: handle with cookie and session management
			tpl.ExecuteTemplate(w, "index.gohtml", nu)
			return
		}
		data.WrongLogin = true
	}
	data.u = nu
	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func signup(w http.ResponseWriter, req *http.Request){
	data := struct{
		U user
		WrongEmail bool
	}{}
	nu := user{}
	if req.Method == http.MethodPost {
		nu.Email = req.FormValue("email")
		nu.Password = req.FormValue("password")
		nu.Username = req.FormValue("username")

		if req.FormValue("email") == "" || req.FormValue("password") == "" || req.FormValue("username") == "" {
			http.Error(w, "You must fill up all the fields", http.StatusUnauthorized)
			return
		}
		err = GetUser(&nu)
		if err == nil {
			data.U = nu
			data.WrongEmail = true
			tpl.ExecuteTemplate(w, "signup.gohtml", data)
			return
		}
		// save in db
		err = NewUser(nu)
		if err != nil {
			log.Fatal(err)
		}
		// TODO create session 
		http.Redirect(w, req, "/index", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func index(w http.ResponseWriter, req *http.Request){
   tpl.ExecuteTemplate(w, "index.gohtml", nil)
}