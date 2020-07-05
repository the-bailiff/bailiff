package main

import (
	"errors"
	"net/http"
	"strconv"
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
	ctrl := NewAuthController(db)
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", ctrl.register)
	mux.HandleFunc("/auth/login", ctrl.login)
	mux.HandleFunc("/auth/logout", ctrl.logout)

	return &App{mux}
}

/**
Controller
*/

type AuthController struct {
	*DB
}

func NewAuthController(db *DB) *AuthController {
	return &AuthController{db}
}

func (c *AuthController) register(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.DB.create(username, password)
	if err != nil {
		if err == ErrAlreadyExists {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("X-SESSION-SET-USER-ID", strconv.Itoa(user.id))
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := c.DB.read(username, password)
	if err != nil {
		if err == ErrNotExists {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("X-SESSION-SET-USER-ID", strconv.Itoa(user.id))
}

func (c *AuthController) logout(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-SESSION-DEL-USER-ID", "")
}

func isAuthorized(r *http.Request) bool {
	userId := r.Header.Get("X-SESSION-USER-ID")
	return len(userId) != 0
}

/**
Database
*/

type User struct {
	id       int
	username string
	password string
}

var ErrAlreadyExists = errors.New("already exists")
var ErrNotExists = errors.New("not exists")

type DB struct {
	*sync.Mutex
	data map[string]*User
	id   int
}

func NewDb() *DB {
	return &DB{
		new(sync.Mutex),
		make(map[string]*User),
		0,
	}
}

func (db *DB) create(username, password string) (*User, error) {
	db.Lock()
	defer db.Unlock()

	_, exists := db.data[username]
	if exists {
		return nil, ErrAlreadyExists
	}

	db.id++
	user := &User{db.id, username, password}
	db.data[username] = user

	return user, nil
}

func (db *DB) read(username, password string) (*User, error) {
	db.Lock()
	defer db.Unlock()

	user, exists := db.data[username]
	if !exists {
		return nil, ErrNotExists
	}
	if user.password != password {
		return nil, ErrNotExists
	}

	return user, nil
}
