// Package v1 — RESTful 风格
//
// 核心特征：
//   URL 只有名词（/tasks、/tasks/{id}），没有动词
//   HTTP 方法区分操作：GET=读  POST=创建  PUT=更新  DELETE=删除
//   同一资源同一 URL，不同方法不同行为

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test/db"
	"test/entity"
)

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// URL 全是名词，HTTP 方法区分语义
	mux.HandleFunc("GET    /api/v1/tasks", h.listTasks)
	mux.HandleFunc("POST   /api/v1/tasks", h.createTask)
	// mux.HandleFunc("GET    /api/v1/tasks/{id}", h.getTask)
	mux.HandleFunc("PUT    /api/v1/tasks/{id}", h.updateTask)
	// mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.deleteTask)
}

// ========== Handler ==========

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.db.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 获取新值
	var newTask entity.Task
	newTask, err := h.db.CreateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// 返回新值
	json.NewEncoder(w).Encode(newTask)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	// 从路径参数中获取 ID
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// body 解析为 task
	var body struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id == 0 && body.ID != 0 {
		id = body.ID
	}

	task := entity.Task{
		ID:     id,
		Name:   body.Name,
		Status: body.Status,
	}
	err = h.db.UpdateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
