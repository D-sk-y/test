package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("App Start:", time.Now().Format("2006-01-02 15:04:05"))

	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化
	go initialize(ctx)

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("App PowerOff")
}

func initialize(ctx context.Context) {
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
