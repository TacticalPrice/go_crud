package main

import(
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Item struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var items []Item

func main() {
	router := mux.NewRouter()
    
	items = append(items, Item{ID: "1", Name: "Item 1",})
	items = append(items, Item{ID: "2", Name: "Item 2",})
	items = append(items, Item{ID: "3", Name: "Item 3",})


	router.HandleFunc("/items" , GetItems).Methods("GET")
	router.HandleFunc("/items/{id}" , GetItem).Methods("GET")
	router.HandleFunc("/items", CreateItem).Methods("POST")
	router.HandleFunc("/items{id}" , UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}" , DeleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000" , router))
}

func GetItems(w http.ResponseWriter , r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter , r *http.Request) {
	params := mux.Vars(r)
	for _, item := range items {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func CreateItem(w http.ResponseWriter , r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	items = append(items , item)
	json.NewEncoder(w).Encode(items)
}

func UpdateItem(w http.ResponseWriter , r *http.Request) {
	params := mux.Vars(r)
	for i , item := range items {
		if item.ID == params["id"] {
			items = append(items[:i] , items[i+1:]...)
			var newItem Item
			_ = json.NewDecoder(r.Body).Decode(&newItem)
			newItem.ID = params["id"]
			items = append(items , newItem)
			json.NewEncoder(w).Encode(items)
			return
		}
	}
	json.NewEncoder(w).Encode(items)
}

func DeleteItem(w http.ResponseWriter , r *http.Request) {
	params := mux.Vars(r)
	for i, item := range items {
		if item.ID == params["id"] {
			items = append(items[:i] , items[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(items)
}