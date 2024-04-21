package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func logout(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		return
	}
	s := session{
		Sid: c.Value,
	}
	_ = DeleteSessionBySid(&s)
	c.MaxAge = -1
	http.SetCookie(w, c)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func login(w http.ResponseWriter, req *http.Request) {
	data := struct {
		u          user
		ErrorMessage string
	}{}
	nu := user{}
	if req.Method == http.MethodPost {
		nu.Email = req.FormValue("email")
		if req.FormValue("email") == "" || req.FormValue("password") == "" {
			// data.ErrorMessage = "password and email do not match"
			// tpl.ExecuteTemplate(w, "login.gohtml", data)
			http.Error(w, "Email and Password do not match", http.StatusUnauthorized)
			return
		}
		err = GetUser(&nu)
		if err != nil {
			fmt.Println(err)
		}
		// TODO encrypt password
		err = bcrypt.CompareHashAndPassword(nu.Password,[]byte(req.FormValue("password")))
		if err == nil {
			err = userLogin(w, req, nu)
			if err == nil {
				return
			} 
		}
		data.ErrorMessage = "password and email do not match"
	}
	
	data.u = nu
	tpl.ExecuteTemplate(w, "login.gohtml", data)
}

func signup(w http.ResponseWriter, req *http.Request) {
	data := struct {
		U          user
		WrongEmail bool
	}{}
	nu := user{}
	if req.Method == http.MethodPost {
		nu.Email = req.FormValue("email")
		if err != nil {
			log.Fatal("could not sign up with password")
		}

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
		nu.Username = req.FormValue("username")
		nu.Password, err = bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost) 
		err = CreateUser(nu)
		if err != nil {
			log.Fatal(err)
		}
		userLogin(w, req, nu)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}