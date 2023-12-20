package Task

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DNS = "root:root123@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"

type Task struct {
	gorm.Model /*
		//Id       int       `json:"id"`
		username string    `json:"username"`
		taskname string    `json:"taskname"`
		Status   int       `json:"status"`
		Priority int       `json:"priority"`
		Deadline time.Time `json:"deadline"`
		isvalid  bool      `json: "isvalid"`*/
	//gorm.Model

	Username string     `json:"username" gorm:"column:username"`
	Taskname string     `json:"taskname" gorm:"column:taskname"`
	Status   int        `json:"status"`
	Priority int        `json:"priority"`
	Deadline *time.Time `json:"deadline"` // Use a pointer to time.Time to allow NULL
	Isvalid  bool       `json:"isvalid" gorm:"column:isvalid"`
}

func InitialMigration() {
	/*DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&Task{})*/
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		log.Println("Error connecting to the database:", err)
		panic("Cannot connect to DB")
	}

	log.Println("Connected to the database")

	err = DB.AutoMigrate(&Task{})
	if err != nil {
		log.Println("Error auto-migrating tables:", err)
	} else {
		log.Println("Tables auto-migrated successfully")
	}
}

func View_tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task_tbl []Task
	if err := DB.Find(&task_tbl).Error; err != nil {
		log.Println("Error fetching tasks:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}
	fmt.Print("data showing")
	json.NewEncoder(w).Encode(task_tbl)
}

/*
func Create_task(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task_tbl Task
	json.NewDecoder(r.Body).Decode(&task_tbl)
	DB.Create(&task_tbl)
	json.NewEncoder(w).Encode(task_tbl)
	fmt.Print("data inserted...")

}*/
func Create_task(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task_tbl Task
	if err := json.NewDecoder(r.Body).Decode(&task_tbl); err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Log the received task data
	log.Printf("Received Task: %+v\n", task_tbl)

	// Attempt to create the task
	if err := DB.Create(&task_tbl).Error; err != nil {
		log.Println("Error creating task:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log a success message
	log.Println("Task created successfully")

	// Respond with the created task
	json.NewEncoder(w).Encode(task_tbl)
}

func Get_task(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var task_tbl Task
	DB.First(&task_tbl, params["id"])
	json.NewEncoder(w).Encode(task_tbl)
}

func Update_task(w http.ResponseWriter, r *http.Request) {
	// Parse the task ID from the request parameters
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the JSON request body
	var updateData map[string]int
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Extract the new status value from the JSON data
	newStatus, ok := updateData["status"]
	if !ok {
		http.Error(w, "Status field is required", http.StatusBadRequest)
		return
	}

	// Update the task status in the database
	result := DB.Model(&Task{}).Where("id = ?", taskID).Update("status", newStatus)
	if result.Error != nil {
		http.Error(w, "Failed to update task status", http.StatusInternalServerError)
		return
	}

	// Check if the task was found and updated
	if result.RowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task status updated successfully")
}

/* json values like

{
    "id": 3,
    "username": "shiva",
    "taskname": "create v2 project",
    "status": 0,
    "priority": 1,
    "deadline": "2023-12-19T23:59:59Z",
    "isvalid": true
}
*/
