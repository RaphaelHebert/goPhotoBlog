package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

type user struct {
	Email string
	Password []byte
}

func init(){
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main (){
	http.HandleFunc("/", login)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, req *http.Request){
 	var nu user
	if req.Method == http.MethodPost {
		nu.Email = req.FormValue("email")
		password, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			log.Fatal(err)
		}
		nu.Password = password
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nu)
}