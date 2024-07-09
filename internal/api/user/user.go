package user

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
)

type userBody struct {
	PassportNumber *string `json:"passportNumber,omitempty"`
	Surname        *string `json:"surname,omitempty"`
	Name           *string `json:"name,omitempty"`
	Patronymic     *string `json:"patronymic,omitempty"`
	Address        *string `json:"address,omitempty"`
}

func New(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var newUser userBody
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

	pn, err := validatePassport(*newUser.PassportNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		PassportNumber: pn,
		Surname:        safeDereference(newUser.Surname),
		Name:           safeDereference(newUser.Name),
		Patronymic:     safeDereference(newUser.Patronymic),
		Address:        safeDereference(newUser.Address),
	}

	id, err := db.NewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdUser, err := db.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userToUpdate userBody

	err = json.NewDecoder(r.Body).Decode(&userToUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pn string
	if userToUpdate.PassportNumber != nil {
		pn, err = validatePassport(*userToUpdate.PassportNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	user := models.User{
		ID:             id,
		PassportNumber: pn,
		Surname:        safeDereference(userToUpdate.Surname),
		Name:           safeDereference(userToUpdate.Name),
		Patronymic:     safeDereference(userToUpdate.Patronymic),
		Address:        safeDereference(userToUpdate.Address),
	}

	user.ID = id
	user.PassportNumber = pn

	err = db.UpdateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := db.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
