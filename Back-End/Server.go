package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Reyan-1103"
	dbname   = "postgres"
)

// Student struct to hold student data
type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Studies struct {
	StudentId int `json:"studentid"`
	CourseId int `json:"courseid"`
}

type Studentz struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Cource struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Seats int `json:"seats"`
}
var CourceName=""
var CourceSeats=0
// In-memory database (slice) to store students
var cources []Cource
var students []Student
var studs []Studentz
var studies []Studies

var (
	studentsMutex sync.Mutex
	courcesMutex  sync.Mutex
	studiesMutex  sync.Mutex
)

var dsn = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname,
)
var db, err = sql.Open("postgres", dsn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

var clients = make(map[*websocket.Conn]bool) 
var broadcast = make(chan Cource)           

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	for {
		var updatedCource Cource
		err := conn.ReadJSON(&updatedCource)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(clients, conn)
			break
		}
	}
}

func broadcastUpdates() {
	for {
		updatedCource := <-broadcast
		for client := range clients {
			err := client.WriteJSON(updatedCource)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// func createTable(db *sql.DB) {
// 	query := `
// 	CREATE TABLE IF NOT EXISTS Studies (
// 		StudentId INT,
// 		CourseId INT,
// 		FOREIGN KEY (StudentId) REFERENCES Students(id),
// 		FOREIGN KEY (CourseId) REFERENCES Cources(id)
// 	);`
// // // 	query = `
// // // 	CREATE TABLE IF NOT EXISTS Students (
// // // 		id SERIAL PRIMARY KEY,
// // // 		name TEXT NOT NULL,
// // // 		password TEXT
// // // 	);`
// // // 	query = `
// // // 	CREATE TABLE IF NOT EXISTS Cources (
// // // 		id SERIAL PRIMARY KEY,
// // // 		name TEXT NOT NULL,
// // // 		seats int,
// // // 	);`
// 	_, err := db.Exec(query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Table created successfully")
// }

func removeCourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var coursesToRemove []Cource
	if err := json.NewDecoder(r.Body).Decode(&coursesToRemove); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, course := range coursesToRemove {
		query := `DELETE FROM Cources WHERE id = $1`
		_, err := db.Exec(query, course.ID)
		if err != nil {
			log.Printf("Failed to remove course ID %d: %v", course.ID, err)
		}
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Courses removed successfully!"})
}

// 🚀 POST /addStudent → Adds a new student
func addStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body
	var newStudent Student
	if err := json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	studentsMutex.Lock()
	students = append(students, newStudent)
	studentsMutex.Unlock()

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Student added successfully!"})
}
func updateCource(db *sql.DB, cour []Cource) {
	for i := 0; i < len(cour); i++ {
		query := `UPDATE Cources SET seats = $1 WHERE id = $2;`
		_, err := db.Exec(query, cour[i].Seats, cour[i].ID)
		if err != nil {
			log.Printf("Update failed for Cource ID %d: %v", cour[i].ID, err)
		} else {
			// Broadcast the updated course to all WebSocket clients
			broadcast <- cour[i]
		}
	}
	fmt.Println("Cources updated successfully")
}

func updateCourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body
	var courceList []Cource
	if err := json.NewDecoder(r.Body).Decode(&courceList); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the courses in the database
	updateCource(db, courceList)

	// Respond with a success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Courses updated successfully!"})
}

// 🚀 GET /students → Retrieves all students
func getStudentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch students directly from the database
	studs := getStudents(db)

	// Respond with the list of students
	json.NewEncoder(w).Encode(studs)
}

// 🚀 GET /students → Retrieves all cources
func getCourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch courses directly from the database
	cources := getCources(db)

	// Respond with the list of courses
	json.NewEncoder(w).Encode(cources)
}

// 🚀 GET /Studies → Retrieves all cources
func getStudiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch studies directly from the database
	studies := getStudies(db)

	// Respond with the list of studies
	json.NewEncoder(w).Encode(studies)
}

func getStudents (db *sql.DB) []Studentz{
	rows, err := db.Query("SELECT id, name,password FROM Students")
	if (err != nil) {
		log.Fatal(err)
	}
	defer rows.Close()
	var students []Studentz

	for rows.Next() {
		var student Studentz
		rows.Scan(&student.ID, &student.Name, &student.Password)
		students = append(students, student)
		studs=append(studs, student)
	}

	return students
}
func getStudies(db *sql.DB) []Studies {
	rows, err := db.Query("SELECT StudentId, CourseId FROM Studies") // Fixed column name
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var study []Studies

	for rows.Next() {
		var stu Studies
		rows.Scan(&stu.StudentId, &stu.CourseId)
		study = append(study, stu)
		}
	studies = study
	return study
}
func getCources (db *sql.DB) []Cource{
	rows, err := db.Query("SELECT id, name,seats FROM Cources")
	if (err != nil) {
		log.Fatal(err)
	}
	defer rows.Close()
	var cource []Cource

	for rows.Next() {
		var cour Cource
		rows.Scan(&cour.ID, &cour.Name, &cour.Seats)
		cource=append(cource, cour)
	}
	cources=cource	
	return cource
}
func insertCource(db *sql.DB, name string, seats int) {
	
	query := `INSERT INTO Cources (name,seats) VALUES ($1,$2) RETURNING id`
	var id int
	err := db.QueryRow(query, name,seats).Scan(&id)
	if (err != nil) {
		log.Fatal("Insert failed:", err)
	}
}

func insertStudies(db *sql.DB, studentId int, courseId int) {
	query := `INSERT INTO Studies (StudentId, CourseId) VALUES ($1, $2)`
	_, err := db.Exec(query, studentId, courseId)
	if err != nil {
		log.Printf("Insert into Studies failed: %v", err)
		return
	}

	// Decrease seat count by 1 for the registered course
	updateSeats(db, courseId, -1)
}

func addCourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body
	var newCource Cource
	if err := json.NewDecoder(r.Body).Decode(&newCource); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the course into the database
	insertCource(db, newCource.Name, newCource.Seats)

	// Broadcast the new course to all WebSocket clients
	broadcast <- newCource

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Course added successfully!"})
}

func removeStudiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var studiesToRemove []Studies
	if err := json.NewDecoder(r.Body).Decode(&studiesToRemove); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, study := range studiesToRemove {
		query := `DELETE FROM Studies WHERE CourseId = $1 AND StudentId = $2`
		_, err := db.Exec(query, study.CourseId, study.StudentId)
		if err != nil {
			log.Printf("Failed to remove study record (CourseId: %d, StudentId: %d): %v", study.CourseId, study.StudentId, err)
			continue
		}

		// Increase seat count by 1 for the dropped course
		updateSeats(db, study.CourseId, 1)
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Studies removed successfully!"})
}

func addStudiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body
	var studiesToAdd []Studies
	if err := json.NewDecoder(r.Body).Decode(&studiesToAdd); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, study := range studiesToAdd {
		insertStudies(db, study.StudentId, study.CourseId)
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Studies added successfully!"})
}

func updateSeats(db *sql.DB, courseId int, seatChange int) {
	query := `UPDATE Cources SET seats = seats + $1 WHERE id = $2`
	result, err := db.Exec(query, seatChange, courseId)
	if err != nil {
		log.Printf("Failed to update seats for CourseId %d: %v", courseId, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to fetch rows affected for CourseId %d: %v", courseId, err)
		return
	}

	if rowsAffected == 0 {
		log.Printf("No rows updated for CourseId %d. Ensure the course exists.", courseId)
	} else {
		log.Printf("Successfully updated seats for CourseId %d by %d.", courseId, seatChange)

		// Fetch the updated course details
		query = `SELECT id, name, seats FROM Cources WHERE id = $1`
		var updatedCourse Cource
		err := db.QueryRow(query, courseId).Scan(&updatedCourse.ID, &updatedCourse.Name, &updatedCourse.Seats)
		if err != nil {
			log.Printf("Failed to fetch updated course details for CourseId %d: %v", courseId, err)
			return
		}

		// Broadcast the updated course to all WebSocket clients
		broadcast <- updatedCourse
	}
}

// 🚀 POST /login → Handles student login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginData struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	for _, student := range studs {
		if student.ID == loginData.ID && student.Password == loginData.Password {
			json.NewEncoder(w).Encode(map[string]string{"message": "Login successful!"})
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
func getInfo() {
    var wg sync.WaitGroup

    wg.Add(3) 

    go func() {
        defer wg.Done()
        getStudents(db)
    }()

    go func() {
        defer wg.Done()
        getCources(db)
    }()

    go func() {
        defer wg.Done()
        getStudies(db)
    }()

    wg.Wait()
}

func main() {
	if err != nil {
		log.Fatal(err)
	}
	
	// Initialize some sample students
	students = []Student{
		{Name: "Alice", Age: 22},
		{Name: "Bob", Age: 24},
	}

	// Create HTTP router
	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/removeCource", removeCourceHandler) // POST

	mux.HandleFunc("/addStudent", addStudentHandler) // POST

	mux.HandleFunc("/updateCource", updateCourceHandler) // POST

	mux.HandleFunc("/students", getStudentsHandler)  // GET

	mux.HandleFunc("/cources",getCourcesHandler)  // GET

	mux.HandleFunc("/studies",getStudiesHandler)  // GET

	mux.HandleFunc("/addCource", addCourceHandler)  // GET
	
	mux.HandleFunc("/removeStudies", removeStudiesHandler) // POST

	mux.HandleFunc("/addStudies", addStudiesHandler) // POST

	mux.HandleFunc("/login", loginHandler) // POST

	// WebSocket endpoint
	mux.HandleFunc("/ws", websocketHandler)

	// Start the WebSocket broadcaster
	go broadcastUpdates()

	// Configure CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"}, // Allow all localhost ports
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	
	
	
		getInfo()

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", corsHandler)
}
