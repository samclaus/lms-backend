package main

import "crypto/sha256"

func updateUsername123(s *Server) {
	var user UserInfo
	s.Database.Take(&user, &UserInfo{Username: "username123"})

	var saltedPassword []byte = append(user.Salt, "password123"...)
	var passwordHash [32]byte = sha256.Sum256(saltedPassword)
	numRounds := 0

	for i := 1; i < int(user.Rounds); i++ {
		passwordHash = sha256.Sum256(passwordHash[0:])
		numRounds++
	}

	user.PasswordHash = passwordHash[0:]
	s.Database.Save(&user)
}
