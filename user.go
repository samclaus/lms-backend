package main

import "fmt"

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
func HandleLogin(r Request, s *Server) bool {
	var loginData Login
	loginData.Username = r.RequestData["username"]
	loginData.Password = r.RequestData["password"]
	fmt.Print("Login request received from", loginData.Username, "\n")

	var user UserInfo
	queryResult := s.Database.Take(&user, &UserInfo{Username: loginData.Username})
	if queryResult.Error != nil {
		fmt.Printf("User %v not found, closing connection\n", loginData.Username)
		return false
	}

	fmt.Printf("User %v found in database, authenticated\n", loginData.Username)
	// fmt.Print(user)
	return true
}
