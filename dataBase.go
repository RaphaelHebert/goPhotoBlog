package main

import (
	"errors"
	"fmt"
)

// getUser takes a pointer to a user object as argument.
// its query the db for the corresponding user
func GetUser(u *user) error {
	query := fmt.Sprintf("SELECT * FROM users.test where email = '%s'", u.Email)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		// if more than one user returned check your DB setting as email should be unique
		if count > 0 {
			return errors.New("non unique email")
		}
		err := rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		CheckError(err)
		count++
	}
	if count == 0 {
		return errors.New("no match")
	}
	return nil
}