package api

import (
	"context"
	"fmt"
	"test/pkg/mod"
)

type Gateway struct {
	bus *mod.Bus
	ctx context.Context
}

func NewGateway(bus *mod.Bus, ctx context.Context) *Gateway {
	return &Gateway{
		bus: bus,
		ctx: ctx,
	}
}

func (g *Gateway) GetContext() context.Context {
	return g.ctx
}

func (g *Gateway) Initialize(c context.Context) error {
	fmt.Println("Gateway initialized with context:", g.GetContext())
	return nil
}

func (g *Gateway) DeInitialize() error {
	fmt.Println("Gateway deinitialized")
	return nil
}
