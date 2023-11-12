package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:NEERA@ujjwa12@tcp(localhost:3306)/yourdatabase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	router.HandleFunc("/items", GetItems).Methods("GET")
	router.HandleFunc("/items/{id}", GetItem).Methods("GET")
	router.HandleFunc("/items", CreateItem).Methods("POST")
	router.HandleFunc("/items/{id}", UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", DeleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	row := db.QueryRow("SELECT id, name FROM items WHERE id = ?", params["id"])
	var item Item
	err := row.Scan(&item.ID, &item.Name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)

	_, err := db.Exec("INSERT INTO items (name) VALUES (?)", newItem.Name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedItem Item
	_ = json.NewDecoder(r.Body).Decode(&updatedItem)

	_, err := db.Exec("UPDATE items SET name = ? WHERE id = ?", updatedItem.Name, params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM items WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
