package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/internal/database"
	"go_final_project/internal/services"
	"io"
	"log"
	"net/http"
	"time"
)

// Task является структкурой необходимой для сериализации и десериализации задачи.
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Result является структкурой необходимой для сериализации http ответа сревера.
type Result struct {
	Tasks []Task `json:"tasks"`
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

var db *sql.DB

// Authorization получает пароль через https, и если он является валидным, 
// то отправляет http ответ, содержащий сгенерированный токен.
func Authorization(w http.ResponseWriter, r *http.Request) {
	var (
		signedToken string
		resp        []byte
		// В переменной takenMap будет содержаться структура с паролем, полученная через https.
		takenMap map[string]string
	)

	body := r.Body
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		log.Println(err.Error())
		w.Write(resp)
		return
	}

	err = json.Unmarshal(data, &takenMap)
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		log.Println(err.Error())
		w.Write(resp)
		return
	}

	signedToken, err = services.GetJWT(takenMap)
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		log.Println(err.Error())
		w.Write(resp)
		return
	}
	resp, _ = json.Marshal(Result{Token: signedToken})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// GetNextDate получает занчения параметров now, date, repeat из парметров запроса и 
// с их помощью возвращает http ответ, содержащий следующую ближайшую дату. 
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	res, err := services.NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// GetTasks возвращает http ответ, содержащий список всех сущетсвующих задач.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	if !services.CheckJWT(w, r) {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	var (
		err   error
		tasks = []Task{}
		resp  []byte
	)

	db, err = sql.Open("sqlite", "../cmd/scheduler.db")
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}
	defer db.Close()

	tasksTmp, err := services.GetTasks(db, r)
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	for _, v := range tasksTmp {
		tasks = append(tasks, Task{Id: v.Id, Date: v.Date, Title: v.Title, Comment: v.Comment, Repeat: v.Repeat})
	}
	resp, _ = json.Marshal(Result{Tasks: tasks})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// DoneTask завершает или обноввляет дату задачи, если поле repeat не пустое и
// вохвращает пустой json http ответа в случае успешного завершения. 
func DoneTask(w http.ResponseWriter, r *http.Request) {
	if !services.CheckJWT(w, r) {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	var (
		task Task
		resp []byte
		id   string
		row  *sql.Row
		err  error
	)

	db, err = sql.Open("sqlite", "../cmd/scheduler.db")
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}
	defer db.Close()

	if r.FormValue("id") == "" {
		resp, _ = json.Marshal(Result{Error: "id не указан или указан некорректно"})
		w.Write(resp)
		log.Println("Id не указан или указан некорректно")
		return
	}

	id = r.FormValue("id")
	row = database.SearchTask(db, id)
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	if r.Method == http.MethodPost {
		if task.Repeat != "" {
			nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				resp, _ = json.Marshal(Result{Error: err.Error()})
				w.Write(resp)
				log.Println(err.Error())
				return
			}

			updateTask := &database.Task{
				Id:      task.Id,
				Date:    nextDate,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			}

			err = updateTask.UpdateTask(db)
			if err != nil {
				resp, _ = json.Marshal(Result{Error: err.Error()})
				w.Write(resp)
				log.Println(err.Error())
				return
			}
		} else {
			err = database.DeleteTask(db, id)
			if err != nil {
				resp, _ = json.Marshal(Result{Error: err.Error()})
				w.Write(resp)
				log.Println(err.Error())
				return
			}
		}
	}

	resp, _ = json.Marshal(Task{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// UpdateTasks возвращает id новой задачи в случае POST запроса,
// вовзращает задачу в случае GET запроса
// и возвращает пустой json http ответа в случае успешного завершения изменения задачи
// с помощью метода PUT и DELETE.
func UpdateTasks(w http.ResponseWriter, r *http.Request) {
	if !services.CheckJWT(w, r) {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	var (
		id   string
		resp []byte
		err  error
		task services.Task
	)

	db, err = sql.Open("sqlite", "../cmd/scheduler.db")
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodPost:
		id, err = services.PostTask(db, w, r)
		resp, _ = json.Marshal(Result{Id: id})
	case http.MethodGet:
		task, err = services.GetTask(db, w, r)
		resp, _ = json.Marshal(Task{Id: task.Id, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat})
	case http.MethodPut:
		resp, err = services.EditTask(db, w, r)
	case http.MethodDelete:
		resp, err = services.DeleteTask(db, w, r)
	}
	if err != nil {
		resp, _ = json.Marshal(Result{Error: err.Error()})
		w.Write(resp)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}