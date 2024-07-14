package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-effective-mobile/internal/logger"
	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
	"io"
	"net/http"
	"strconv"
	"time"
)

func New(w http.ResponseWriter, r *http.Request) {
	userIDraw := chi.URLParam(r, "userId")

	userID, err := strconv.Atoi(userIDraw)
	if err != nil {
		logger.Info("task.New Parse userId", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.GetUser(userID)
	if err != nil {
		logger.Info("task.New User not found", "error", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Info("task.New Read body", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var newTask *models.Task
	if len(body) == 0 {
		newTask = &models.Task{}
	} else {
		err = json.Unmarshal(body, &newTask)
		if err != nil {
			logger.Info("task.New Parse body", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	newTask.UserID = userID

	taskID, err := db.NewTask(newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.StartTask(taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write([]byte(strconv.Itoa(taskID))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userIDraw := chi.URLParam(r, "userId")
	startTimeRaw := r.URL.Query().Get("start_time")
	endTimeRaw := r.URL.Query().Get("end_time")

	userID, err := strconv.Atoi(userIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var startTime, endTime *string

	if startTimeRaw != "" {
		parsedStartTime, err := time.Parse(time.DateTime, startTimeRaw)
		if err != nil {
			http.Error(w, "Invalid start time", http.StatusBadRequest)
			return
		}
		s := parsedStartTime.Format(time.RFC3339)
		startTime = &s
	}

	if endTimeRaw != "" {
		parsedEndTime, err := time.Parse(time.DateTime, endTimeRaw)
		if err != nil {
			http.Error(w, "Invalid end time", http.StatusBadRequest)
			return
		}
		s := parsedEndTime.Format(time.RFC3339)
		endTime = &s
		//RFC3339TimeString, err := time.Parse(time.RFC3339, RFC3339Time)
		//endTime = &RFC3339TimeString
	}

	//if startTime != nil && endTime != nil && startTime.After(*endTime) {
	//	http.Error(w, "Start time cannot be after end time", http.StatusBadRequest)
	//	return
	//}

	tasks, err := db.GetUserTasks(userID, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Start(w http.ResponseWriter, r *http.Request) {
	taskIDraw := chi.URLParam(r, "taskId")
	userIDraw := chi.URLParam(r, "userId")

	taskID, err := strconv.Atoi(taskIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.StartTask(taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := map[string]string{
		"message": fmt.Sprintf("Task %s for user %s has been started", taskIDraw, userIDraw),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	taskIDraw := chi.URLParam(r, "taskId")
	userIDraw := chi.URLParam(r, "userId")

	taskID, err := strconv.Atoi(taskIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	duration, err := db.StopTask(taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if duration == "" {
		http.Error(w, fmt.Sprint("task not found"), http.StatusNoContent)
		return
	}

	// TODO: add handling wrong user-task id's relative
	msg := map[string]string{
		"message": fmt.Sprintf("Task %s for user %s has been stopped. Total duration: %s", taskIDraw, userIDraw, duration),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
