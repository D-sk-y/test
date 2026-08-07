// Package v2 — RPC-over-HTTP 风格
//
// 核心特征：
//   URL 含动词，表达"要做什么"（/startTask、/stopTask、/listTasks）
//   几乎所有操作用 POST
//   每个动作映射一条独立路由
//   响应格式与 v1 完全相同（共享 Store），对照时只看 HTTP 层差异

package v2

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
	// URL 带动词，全部 POST —— 这是 RPC-over-HTTP 的标志
	mux.HandleFunc("POST /api/v2/listTasks", h.listTasks)
	mux.HandleFunc("POST /api/v2/createTask", h.createTask)
	mux.HandleFunc("POST /api/v2/getTask", h.getTask)
	mux.HandleFunc("POST /api/v2/updateTask", h.updateTask)
	mux.HandleFunc("POST /api/v2/deleteTask", h.deleteTask)
}

// ========== Handler ==========
// 注意：URL 语义和 v1 完全相反。v1 是 GET /tasks，v2 是 POST /listTasks

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
	writeJSON(w, http.StatusOK, task) // RPC 风格通常统一返回 200
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	task, err := h.store.Get(body.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	task, err := h.store.Update(body.ID, body.Name, body.Status)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := h.store.Delete(body.ID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": body.ID, "deleted": "true"})
}

// ========== 响应工具（与 v1 完全相同） ==========

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
