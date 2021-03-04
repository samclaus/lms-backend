package main

type UserInfo struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `gorm:"username" json:"username"`
	Salt     []byte `gorm:"salt" json:"-"`
	Rounds   uint   `gorm:"rounds" json:"-"`
	Password []byte `gorm:"password" json:"-"`
}
