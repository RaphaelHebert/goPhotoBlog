package main

import (
	"net/http"
)

// GuestGuard prevents logged in users to access routes reserved to users that are not logged.
func guestGuard(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !IsAlreadyLoggedIn(w, req) {
			//http.Error(w, "not logged in", http.StatusUnauthorized)
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return // don't call original handler
		}
		h.ServeHTTP(w, req)
	})
}

// AuthGuard prevents logged in users to access routes reserved to users that are not logged.
func authGuard(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if IsAlreadyLoggedIn(w, req) {
			//http.Error(w, "not logged in", http.StatusUnauthorized)
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return // don't call original handler
		}
		h.ServeHTTP(w, req)
	})
}