package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	app := NewApp()

	if err := http.ListenAndServe(":80", app); err != nil {
		panic(err)
	}
}

/**
App
*/

type App struct {
	*http.ServeMux
}

func NewApp() *App {
	db := NewDb()
	controller := NewController(db)
	mux := http.NewServeMux()

	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.create(w, r)
		case http.MethodGet:
			controller.getAll(w, r)
		}
	})

	return &App{mux}
}

/**
Controller
*/

type Controller struct {
	db *DB
}

func NewController(db *DB) *Controller {
	return &Controller{db}
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		handleUserIDError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value := r.FormValue("value")

	c.db.create(userID, value)

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) getAll(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		handleUserIDError(w, err)
		return
	}

	values := c.db.read(userID)

	w.Write([]byte(strings.Join(values, "\n")))
	w.WriteHeader(http.StatusOK)
}

var ErrEmptyUserID = errors.New("empty user id")
var ErrBrokenUserID = errors.New("broken user id")

func getUserID(r *http.Request) (int, error) {
	userId := r.Header.Get("X-SESSION-USER-ID")
	if len(userId) == 0 {
		return 0, ErrEmptyUserID
	}

	v, err := strconv.Atoi(userId)
	if err != nil {
		return 0, ErrBrokenUserID
	}

	return v, nil
}

func handleUserIDError(w http.ResponseWriter, err error) {
	switch err {
	case ErrEmptyUserID:
		w.WriteHeader(http.StatusUnauthorized)
	case ErrBrokenUserID:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/**
Database
*/

type DB struct {
	*sync.Mutex
	data map[int][]string
}

func NewDb() *DB {
	return &DB{
		new(sync.Mutex),
		make(map[int][]string),
	}
}

func (db *DB) create(userID int, value string) {
	db.Lock()
	defer db.Unlock()

	prev, exists := db.data[userID]
	if !exists {
		db.data[userID] = []string{value}
		return
	}

	db.data[userID] = append(prev, value)
}

func (db *DB) read(userID int) []string {
	db.Lock()
	defer db.Unlock()

	list, exists := db.data[userID]
	if !exists {
		return []string{}
	}

	return list
}
