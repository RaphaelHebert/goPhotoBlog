package main

import (
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

type user struct {
	Email string
	Password string
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
		fmt.Print(req.FormValue(`email`))
		nu.Email = req.FormValue("email")
		nu.Password = bcrypt.GenerateFromPassword([]byte(req.FormValue("password")))
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nu)
}