package user

import (
	"encoding/json"
	"fmt"
	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type newUserBody struct {
	PassportNumber string `json:"passportNumber"`
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
		"PassportNumber": u.PassportNumber,
	}

	for fieldName, fieldValue := range rfs {
		if fieldValue == "" {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	return emptyFields
}

func validatePassport(s string) (string, error) {
	re := regexp.MustCompile(`^(\d{4})\s?(\d{6})$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return "", fmt.Errorf("invalid passport number")
	}

	pn := strings.Join(matches[1:], " ")

	return pn, nil
}
