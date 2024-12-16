package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"mime"
	"net/http"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", insertTask)
	r.Get("/tasks/{id}", getTaskID) //to my opinion need - /task/{id}
	r.Delete("/tasks/{id}", deleteTaskID)

	if err := http.ListenAndServe(":8000", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("/tasks: method - GET at %s\n", r.URL.Path)

	if len(tasks) < 1 {
		http.Error(w, fmt.Sprint("empty tasks"), http.StatusNoContent)
		return
	}
	renderJson(w, tasks)
}

func insertTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("/tasks: method - POST at %s\n", r.URL.Path)

	format := r.Header.Get("Content-type")
	mediaType, _, err := mime.ParseMediaType(format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, fmt.Sprintf("incorrect media-type: %s", mediaType), http.StatusUnsupportedMediaType)
		return
	}
	dataDec := json.NewDecoder(r.Body)
	dataDec.DisallowUnknownFields()

	var task Task
	err = dataDec.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ex := tasks[task.ID]; ex {
		http.Error(w, fmt.Sprintf("task with id=%s, exist", task.ID), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.WriteHeader(http.StatusOK)
}

func getTaskID(w http.ResponseWriter, r *http.Request) {
	log.Printf("/tasks/{id}: metod - GET at %s\n", r.URL.Path)

	id := chi.URLParam(r, "id")
	task, ex := tasks[id]
	if !ex {
		http.Error(w, fmt.Sprintf("task with id=%s, not exist", id), http.StatusBadRequest)
		return
	}
	renderJson(w, task)
}

func deleteTaskID(w http.ResponseWriter, r *http.Request) {
	log.Printf("/tasks/{id}: metod - DELETE at %s\n", r.URL.Path)

	id := chi.URLParam(r, "id")
	_, ex := tasks[id]
	if !ex {
		http.Error(w, fmt.Sprintf("task with id=%s, not exist", id), http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func renderJson(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//http.StatusOK add in ~> func (w *response) write(lenData int, dataB []byte, dataS string) (n int, err error)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
