package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// UserInfo represents user authentication information stored in the sqlite3 database.
// Salt and rounds are randomly generated per user when a new user is created.
// The password is NEVER stored in the database, only the hashed version.
type UserInfo struct {
	ID           string `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"username" json:"username"`
	Salt         []byte `gorm:"salt" json:"-"`
	Rounds       uint   `gorm:"rounds" json:"-"`
	PasswordHash []byte `gorm:"password_hash" json:"-"`
}

// Login information sent by the client, including the username and the password
// before it has been salted and hashed.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleLogin is the login request handler function. It takes a Request
// and a pointer to the GORM database object as a parameter.
// TODO: replace success bools and message strings into an error or outcome type
func HandleLogin(r WebSocketRequest, s *Server) (success bool, authenticatedUser UserInfo, message string) {

	var loginData Login
	loginData.Username = r.RequestData["username"]
	loginData.Password = r.RequestData["password"]
	fmt.Printf("Login request received from %v\n", loginData.Username)

	var user UserInfo
	user.Username = loginData.Username
	queryResult := s.Database.Take(&user, &user)
	if queryResult.Error != nil {
		fmt.Printf("User %v not found, closing connection\n", loginData.Username)
		return false, user, "Authentication Failed: username not found"
	}

	fmt.Printf("User %v found in database, attempting to authenticate\n", loginData.Username)
	var saltedPassword []byte = append(user.Salt, loginData.Password...)
	var passwordHash [32]byte = sha256.Sum256(saltedPassword)

	for i := 1; i < int(user.Rounds); i++ {
		passwordHash = sha256.Sum256(passwordHash[0:])
	}

	storedHash := user.PasswordHash
	fmt.Printf("Hashed Password: %v\n", passwordHash)
	fmt.Printf("Stored Hash:     %v\n", storedHash)

	if bytes.Equal(passwordHash[0:], storedHash[0:]) {
		fmt.Printf("User %v authenticated successfully!", user.Username)
		return true, user, "Authentication Successful! Welcome, " + user.Username
	}

	fmt.Printf("Incorrect password for user %v, closing connection", user.Username)
	return false, user, "Authentication Failed: incorrect password"
}
