package user

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
)

type newUserBody struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic,omitempty"`
	Address    string `json:"address"`
	Passport   string `json:"passport,omitempty"`
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

	emptyFields := requiredFields(&newUser)

	if len(emptyFields) > 0 {
		errorMessage := "Required fields are empty: " + strings.Join(emptyFields, ", ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	user := &models.User{
		Surname:    newUser.Surname,
		Name:       newUser.Name,
		Patronymic: newUser.Patronymic,
		Address:    newUser.Address,
		Passport:   newUser.Passport,
		CreatedAt:  time.Time{},
	}

	id, err := db.NewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func requiredFields(u *newUserBody) []string {
	var emptyFields []string

	rfs := map[string]string{
		"Surname": u.Surname,
		"Name":    u.Name,
		"Address": u.Address,
	}

	for fieldName, fieldValue := range rfs {
		if fieldValue == "" {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	return emptyFields
}
