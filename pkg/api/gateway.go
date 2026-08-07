package api

import (
	"context"
	"fmt"
	"net/http"
	"test/pkg/mod"
	"time"

	v1 "test/pkg/api/v1"
	v2 "test/pkg/api/v2"
)

type Gateway struct {
	bus    *mod.Bus
	ctx    context.Context
	server *http.Server
	store  *mod.Store
}

func NewGateway(bus *mod.Bus, ctx context.Context) *Gateway {
	return &Gateway{
		bus:   bus,
		ctx:   ctx,
		store: mod.NewStore(),
	}
}

func (g *Gateway) GetContext() context.Context {
	return g.ctx
}

func (g *Gateway) Initialize(c context.Context) error {
	fmt.Println("Gateway initialized with context:", g.GetContext())
	mux := http.NewServeMux()
	v1.NewHandler(g.store).RegisterRoutes(mux)
	v2.NewHandler(g.store).RegisterRoutes(mux)
	g.server = &http.Server{Addr: ":18080", Handler: mux}
	go func() {
		fmt.Println("HTTP server listening on :18080")
		fmt.Println("  RESTful v1:  GET/POST/PUT/DELETE /api/v1/tasks[/{id}]")
		fmt.Println("  RPC     v2:  POST /api/v2/{listTasks,createTask,getTask,updateTask,deleteTask}")
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Gateway error:", err)
		}
	}()
	return nil
}

func (g *Gateway) DeInitialize() error {
	fmt.Println("Gateway deinitialized")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return g.server.Shutdown(ctx)
}
