package mod

import (
	"context"
	"sync"
)

// 接口
type Module interface {
	Initialize(c context.Context) error
	DeInitialize() error
}

type ModuleBase struct {
	modules map[string]Module
	mutex   sync.RWMutex
}

func NewModuleBase() *ModuleBase {
	return &ModuleBase{
		modules: make(map[string]Module),
	}
}

func (mb *ModuleBase) SetModule(name string, module Module) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	mb.modules[name] = module
}

func (mb *ModuleBase) GetModule(name string) (Module, bool) {
	mb.mutex.RLock()
	defer mb.mutex.RUnlock()
	module, exists := mb.modules[name]
	return module, exists
}
func (mb *ModuleBase) RemoveModule(name string) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	delete(mb.modules, name)
}
func (mb *ModuleBase) UploadModule(name string, module Module) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	mb.SetModule(name, module)
}

func (mb *ModuleBase) InitializeAll(c context.Context) error {
	mb.mutex.RLock()
	defer mb.mutex.RUnlock()
	for _, module := range mb.modules {
		if err := module.Initialize(c); err != nil {
			return err
		}
	}
	return nil
}

func (mb *ModuleBase) DeInitializeAll() error {
	mb.mutex.RLock()
	defer mb.mutex.RUnlock()
	for _, module := range mb.modules {
		if err := module.DeInitialize(); err != nil {
			return err
		}
	}
	return nil
}
