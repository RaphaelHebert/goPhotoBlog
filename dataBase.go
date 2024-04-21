package main

import (
	"errors"
	"fmt"
)

// getUser takes a pointer to a user object as argument.
// its query the db for the corresponding user
func GetUser(u *user) error {
	query := fmt.Sprintf("SELECT * FROM users.users where email = '%s'", u.Email)
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
		err := rows.Scan(&u.Id, &u.Email, &u.Password, &u.Username)
		CheckError(err)
		count++
	}
	if count == 0 {
		return errors.New("no match")
	}
	return nil
}

func CreateUser(u user) error {
	stmt, err := db.Prepare("INSERT INTO users (email, username, password) VALUES (?, ?, ?)")
	CheckError(err)
	defer stmt.Close()

	r, err := stmt.Exec(u.Email, u.Username, u.Password)
	CheckError(err)
	n, err := r.RowsAffected()

	fmt.Printf("users INSERT %v rows\n", n)
	return err
}

func GetSessionBySid(s *session) error {
	query := fmt.Sprintf("SELECT * FROM users.session where sid = '%s'", s.Sid)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()


	count := 0
	for rows.Next() {
		// if more than one session returned check your DB setting as session should be unique
		if count > 0 {
			return errors.New("non unique email")
		}
		err := rows.Scan(&s.Sid, &s.Email, &s.LastUpdated)
		CheckError(err)
		count++
	}
	if count == 0 {
		return errors.New("no match")
	}
	return nil
}

func GetSessionByEmail(s *session) error {
	query := fmt.Sprintf("SELECT * FROM users.session where email = '%s'", s.Email)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()


	count := 0
	for rows.Next() {
		// if more than one session returned check your DB setting as session should be unique
		if count > 0 {
			return errors.New("non unique email")
		}
		err := rows.Scan(&s.Sid, &s.Email, &s.LastUpdated)
		CheckError(err)
		count++
	}
	if count == 0 {
		return errors.New("no match")
	}
	return nil
}

func DeleteSession(s *session) error {
	query := fmt.Sprintf("DELETE FROM session WHERE email='%s';", s.Email)

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec()
	CheckError(err)

	n, err := r.RowsAffected()
	CheckError(err)

	fmt.Printf("session DELETED %v rows\n", n)

	return nil
}

func DeleteSessionBySid(s *session) error {
	query := fmt.Sprintf("DELETE FROM session WHERE sid='%s';", s.Sid)

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec()
	CheckError(err)

	n, err := r.RowsAffected()
	CheckError(err)

	fmt.Printf("session DELETED %v rows\n", n)

	return nil
}

func CreateSession(s *session) error {
	stmt, err := db.Prepare("INSERT INTO session (sid, email, lastUpdated) VALUES (?, ?, ?)")
	CheckError(err)
	defer stmt.Close()

	r, err := stmt.Exec(s.Sid, s.Email, s.LastUpdated)
	CheckError(err)
	n, err := r.RowsAffected()

	fmt.Printf("session INSERT %v rows\n", n)
	return err
}

func UpdateSession(s *session) error {
	query := fmt.Sprintf(`UPDATE session SET lastUpdated="%v" WHERE sid="%s";`, s.LastUpdated, s.Sid)
	stmt, err := db.Prepare(query)
	CheckError(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	CheckError(err)

	n, err := r.RowsAffected()
	CheckError(err)

	fmt.Printf("session UPDATED %v rows\n", n)
	return err
}