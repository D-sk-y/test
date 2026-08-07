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

	"test/pkg/mod"
)

type Handler struct {
	store *mod.Store
}

func NewHandler(store *mod.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// URL 全是名词，HTTP 方法区分语义
	mux.HandleFunc("GET    /api/v1/tasks", h.listTasks)
	mux.HandleFunc("POST   /api/v1/tasks", h.createTask)
	mux.HandleFunc("GET    /api/v1/tasks/{id}", h.getTask)
	mux.HandleFunc("PUT    /api/v1/tasks/{id}", h.updateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.deleteTask)
}

// ========== Handler ==========

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	task := h.store.Create(body.Name)
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	task, err := h.store.Update(id, body.Name, body.Status)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id, "deleted": "true"})
}

// ========== 响应工具 ==========

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"message": msg,
	})
}
