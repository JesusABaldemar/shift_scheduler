package main
import(
	"log"
	"net/http"
	"database/sql"
	_ "modernc.org/sqlite"
)

func main(){
	Users, err := sql.Open("sqlite", "./data/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer Users.Close()
	{
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		schedule TEXT,
		admin BOOLEAN
	)`
	
	if _, err := Users.Exec(query); err != nil {
		log.Fatal(err)

		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home)
	mux.HandleFunc("POST /signin", signIn)
	mux.HandleFunc("POST /login", logIn)
	mux.HandleFunc("GET /schedule", schedule)

	log.Print("Starting server on :4000")

	erro := http.ListenAndServe(":4000", mux)
	log.Fatal(erro)
}

