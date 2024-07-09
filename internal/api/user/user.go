package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
)

type newUserBody struct {
	PassportNumber string `json:"passportNumber"`
}

type newUserResponse struct {
	ID int `json:"id"`
}

func New(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var newUser newUserBody
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	emptyFields := requiredFieldsNewUser(&newUser)

	if len(emptyFields) > 0 {
		errorMessage := "Required fields are empty: " + strings.Join(emptyFields, ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	pn, err := validatePassport(newUser.PassportNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		PassportNumber: pn,
	}

	id, err := db.NewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := newUserResponse{ID: id}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func Get(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.GetUser(id)
	if errors.Is(err, db.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Update(w http.ResponseWriter, r *http.Request) {

}
