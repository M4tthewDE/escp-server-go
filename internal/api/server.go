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

func (h Handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getUser(w, r)
		return
	}

	if r.Method == "POST" {
		h.saveUser(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

func (h Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var user db.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	savedUser, err := h.dbHandler.GetUser(user.Name)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if savedUser.Pass == user.Pass {
		http.Error(w, "Correct pass", http.StatusOK)
		return
	}

	http.Error(w, "", http.StatusForbidden)
}

func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("user")

	_, err := h.dbHandler.GetUser(userName)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	http.Error(w, "User found", http.StatusFound)
}

func (h Handler) saveUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var user db.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	err = h.dbHandler.SaveUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not save user", http.StatusInternalServerError)

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
