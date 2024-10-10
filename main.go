package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Employee structure
type Employee struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

var db *sql.DB

// Connect to the PostgreSQL database
func connectDB() {
	var err error
	connStr := "user=postgres password=admin123 dbname=employee_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}
	fmt.Println("Connected to the database successfully!")
}

// Create an employee
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO employees (name, position) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, emp.Name, emp.Position).Scan(&emp.ID)
	if err != nil {
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// Get all employees
func getEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, position FROM employees")
	if err != nil {
		http.Error(w, "Failed to retrieve employees", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Position)
		if err != nil {
			http.Error(w, "Failed to scan employee", http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// Update an employee by ID
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE employees SET name=$1, position=$2 WHERE id=$3`
	result, err := db.Exec(query, emp.Name, emp.Position, emp.ID)
	if err != nil {
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

// Delete an employee by ID
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM employees WHERE id=$1`
	result, err := db.Exec(query, emp.ID)
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employee with ID %d deleted", emp.ID)
}

func main() {
	// Connect to the database
	connectDB()
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/employees", getEmployees)     // GET to retrieve all employees
	http.HandleFunc("/employees/create", createEmployee) // POST to create a new employee
	http.HandleFunc("/employees/update", updateEmployee) // PUT to update an existing employee
	http.HandleFunc("/employees/delete", deleteEmployee) // DELETE to delete an employee

	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
