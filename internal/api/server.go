package api

import (
	"encoding/json"
	"fmt"
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

func (h Handler) HandleRanking(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.SetRanking(w, r)
		return
	}

	if r.Method == "GET" {
		h.GetRanking(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h Handler) SetRanking(w http.ResponseWriter, r *http.Request) {
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

func (h Handler) GetRanking(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	if user == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ranking, err := h.dbHandler.GetRanking(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not save result", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(ranking)
	if err != nil {
		http.Error(w, "Could not return countries", http.StatusInternalServerError)

		return
	}
}

func (h Handler) HandleLock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetLock(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h Handler) GetLock(w http.ResponseWriter, r *http.Request) {
	lock, err := h.dbHandler.GetLock()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get lock", http.StatusInternalServerError)

		return
	}

	fmt.Fprintln(w, lock)
}
