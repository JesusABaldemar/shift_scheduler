package main

import (
	//	"fmt"
//	"encoding/json"
	"net/http"
	//	"strconv"
	"html/template"
	"database/sql"
	_ "modernc.org/sqlite"
	"fmt"
	"log"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, signIn)
	
}

func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sign In Request Received")
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	Users, err := sql.Open("sqlite", "./data/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer Users.Close()

	user := createUser(r.FormValue("username"), r.FormValue("password"))

	query := queryUserToSQL(*user)
	_, err = Users.Exec(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("User signed in successfully")
	
	cookie := &http.Cookie{
		Name:  "username",
		Value: r.FormValue("username"),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/schedule", http.StatusSeeOther)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	Users, err := sql.Open("sqlite", "./data/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer Users.Close()

	query := "SELECT username, password, schedule, admin FROM users WHERE username = ? AND password = ?"

	var username string

	var password string

	var schedule string

	var admin bool

	err = Users.QueryRow(query, r.FormValue("username"), r.FormValue("password"),).Scan(&username, &password, &schedule, &admin)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("User logged in successfully")

	cookie := &http.Cookie{
		Name:  "username",
		Value: r.FormValue("username"),
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/schedule", http.StatusSeeOther)
}

func schedule(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	username := cookie.Value

	users, err := sql.Open("sqlite", "./data/users.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer users.Close()

	var scheduleData string

	err = users.QueryRow(
		"SELECT schedule FROM users WHERE username = ?",
		username,
	).Scan(&scheduleData)

	if err != nil {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	schedule := fromJson(scheduleData)

	var table string

	for i := 0; i < 60; i++ {
		table += "<tr>"
		table += fmt.Sprintf("<td>%d</td>", i+1)

		if schedule.Monday[i] {
			table += fmt.Sprintf(`
			<td
				class="available"
				hx-post = "/select-start"
				hx-vals='{"day":"Monday","slot":"%d"}'
				hx-target="#schedule-container"
			></td>
			`, i)
		} else {
			table += fmt.Sprintf(`
			<td
				hx-post="/select-start"
				hx-vals='{"day":"Monday","slot":"%d"}'
				hx-target="#schedule-container">
			`, i)
		}

		if schedule.Tuesday[i] {
			table += fmt.Sprintf(`
			<td
				class="available"
				hx-post = "/select-start"
				hx-vals='{"day":"Tuesday","slot":"%d"}'
				hx-target="#schedule-container"
			></td>
			`, i)
		} else {
			table += fmt.Sprintf(`
			<td
				hx-post="/select-start"
				hx-vals='{"day":"Tuesday","slot":"%d"}'
				hx-target="#schedule-container">
			`, i)
		}

		if schedule.Wednesday[i] {
			table += fmt.Sprintf(`
			<td
				class="available"
				hx-post = "/select-start"
				hx-vals='{"day":"Wednesday","slot":"%d"}'
				hx-target="#schedule-container"
			></td>
			`, i)
		} else {
			table += fmt.Sprintf(`
			<td
				hx-post="/select-start"
				hx-vals='{"day":"Wednesday","slot":"%d"}'
				hx-target="#schedule-container">
			`, i)
		}

		if schedule.Thursday[i] {
			table += fmt.Sprintf(`
			<td
				class="available"
				hx-post = "/select-start"
				hx-vals='{"day":"Thursday","slot":"%d"}'
				hx-target="#schedule-container"
			></td>
			`, i)
		} else {
			table += fmt.Sprintf(`
			<td
				hx-post="/select-start"
				hx-vals='{"day":"Thursday","slot":"%d"}'
				hx-target="#schedule-container">
			`, i)
		}

		if schedule.Friday[i] {
			table += fmt.Sprintf(`
			<td
				class="available"
				hx-post = "/select-start"
				hx-vals='{"day":"Friday","slot":"%d"}'
				hx-target="#schedule-container"
			></td>
			`, i)
		} else {
			table += fmt.Sprintf(`
			<td
				hx-post="/select-start"
				hx-vals='{"day":"Friday","slot":"%d"}'
				hx-target="#schedule-container">
			`, i)
		}

		table += "</tr>"
	}

	data := struct {
		Username string
		Table    template.HTML
		MondayTime	float64
		TuesdayTime float64
		WednesdayTime float64
		ThursdayTime float64
		FridayTime float64
		WeeklyTime float64
	}{
		Username: username,
		Table:    template.HTML(table),
		MondayTime: convertTimeSlotToHours(dailyTotal(&schedule, "Monday")),
		TuesdayTime: convertTimeSlotToHours(dailyTotal(&schedule, "Tuesday")),
		WednesdayTime: convertTimeSlotToHours(dailyTotal(&schedule, "Wednesday")),
		ThursdayTime: convertTimeSlotToHours(dailyTotal(&schedule, "Thursday")),
		FridayTime: convertTimeSlotToHours(dailyTotal(&schedule, "Friday")),
		WeeklyTime: convertTimeSlotToHours(weeklyTotal(&schedule)),
	}

	tmpl, err := template.ParseFiles("templates/schedule.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}


