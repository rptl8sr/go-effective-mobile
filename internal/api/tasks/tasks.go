package tasks

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
	"io"
	"net/http"
	"strconv"
)

func New(w http.ResponseWriter, r *http.Request) {
	userIDraw := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var newTask *models.Task
	err = json.Unmarshal(body, &newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask.UserID = userID

	taskID, err := db.NewTask(newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.StartTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(taskID)))
}

func Get(w http.ResponseWriter, r *http.Request) {
	taskIDraw := chi.URLParam(r, "taskId")

	taskID, err := strconv.Atoi(taskIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(taskID, nil)
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
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
