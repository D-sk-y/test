// Package v2 — RPC-over-HTTP 风格
//
// 核心特征：
//   URL 含动词，表达"要做什么"（/startTask、/stopTask、/listTasks）
//   几乎所有操作用 POST
//   每个动作映射一条独立路由
//   响应格式与 v1 完全相同（共享 Store），对照时只看 HTTP 层差异

package v2

import (
	"net/http"
	"test/db"
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
	// URL 带动词，全部 POST —— 这是 RPC-over-HTTP 的标志
	// mux.HandleFunc("POST /api/v2/listTasks", h.listTasks)
	// mux.HandleFunc("POST /api/v2/createTask", h.createTask)
	// mux.HandleFunc("POST /api/v2/getTask", h.getTask)
	// mux.HandleFunc("POST /api/v2/updateTask", h.updateTask)
	// mux.HandleFunc("POST /api/v2/deleteTask", h.deleteTask)
}

// ========== Handler ==========
// 注意：URL 语义和 v1 完全相反。v1 是 GET /tasks，v2 是 POST /listTasks
