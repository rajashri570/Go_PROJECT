package Task

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
