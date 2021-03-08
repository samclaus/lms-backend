package main

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
