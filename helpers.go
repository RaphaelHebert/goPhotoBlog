package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CheckError(e error){
	if e != nil {
		log.Fatal(e)
	}
}

// IsAlreadyLoggedIn if the user making the request is logged in .
func IsAlreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	s := session{
		Sid: c.Value,
	}
	c.MaxAge = expireTime
	http.SetCookie(w, c)
	err = GetSessionBySid(&s)
	if err != nil {
		return false
	}
	// TODO: check if lastUpdate is older than session lifeSpan
 	return true
}

// makeSessionCookie create a session cookie.
func makeSessionCookie () *http.Cookie {
	sID := uuid.NewString()
	c := &http.Cookie{
		Name: "session",
		Value: sID,
	}
	return c
}

func userLogin(w http.ResponseWriter, req *http.Request, u user) error {
	// create session
	c := makeSessionCookie()
	c.MaxAge = expireTime
	http.SetCookie(w, c)

	s := session{
		Email: u.Email,
	}
    err := GetSessionByEmail(&s)
	if err == nil {
		err = DeleteSession(&s)
		if err != nil {
			return errors.New("could not log you in")
		}
	}
	s.Sid = c.Value
	s.LastUpdated = time.Now().UTC()
	err = CreateSession(&s)

	if err != nil {
		// TODO handle error to redirect to sign up with message for the user without ending main process
		return errors.New("could not log you in")
	}
	http.Redirect(w, req, "/index", http.StatusSeeOther)
	return nil
}