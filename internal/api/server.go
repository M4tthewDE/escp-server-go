package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/m4tthewde/escp-server-go/internal/db"
)

type Handler struct {
	dbHandler db.DatabaseHandler
}

func NewHandler() Handler {
	dbHandler := db.NewDatabaseHandler()

	handler := Handler{
		dbHandler: dbHandler,
	}

	return handler
}

func (h Handler) GetCountries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	countries, err := h.dbHandler.GetCountries()
	if err != nil {
		http.Error(w, "Could not fetch countries", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(countries)
	if err != nil {
		http.Error(w, "Could not return countries", http.StatusInternalServerError)

		return
	}
}

func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userName := r.URL.Query().Get("user")

	user, err := h.dbHandler.GetUser(userName)
	if err.Error() == "User not found" {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Could not return user", http.StatusInternalServerError)

		return
	}
}

func (h Handler) SetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var result db.ResultDto

	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	err = h.dbHandler.SaveResult(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not save result", http.StatusInternalServerError)

		return
	}
}

func (h Handler) SetRanking(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var ranking db.RankingDto

	err := json.NewDecoder(r.Body).Decode(&ranking)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	err = h.dbHandler.SaveRanking(ranking)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not save result", http.StatusInternalServerError)

		return
	}
}
