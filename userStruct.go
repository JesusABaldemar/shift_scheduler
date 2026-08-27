package main

import (
	"fmt"
)

type User struct {
	Username string
	Password string
	Schedule Schedule
	admin	bool
}

func submitSchedule(user *User, schedule Schedule) {
	user.Schedule = schedule
}


func createUser(username string, password string) *User {
	newUser := User{
		Username: username,
		Password: password,
		Schedule: *createSchedule(),
		admin:    false,
	}
	return &newUser
}

func queryUserToSQL(user User) string {
	query := "INSERT INTO users (username, password, schedule, admin) VALUES "
	query += fmt.Sprintf("('%s', '%s', '%s', %t)", user.Username, user.Password, toJson(&user.Schedule), user.admin)
	query += ";"
	return query
}