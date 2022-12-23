package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	dapatId := r.Context().Value("id").(string)
	if dapatId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	cekIdInt, _ := strconv.Atoi(dapatId)
	idTask := r.URL.Query().Get("task_id")

	if idTask == "" {
		res, err := t.taskService.GetTasks(r.Context(), cekIdInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
	idTaskInt, _ := strconv.Atoi(idTask)
	if idTask != "" {
		res, err := t.taskService.GetTaskByID(r.Context(), idTaskInt)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	if task.Description == "" || task.Title == "" || task.CategoryID == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	dapatId := r.Context().Value("id").(string)
	if dapatId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idCekInt, _ := strconv.Atoi(dapatId)
	res, err := t.taskService.StoreTask(r.Context(), &entity.Task{
		ID:          task.ID,
		Description: task.Description,
		Title:       task.Title,
		UserID:      idCekInt,
		CategoryID:  task.CategoryID,
	})

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	hasil := map[string]interface{}{
		"user_id": idCekInt,
		"message": "success create new task",
		"task_id": res.ID,
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(hasil)
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idTask := r.URL.Query().Get("task_id")
	dapatId := r.Context().Value("id").(string)
	if dapatId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).
			Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	cekTaskIdInteger, _ := strconv.Atoi(idTask)
	err := t.taskService.DeleteTask(r.Context(), cekTaskIdInteger)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).
			Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	id, _ := strconv.Atoi(dapatId)
	res := map[string]interface{}{
		"user_id": id,
		"message": "success delete task",
		"task_id": cekTaskIdInteger,
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	dapatId := r.Context().Value("id").(string)
	if dapatId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	cekIdTask, _ := strconv.Atoi(dapatId)
	res, err := t.taskService.UpdateTask(r.Context(), &entity.Task{
		ID:          task.ID,
		Description: task.Description,
		Title:       task.Title,
		UserID:      cekIdTask,
		CategoryID:  task.CategoryID,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).
			Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).
		Encode(map[string]interface{}{
			"user_id": cekIdTask,
			"message": "success update task",
			"task_id": res.ID,
		})
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	idUser := r.Context().Value("id")
	idLoginUser, err := strconv.Atoi(idUser.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		UserID:     int(idLoginUser),
		CategoryID: task.CategoryID,
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idUser,
		"message": "success update task category",
		"task_id": task.ID,
	})
}
