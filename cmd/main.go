package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test/pkg/api"
	"test/pkg/mod"
)

func main() {

	fmt.Println("App Start:", time.Now().Format("2006-01-02 15:04:05"))

	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	modeules := mod.NewModuleBase()

	// 初始化消息
	message := mod.NewMessage("main", "init")
	fmt.Println("Message: ", message.Name, message.Msg, message.Time)

	// 初始化消息总线
	bus := mod.GetBus()
	// 初始数系统
	system := mod.NewSystem(bus, ctx, cancel)
	// 初始化牛马
	worker := mod.NewWorker(bus, ctx, cancel)
	// 初始化api网关
	gateway := api.NewGateway(bus, ctx)

	modeules.SetModule("ApiGateway", gateway)
	modeules.SetModule("System", system)
	modeules.SetModule("Worker", worker)

	// 初始化
	go initialize(ctx, modeules)

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 取消上下文，通知所有模块退出
	deinitialize(modeules)
	fmt.Println("App PowerOff")
}

func initialize(ctx context.Context, modules *mod.ModuleBase) {
	modules.InitializeAll(ctx)

	// 每秒发布一次模块状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// 当上下文被取消时，退出循环
		case <-ctx.Done():
			return
		// 当定时器触发时，打印当前时间
		case <-ticker.C:
			fmt.Println("App Status:", time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}

func deinitialize(modules *mod.ModuleBase) {
	modules.DeInitializeAll()
}
