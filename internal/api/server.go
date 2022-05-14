package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func (h Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	user := r.URL.Query().Get("user")

	if user == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ranking, err := h.dbHandler.GetRanking(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get ranking", http.StatusInternalServerError)

		return
	}

	adminRanking, err := h.dbHandler.GetRanking("admin")
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get ranking", http.StatusInternalServerError)

		return
	}

	result := calcResult(ranking, adminRanking)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Could not return result", http.StatusInternalServerError)

		return
	}
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
	var ranking db.Ranking

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

func (h Handler) HandleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetDone(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h Handler) GetDone(w http.ResponseWriter, r *http.Request) {
	lock, err := h.dbHandler.GetDone()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get done", http.StatusInternalServerError)

		return
	}

	fmt.Fprintln(w, lock)
}

type Result struct {
	Points   int
	Accuracy map[int]int
}

func calcResult(ranking *db.Ranking, adminRanking *db.Ranking) *Result {
	accuracy := make(map[int]int)
	points := 0

	for i, country := range ranking.Ranking {
		adminIndex := find(country.Name, adminRanking)

		delta := int(math.Abs(float64(adminIndex - i)))
		accuracy[i] = delta

		if delta == 0 {
			points += 3
		}

		if delta == 1 {
			points += 2
		}

		if delta == 2 {
			points += 1
		}
	}

	return &Result{
		Points:   points,
		Accuracy: accuracy,
	}
}

func find(name string, ranking *db.Ranking) int {
	for i, country := range ranking.Ranking {
		if country.Name == name {
			return i
		}
	}

	return -1
}
