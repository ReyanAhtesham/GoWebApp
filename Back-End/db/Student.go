// package main

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"

// 	_ "github.com/lib/pq"
// 	"github.com/rs/cors"
// )

// // Define request and response structures
// type RequestData struct {
// 	Message string `json:"message"`
// }

// type ResponseData struct {
// 	Reply string `json:"reply"`
// }

// // Handler function to process user messages

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "Reyan-1103"
// 	dbname   = "postgres"
// )
// // Table Studies{
// // 	StudentId integar
// // 	CourceId integar
// //   }
// //   Table Cources{
// // 	id integar
// // 	name varchar
// // 	seats int
// //   }
// //   Table Students {
// // 	id integer [primary key]
// // 	username varchar
// // 	password varchar
// //   } 
// // func createTable(db *sql.DB) {
// // 	query := `
// // 	CREATE TABLE IF NOT EXISTS Studies (
// // 		StudentId INT,
// // 		CourseId INT,
// // 		FOREIGN KEY (StudentId) REFERENCES Students(id),
// // 		FOREIGN KEY (CourseId) REFERENCES Cources(id)
// // 	);`
// // 	query = `
// // 	CREATE TABLE IF NOT EXISTS Students (
// // 		id SERIAL PRIMARY KEY,
// // 		name TEXT NOT NULL,
// // 		password TEXT
// // 	);`
// // 	query = `
// // 	CREATE TABLE IF NOT EXISTS Cources (
// // 		id SERIAL PRIMARY KEY,
// // 		name TEXT NOT NULL,
// //		seats int,
// // 	);`
// // 	_, err := db.Exec(query)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	fmt.Println("Table created successfully")
// // }

// func insertStudent(db *sql.DB, name string, password string) {
// 	query := `INSERT INTO Students (name,password) VALUES ($1,$2) RETURNING id`
// 	var id int
// 	err := db.QueryRow(query, name,password).Scan(&id)
// 	if err != nil {
// 		log.Fatal("Insert failed:", err)
// 	}
// 	fmt.Printf("Inserted user with ID: %d\n", id)
// }
// func insertCource(db *sql.DB, name string, seats int) {
// 	query := `INSERT INTO Cources (name,seats) VALUES ($1,$2) RETURNING id`
// 	var id int
// 	err := db.QueryRow(query, name,seats).Scan(&id)
// 	if err != nil {
// 		log.Fatal("Insert failed:", err)
// 	}
// 	fmt.Printf("Inserted user with ID: %d\n", id)
// }
// type Student struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// }

// type SelectionResponse struct {
// 	Success bool `json:"success"`
// }

// func getStudents (db *sql.DB) []Student{
// 	rows, err := db.Query("SELECT id, name,password FROM Students")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	var students []Student

// 	for rows.Next() {
// 		var student Student
// 		rows.Scan(&student.ID, &student.Name, &student.Password)
// 		students = append(students, student)
// 		fmt.Printf("ID: %d, Name: %s, Password: %s\n",student.ID, student.Name,student.Password)
// 	}
// 	return students
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Read request body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	// Parse JSON
// 	var request RequestData
// 	if err := json.Unmarshal(body, &request); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	// Create response based on user input
// 	var response ResponseData

// 	switch request.Message {
// 	case "hello":
// 		response.Reply = "Hello! How can I help you?"
// 	case "bye":
// 		response.Reply = "Goodbye! Have a nice day!"
// 	default:
// 		response.Reply = "I don't understand that message."
// 	}	


// 	// Send JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func sendStudents(students []Student) bool {
// 	jsonData, _ := json.Marshal(students)

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", "http://localhost:8080/validate", bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Access-Control-Allow-Origin", "*") // Allow all origins (adjust as needed)
// 	req.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	req.Header.Set("Access-Control-Allow-Headers", "Content-Type")

// 	resp, _ := client.Do(req)
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	var result SelectionResponse
// 	json.Unmarshal(body, &result)

// 	return result.Success
// }


// func main00() {

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/chat", handler)

// 	// Configure CORS
// 	corsHandler := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow React frontend
// 		AllowedMethods:   []string{"POST", "OPTIONS"},
// 		AllowedHeaders:   []string{"Content-Type"},
// 		AllowCredentials: true,
// 	}).Handler(mux)

// 	// Start server
// 	fmt.Println("Server running on http://localhost:8080")
// 	http.ListenAndServe(":8080", corsHandler)

	
// 	dsn := fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname,
// 	)

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// //	createTable(db)
// students := getStudents(db)
// isValid := sendStudents(students)
// if(isValid){
// 	println("Wowww Crazyyy")
// }else{
// 	println("ohooooo")
// }
// //getStudents(db)
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to PostgreSQL successfully!")
// }
